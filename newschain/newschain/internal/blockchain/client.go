package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// newsRegistryABI is the minimal ABI needed to call the functions we use.
// In production, generate this with `abigen` from the compiled contract
// artifact instead of hand-writing it.
const newsRegistryABI = `[
	{"type":"function","name":"publishArticle","stateMutability":"nonpayable",
	 "inputs":[{"name":"articleId","type":"bytes32"},{"name":"contentHash","type":"bytes32"},{"name":"ipfsCID","type":"string"}],
	 "outputs":[]},
	{"type":"function","name":"amendArticle","stateMutability":"nonpayable",
	 "inputs":[{"name":"articleId","type":"bytes32"},{"name":"newContentHash","type":"bytes32"},{"name":"newIpfsCID","type":"string"}],
	 "outputs":[]},
	{"type":"function","name":"verify","stateMutability":"view",
	 "inputs":[{"name":"articleId","type":"bytes32"},{"name":"candidateHash","type":"bytes32"}],
	 "outputs":[{"name":"","type":"bool"}]},
	{"type":"function","name":"getArticle","stateMutability":"view",
	 "inputs":[{"name":"articleId","type":"bytes32"}],
	 "outputs":[{"name":"","type":"tuple","components":[
		{"name":"author","type":"address"},
		{"name":"contentHash","type":"bytes32"},
		{"name":"ipfsCID","type":"string"},
		{"name":"publishedAt","type":"uint256"},
		{"name":"editCount","type":"uint256"},
		{"name":"retracted","type":"bool"}
	 ]}]}
]`

type Client struct {
	eth        *ethclient.Client
	contract   *bind.BoundContract
	address    common.Address
	privateKey *ecdsa.PrivateKey
	chainID    *big.Int
}

// Config holds everything needed to talk to the chain.
type Config struct {
	RPCURL          string // e.g. an Alchemy/Infura/Base RPC endpoint
	ContractAddress string // deployed NewsRegistry address
	PrivateKeyHex   string // publisher wallet's private key (server-side signer)
	ChainID         int64  // e.g. 8453 for Base mainnet, 84532 for Base Sepolia
}

func NewClient(cfg Config) (*Client, error) {
	ethC, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("dial rpc: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(newsRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("parse abi: %w", err)
	}

	addr := common.HexToAddress(cfg.ContractAddress)
	bound := bind.NewBoundContract(addr, parsedABI, ethC, ethC, ethC)

	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKeyHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	return &Client{
		eth:        ethC,
		contract:   bound,
		address:    addr,
		privateKey: privKey,
		chainID:    big.NewInt(cfg.ChainID),
	}, nil
}

func (c *Client) transactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(c.privateKey, c.chainID)
	if err != nil {
		return nil, err
	}
	opts.Context = ctx
	return opts, nil
}

// ArticleIDFromSlug derives a deterministic on-chain article ID from an
// author address + slug, so both server and chain agree on IDs without
// needing a round trip.
func ArticleIDFromSlug(author common.Address, slug string) [32]byte {
	data := append(author.Bytes(), []byte(slug)...)
	return crypto.Keccak256Hash(data)
}

// HashContent computes the keccak256 hash of an article body — the same
// hash gets anchored on-chain and recomputed at read time to detect tampering.
func HashContent(body string) [32]byte {
	return crypto.Keccak256Hash([]byte(body))
}

// PublishArticle anchors a new article's hash + IPFS CID on-chain.
// Returns the transaction hash so the caller can store it for reference.
func (c *Client) PublishArticle(ctx context.Context, articleID, contentHash [32]byte, ipfsCID string) (string, error) {
	opts, err := c.transactOpts(ctx)
	if err != nil {
		return "", err
	}
	tx, err := c.contract.Transact(opts, "publishArticle", articleID, contentHash, ipfsCID)
	if err != nil {
		return "", fmt.Errorf("publishArticle tx: %w", err)
	}
	return tx.Hash().Hex(), nil
}

// AmendArticle updates the on-chain hash for a correction/edit.
func (c *Client) AmendArticle(ctx context.Context, articleID, newContentHash [32]byte, newIpfsCID string) (string, error) {
	opts, err := c.transactOpts(ctx)
	if err != nil {
		return "", err
	}
	tx, err := c.contract.Transact(opts, "amendArticle", articleID, newContentHash, newIpfsCID)
	if err != nil {
		return "", fmt.Errorf("amendArticle tx: %w", err)
	}
	return tx.Hash().Hex(), nil
}

// Verify checks whether a given content hash matches what's anchored on-chain
// for this article ID. This is what powers the "Verified on-chain" badge.
func (c *Client) Verify(ctx context.Context, articleID, candidateHash [32]byte) (bool, error) {
	var out []interface{}
	opts := &bind.CallOpts{Context: ctx}
	out = append(out, new(bool))
	err := c.contract.Call(opts, &out, "verify", articleID, candidateHash)
	if err != nil {
		return false, fmt.Errorf("verify call: %w", err)
	}
	result, ok := out[0].(*bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type")
	}
	return *result, nil
}
