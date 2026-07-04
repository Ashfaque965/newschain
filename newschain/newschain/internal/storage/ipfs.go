package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

// HTTPIPFSPinner talks to an IPFS HTTP API (e.g. a local Kubo node at
// :5001, or a gateway that exposes the same /api/v0/add endpoint).
// For managed pinning (Pinata, web3.storage) swap this out for their
// respective client SDKs — the Pin(content) interface stays the same.
type HTTPIPFSPinner struct {
	APIURL string // e.g. "http://127.0.0.1:5001"
}

func NewHTTPIPFSPinner(apiURL string) *HTTPIPFSPinner {
	return &HTTPIPFSPinner{APIURL: apiURL}
}

type addResponse struct {
	Hash string `json:"Hash"`
}

func (p *HTTPIPFSPinner) Pin(content []byte) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", "article.txt")
	if err != nil {
		return "", err
	}
	if _, err := part.Write(content); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", p.APIURL+"/api/v0/add", &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ipfs add request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ipfs add returned status %d", resp.StatusCode)
	}

	var out addResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode ipfs response: %w", err)
	}
	return out.Hash, nil
}

// MockIPFSPinner is useful for local dev/testing without a running IPFS node.
// It returns a deterministic fake CID derived from content length + a prefix.
type MockIPFSPinner struct{}

func (MockIPFSPinner) Pin(content []byte) (string, error) {
	return fmt.Sprintf("bafymock%d", len(content)), nil
}
