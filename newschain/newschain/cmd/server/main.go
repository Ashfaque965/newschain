package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kingdev/newschain/internal/api"
	"github.com/kingdev/newschain/internal/blockchain"
	"github.com/kingdev/newschain/internal/storage"
)

func main() {
	dbPath := getEnv("DB_PATH", "newschain.db")
	rpcURL := getEnv("RPC_URL", "https://sepolia.base.org")
	contractAddr := getEnv("CONTRACT_ADDRESS", "")
	privateKey := getEnv("SIGNER_PRIVATE_KEY", "")
	ipfsAPI := getEnv("IPFS_API_URL", "")
	chainID, _ := strconv.ParseInt(getEnv("CHAIN_ID", "84532"), 10, 64) // Base Sepolia default
	port := getEnv("PORT", "8080")

	store, err := storage.New(dbPath)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}
	defer store.Close()

	if contractAddr == "" || privateKey == "" {
		log.Fatal("CONTRACT_ADDRESS and SIGNER_PRIVATE_KEY must be set")
	}

	chain, err := blockchain.NewClient(blockchain.Config{
		RPCURL:          rpcURL,
		ContractAddress: contractAddr,
		PrivateKeyHex:   privateKey,
		ChainID:         chainID,
	})
	if err != nil {
		log.Fatalf("failed to init blockchain client: %v", err)
	}

	var ipfs api.IPFSPinner
	if ipfsAPI == "" {
		log.Println("IPFS_API_URL not set — using MockIPFSPinner (dev mode only)")
		ipfs = storage.MockIPFSPinner{}
	} else {
		ipfs = storage.NewHTTPIPFSPinner(ipfsAPI)
	}

	server := api.NewServer(store, chain, ipfs)

	r := gin.Default()
	server.RegisterRoutes(r)

	log.Printf("NewsChain API listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
