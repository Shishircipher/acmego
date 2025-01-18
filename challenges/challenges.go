package challenges

import (

	"github.com/shishircipher/acmego/client"
	"crypto/sha256"
        "encoding/base64"
	"log"
	"net/http"
	"io"
	"crypto"
	"encoding/json"
	"time"
)

func DNS01Challenges(domain string, authURL string, doer *client.Doer, privateKey crypto.PrivateKey, location string, manager *client.Manager ) (*client.Manager ) {
	jws := client.NewJWS(privateKey,location , manager)
	var challenge interface {}
        response1, err := doer.Get(authURL, challenge)
	log.Println(response1.Header)
	if response1.StatusCode != http.StatusCreated {
                log.Printf("status code: %d", response1.StatusCode)
        }
        bodyBytes, err := io.ReadAll(response1.Body)
        if err != nil {
             log.Printf("failed to bodybytes: %s", err)
        }
        // Log the raw body (optional, useful for debugging)
        log.Printf("Raw Body: %s", string(bodyBytes))
	// Parse JSON response
	var response map[string]interface{}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	var challurl string
	// Extract the token value
	var tokendns string
	if challenges, ok := response["challenges"].([]interface{}); ok {
		for _, challenge := range challenges {
			if challengeMap, ok := challenge.(map[string]interface{}); ok {
				if challengeMap["type"] == "dns-01" {
					tokendns, _ = challengeMap["token"].(string)
					challurl, _ = challengeMap["url"].(string)
					break
				}
			}
		}
	}

	dnstxt, err := jws.GetKeyAuthorization(tokendns)
	if err != nil {
             log.Fatalf("failed to dnstxts: %v", err)
        }
	//log.Println(dnstxt)
        // Define the DNS record details
	dnstxt1 := getTXTValue(dnstxt)
	log.Printf("domain name : %s \n",domain)
	log.Printf(" Paste the text of bracket in domain management portal :- [%s] \n",dnstxt1)
	time.Sleep(300 * time.Second)
	payloadEmptyBytes := []byte("{}")
	responseChallenge, location, manager := client.PostPayload(doer, challurl, payloadEmptyBytes, privateKey, location ,manager)
	log.Println(location)
	log.Printf("responseChallenge: %+v\n", responseChallenge)
	return manager
}



func getTXTValue(keyAuthorization string) string {
    hash := sha256.Sum256([]byte(keyAuthorization))
    return base64.RawURLEncoding.EncodeToString(hash[:])
}
