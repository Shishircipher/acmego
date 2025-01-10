package secure

import (
	"log"
	"testing"
	"net/http"
//	"fmt"
	"bytes"
	"encoding/json"
	"github.com/shishircipher/acmego/client"
	"github.com/shishircipher/acmego/client/internal/nonces"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/shishircipher/acmego/acme"

)

func TestJwsAccountCreation(t *testing.T) {
	client1 := client.CreateDefaultHTTPClient()
	ourUserAgent := "xenolf-acme/4.21.0"
	doer := client.NewDoer(client1, ourUserAgent)
	nonceUrl := "https://acme-staging-v02.api.letsencrypt.org/acme/new-nonce"
	manager := nonces.NewManager(doer, nonceUrl)
	log.Println(manager)
//	nonce := manager.nonces.Nonce()
//	log.Println(nonce)
        privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
	//	return nil, fmt.Errorf("failed to generate ECDSA private key: %w", err)
	        log.Fatal(err)
	}
	jws1 := NewJWS(privateKey, "", manager)
	

	// Define ACME directory URL
	acmeURL := "https://acme-staging-v02.api.letsencrypt.org/acme/new-acct"
	// Create payload for account creation
	payload := map[string]interface{}{
		"termsOfServiceAgreed": true,
		"contact": []string{
			"mailto:admin@example.com",// put tour real contact even in staging environment , if not then - urn:ietf:params:acme:error:invalidContact :: Error creating new account :: contact email has forbidden domain "example.org"
		},
	}
	// Marshal payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}
	jwscontent, err := jws1.SignContent(acmeURL, payloadBytes)
	log.Println(jwscontent)
//	serializedContent := jwscontent.FullSerialize()
        signedBody := bytes.NewBufferString(jwscontent.FullSerialize())
//	req, err := client1.NewRequest("POST", acmeURL, signedBody)
//	req, err := http.NewRequest("POST", acmeURL, bytes.NewBuffer([]byte(serializedContent)))
//	req, err := http.NewRequest("POST", acmeURL, bytes.NewBuffer(jwscontent))
	if err != nil {
		 log.Printf("failed to create POST request: %s", err)
	}
//	req.Header.Set("Content-Type", "application/jose+json")
//	req.Header.Set("User-Agent", "xenolf-acme/4.21.0")
//	resp, err := client1.Do(req)
       // response := 
        var account acme.ExtendedAccount
        resp, err := doer.Post(acmeURL, signedBody, "application/jose+json", account)
// fix this error - json: Unmarshal(non-pointer acme.ExtendedAccount)
	if err != nil {
		log.Printf("failed to send POST request: %s", err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusCreated {
		log.Printf("unexpected status code: %d", resp.StatusCode)
	}

//	log.Println("Account successfully created")
}
