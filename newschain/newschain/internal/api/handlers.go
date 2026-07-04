package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"github.com/kingdev/newschain/internal/auth"
	"github.com/kingdev/newschain/internal/blockchain"
	"github.com/kingdev/newschain/internal/models"
	"github.com/kingdev/newschain/internal/storage"
)

type Server struct {
	store *storage.Store
	chain *blockchain.Client
	ipfs  IPFSPinner
}

// IPFSPinner abstracts content-addressed storage so it can be swapped
// for a real IPFS node / pinning service (Pinata, web3.storage, etc).
type IPFSPinner interface {
	Pin(content []byte) (cid string, err error)
}

func NewServer(store *storage.Store, chain *blockchain.Client, ipfs IPFSPinner) *Server {
	return &Server{store: store, chain: chain, ipfs: ipfs}
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", s.health)

	auth := r.Group("/auth")
	{
		auth.GET("/nonce", s.getNonce)
		auth.POST("/verify", s.verifySignature)
	}

	articles := r.Group("/articles")
	{
		articles.GET("", s.listArticles)
		articles.GET("/:id", s.getArticle)
		articles.GET("/:id/verify", s.verifyArticle)

		authed := articles.Group("")
		authed.Use(s.requireAuth)
		{
			authed.POST("", s.publishArticle)
			authed.PUT("/:id", s.amendArticle)
			authed.DELETE("/:id", s.retractArticle)
		}
	}
}

func (s *Server) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// --- Auth: SIWE-style wallet login ---

func (s *Server) getNonce(c *gin.Context) {
	address := c.Query("address")
	if !auth.IsValidAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address"})
		return
	}
	nonce, err := auth.GenerateNonce(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate nonce"})
		return
	}
	message := auth.SIWEMessage("newschain.io", address, nonce)
	c.JSON(http.StatusOK, gin.H{"nonce": nonce, "message": message})
}

type verifyRequest struct {
	Address   string `json:"address" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

func (s *Server) verifySignature(c *gin.Context) {
	var req verifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := auth.VerifySignature(req.Address, req.Message, req.Signature); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.IssueSessionToken(req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "address": req.Address})
}

func (s *Server) requireAuth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
		return
	}
	token := strings.TrimPrefix(header, "Bearer ")
	claims, err := auth.ParseSessionToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}
	c.Set("wallet", claims.Address)
	c.Next()
}

// --- Articles ---

func (s *Server) listArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	list, err := s.store.ListArticles(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list articles"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"articles": list})
}

func (s *Server) getArticle(c *gin.Context) {
	a, err := s.store.GetArticle(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// verifyArticle recomputes the article's content hash and checks it
// against what's anchored on-chain — this is the tamper-detection endpoint.
func (s *Server) verifyArticle(c *gin.Context) {
	a, err := s.store.GetArticle(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	articleIDBytes, err := hex.DecodeString(strings.TrimPrefix(a.ID, "0x"))
	if err != nil || len(articleIDBytes) != 32 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "corrupt article id"})
		return
	}
	var articleID [32]byte
	copy(articleID[:], articleIDBytes)

	candidateHash := blockchain.HashContent(a.Body)

	verified, err := s.chain.Verify(c.Request.Context(), articleID, candidateHash)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "chain verification failed", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articleId":       a.ID,
		"onChainVerified": verified,
		"recomputedHash":  "0x" + hex.EncodeToString(candidateHash[:]),
		"storedHash":      a.ContentHash,
	})
}

func (s *Server) publishArticle(c *gin.Context) {
	wallet := c.GetString("wallet")

	var req models.PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := slugify(req.Title)
	authorAddr := common.HexToAddress(wallet)
	articleID := blockchain.ArticleIDFromSlug(authorAddr, slug)
	contentHash := blockchain.HashContent(req.Body)

	cid, err := s.ipfs.Pin([]byte(req.Body))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "ipfs pin failed", "detail": err.Error()})
		return
	}

	txHash, err := s.chain.PublishArticle(c.Request.Context(), articleID, contentHash, cid)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "on-chain publish failed", "detail": err.Error()})
		return
	}

	article := &models.Article{
		ID:          "0x" + hex.EncodeToString(articleID[:]),
		Title:       req.Title,
		Slug:        slug,
		Body:        req.Body,
		Author:      wallet,
		ContentHash: "0x" + hex.EncodeToString(contentHash[:]),
		IPFSCid:     cid,
		TxHash:      txHash,
		PublishedAt: time.Now().UTC(),
	}

	if err := s.store.InsertArticle(article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist article"})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func (s *Server) amendArticle(c *gin.Context) {
	wallet := c.GetString("wallet")
	id := c.Param("id")

	existing, err := s.store.GetArticle(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}
	if !strings.EqualFold(existing.Author, wallet) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the original author may amend this article"})
		return
	}

	var req models.AmendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	articleIDBytes, _ := hex.DecodeString(strings.TrimPrefix(existing.ID, "0x"))
	var articleID [32]byte
	copy(articleID[:], articleIDBytes)

	newHash := blockchain.HashContent(req.Body)
	newCID, err := s.ipfs.Pin([]byte(req.Body))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "ipfs pin failed"})
		return
	}

	txHash, err := s.chain.AmendArticle(c.Request.Context(), articleID, newHash, newCID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "on-chain amend failed", "detail": err.Error()})
		return
	}

	if err := s.store.UpdateArticleContent(id, req.Body, "0x"+hex.EncodeToString(newHash[:]), newCID, txHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist amendment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "amended", "txHash": txHash})
}

func (s *Server) retractArticle(c *gin.Context) {
	wallet := c.GetString("wallet")
	id := c.Param("id")

	existing, err := s.store.GetArticle(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}
	if !strings.EqualFold(existing.Author, wallet) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the original author may retract this article"})
		return
	}

	if err := s.store.RetractArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retract"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "retracted"})
}

func slugify(title string) string {
	lower := strings.ToLower(title)
	replaced := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			return r
		case r == ' ' || r == '-':
			return '-'
		default:
			return -1
		}
	}, lower)
	slug := strings.Trim(replaced, "-")

	// append a short random suffix to avoid collisions on duplicate titles
	suffix := make([]byte, 3)
	_, _ = rand.Read(suffix)
	return slug + "-" + hex.EncodeToString(suffix)
}
