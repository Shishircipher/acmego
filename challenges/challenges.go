package challenges

import (
	"fmt"
	"github.com/shishircipher/acmego/client"
	"crypto/sha256"
        "encoding/base64"
	"log"
	"net/http"
	"io"
	"crypto"
	"encoding/json"
	"github.com/shishircipher/acmego/log"
)

func DNS01Challenges(domain string, authURL string, doer *client.Doer, privateKey crypto.PrivateKey, location string, manager *client.Manager ) (*client.Manager ) {
	jws := client.NewJWS(privateKey,location , manager)
	var challenge interface {}
        response1, err := doer.Get(authURL, challenge)
	logger.Info(" %v ",response1.Header)
	if response1.StatusCode != http.StatusCreated {
                logger.Info("status code: %v", response1.StatusCode)
        }
        bodyBytes, err := io.ReadAll(response1.Body)
        if err != nil {
             log.Fatalf("failed to bodybytes: %v", err)
        }
        // Log the raw body (optional, useful for debugging)
        logger.Info("Raw Body: %s", string(bodyBytes))
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
	log.Println(challurl)
	// Define ANSI color codes
    	red := "\033[31m"
    //	green := "\033[32m"
   //	yellow := "\033[33m"
   // 	blue := "\033[34m"
	reset := "\033[0m" // Reset to default color
	
	log.Printf("domain name : %v \n",domain)
	fmt.Printf(" Paste the red text  in domain management portal (time limit is 5 minutes) :- \n")
//	fmt.Printf(red + "%s" + reset, dnstxt1)
	fmt.Printf("%s%s%s", red, dnstxt1, reset)
	fmt.Println()
	fmt.Println("wait for 300 seconds")
	logger.Spinner(300)
	payloadEmptyBytes := []byte("{}")
	responseChallenge, location, manager := client.PostPayload(doer, challurl, payloadEmptyBytes, privateKey, location ,manager)
	logger.Info(" %v ",location)
	logger.Info("responseChallenge: %+v\n", responseChallenge)
	return manager
}



func getTXTValue(keyAuthorization string) string {
    hash := sha256.Sum256([]byte(keyAuthorization))
    return base64.RawURLEncoding.EncodeToString(hash[:])
}
