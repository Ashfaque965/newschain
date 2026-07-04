package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
)

// nonceStore is a simple in-memory nonce cache to prevent signature replay.
// Swap for Redis in a multi-instance deployment.
type nonceStore struct {
	mu     sync.Mutex
	nonces map[string]time.Time
}

var nonces = &nonceStore{nonces: make(map[string]time.Time)}

const nonceTTL = 5 * time.Minute

// GenerateNonce issues a fresh nonce for a wallet address to sign,
// following the SIWE (Sign-In With Ethereum, EIP-4361) pattern.
func GenerateNonce(address string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	nonce := hex.EncodeToString(b)

	nonces.mu.Lock()
	nonces.nonces[strings.ToLower(address)+":"+nonce] = time.Now().Add(nonceTTL)
	nonces.mu.Unlock()

	return nonce, nil
}

func consumeNonce(address, nonce string) bool {
	key := strings.ToLower(address) + ":" + nonce
	nonces.mu.Lock()
	defer nonces.mu.Unlock()

	expiry, ok := nonces.nonces[key]
	if !ok || time.Now().After(expiry) {
		return false
	}
	delete(nonces.nonces, key)
	return true
}

// SIWEMessage builds the human-readable message the wallet signs client-side.
func SIWEMessage(domain, address, nonce string) string {
	return fmt.Sprintf(
		"%s wants you to sign in with your Ethereum account:\n%s\n\nSign in to NewsChain.\n\nNonce: %s\nIssued At: %s",
		domain, address, nonce, time.Now().UTC().Format(time.RFC3339),
	)
}

// VerifySignature checks that `signature` over `message` was produced by
// the private key controlling `address`, and that the nonce embedded in
// the message hasn't been used before (replay protection).
func VerifySignature(address, message, signature string) error {
	sigBytes, err := hexToBytes(signature)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}
	if len(sigBytes) != 65 {
		return errors.New("invalid signature length")
	}
	// Ethereum signatures use v in {27,28}; go-ethereum expects {0,1}.
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	nonce := extractNonce(message)
	if nonce == "" || !consumeNonce(address, nonce) {
		return errors.New("invalid or expired nonce")
	}

	hash := ethSignedMessageHash(message)
	pubKey, err := crypto.SigToPub(hash, sigBytes)
	if err != nil {
		return fmt.Errorf("recover pubkey: %w", err)
	}

	recovered := crypto.PubkeyToAddress(*pubKey)
	if !strings.EqualFold(recovered.Hex(), address) {
		return errors.New("signature does not match address")
	}
	return nil
}

func ethSignedMessageHash(message string) []byte {
	prefixed := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	return crypto.Keccak256([]byte(prefixed))
}

func extractNonce(message string) string {
	for _, line := range strings.Split(message, "\n") {
		if strings.HasPrefix(line, "Nonce: ") {
			return strings.TrimPrefix(line, "Nonce: ")
		}
	}
	return ""
}

func hexToBytes(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	return hex.DecodeString(s)
}

// --- JWT session issuance after successful wallet verification ---

var jwtSecret = []byte("REPLACE_WITH_ENV_SECRET") // load from env/secret manager in production

type Claims struct {
	Address string `json:"address"`
	jwt.RegisteredClaims
}

func IssueSessionToken(address string) (string, error) {
	claims := Claims{
		Address: strings.ToLower(address),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseSessionToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// IsValidAddress is a light sanity check used before nonce generation.
func IsValidAddress(address string) bool {
	return common.IsHexAddress(address)
}
