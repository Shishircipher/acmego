package main

import (
	"fmt"
        "log"
	"encoding/json"
	"github.com/shishircipher/acmego/client"
	"github.com/shishircipher/acmego/challenges"
	"github.com/shishircipher/acmego/registration"
	"github.com/shishircipher/acmego/certificate"
	"github.com/shishircipher/acmego/log"
)

func main () {
	fmt.Println("Welcome to acmego")
        email, domain, accountprivateKey := registration.GetUserAccount()
	fmt.Printf("file content - email %s and domain %s \n", email, domain)
	
	client1 := client.CreateDefaultHTTPClient()
	ourUserAgent := "xenolf-acme/4.21.0"
	
	doer := client.NewDoer(client1, ourUserAgent)
	nonceUrl := "https://acme-staging-v02.api.letsencrypt.org/acme/new-nonce"
	manager := client.NewManager(doer, nonceUrl)
	nonce, errnounce := manager.Nonce()
	if errnounce != nil {
		log.Fatalf("no nounce %v", errnounce)
	}
	logger.Info("%v",nonce)

		 // Define ACME create account URL
        acmeNewAccountUrl := "https://acme-staging-v02.api.letsencrypt.org/acme/new-acct"
       	// Create payload for account creation
       	payloadAccount := map[string]interface{}{
               	"termsOfServiceAgreed": true,
               	"contact": []string{email,
			},
        	}
        	// Marshal payload
        payloadAccountBytes, err := json.Marshal(payloadAccount)
       	if err != nil {
                	log.Fatalf("Failed to marshal payload: %v", err)
        }
        	//log.Printf(" %v ",payloadAccountBytes)
	location := ""
	resp, location1, manager := client.PostPayload(doer, acmeNewAccountUrl, payloadAccountBytes, accountprivateKey, location ,manager)
	logger.Info("Manager: %+v\n", resp)
	logger.Info("Manager: %+v\n", &manager)
	logger.Info("%s", location1)
	accountLocation := location1

        // Create new order
	payloadOrder := map[string]interface{}{
		"identifiers": []map[string]string{
			{"type": "dns", "value": domain},
		},
	}
	payloadOrderBytes, err := json.Marshal(payloadOrder)
        if err != nil {
                log.Fatalf("Failed to marshal payloadOrder: %v", err)
        }
	acmeNewOrderUrl := "https://acme-staging-v02.api.letsencrypt.org/acme/new-order"
        responseOrder, location1, manager := client.PostPayload(doer, acmeNewOrderUrl, payloadOrderBytes, accountprivateKey, location1 ,manager)

        logger.Info("responseOrder: %+v\n", responseOrder)
        logger.Info("Manager: %+v\n", &manager)
        logger.Info("%s", location1)
	authURLs := responseOrder["authorizations"]
	logger.Info(" %s ",authURLs)
	var authURL string
//	authURL = authURLs[0].(string)

	if urls, ok := authURLs.([]interface{}); ok {
    	// Now you can safely index into urls
   	authURL = urls[0].(string)
	} else {
    	log.Fatalf("authURLs is not of type []interface{}")
	}
	manager = challenges.DNS01Challenges(domain, authURL, doer, accountprivateKey, accountLocation, manager)
	
	logger.Info("Manager: %+v\n", &manager)
//	time.Sleep(300 * time.Second)
        
	finalizeStr := responseOrder["finalize"].(string)


	if err := certificate.CSRRequest(finalizeStr, doer, accountLocation, domain ,accountprivateKey,  manager); err != nil {
		fmt.Printf("failed to create certficate %v", err)
	}	

}
