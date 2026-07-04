package models

import "time"

// Article is the off-chain representation stored in SQLite.
// The ContentHash/ArticleID/IPFSCID fields mirror what's anchored
// on-chain in NewsRegistry.sol, so any tampering with Body can be
// detected by recomputing the hash and comparing against the chain.
type Article struct {
	ID          string    `json:"id"`           // hex articleId (bytes32) used on-chain
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Body        string    `json:"body"`
	Author      string    `json:"author"`        // wallet address (0x...)
	ContentHash string    `json:"contentHash"`    // hex keccak256(Body) anchored on-chain
	IPFSCid     string    `json:"ipfsCid"`
	TxHash      string    `json:"txHash"`         // on-chain publish transaction hash
	PublishedAt time.Time `json:"publishedAt"`
	EditCount   int       `json:"editCount"`
	Retracted   bool      `json:"retracted"`
	Verified    bool      `json:"verified"`       // result of last on-chain verification check
}

// PublishRequest is the payload clients send to create an article.
type PublishRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

// AmendRequest is used to correct an existing article.
type AmendRequest struct {
	Body string `json:"body" binding:"required"`
}
