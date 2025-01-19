package config

import (
	"github.com/shishircipher/acmego/client"
//	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"io"

)

// DirectoryResponse represents the structure of the ACME directory response
type DirectoryResponse struct {
	NewNonce    string `json:"newNonce"`
	NewAccount  string `json:"newAccount"`
	NewOrder    string `json:"newOrder"`
	RevokeCert  string `json:"revokeCert"`
	KeyChange   string `json:"keyChange"`
	TermsOfService string `json:"meta.termsOfService"`
}
//# Set this in your environment
//export ACME_ENV=staging

//env := os.Getenv("ACME_ENV")
//var directoryURL string

//if env == "staging" {
//	directoryURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
//} else {
//	directoryURL = "https://acme-v02.api.letsencrypt.org/directory"
//}

func FetchDirectory(directoryURL string, doer *client.Doer) (*DirectoryResponse, error) {
	resp, err := doer.GetResponse(directoryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch directory: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory response: %w", err)
	}

	var directory DirectoryResponse
	if err := json.Unmarshal(body, &directory); err != nil {
		return nil, fmt.Errorf("failed to parse directory JSON: %w", err)
	}

	return &directory, nil
}
