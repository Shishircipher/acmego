package main

import (
	"fmt"
        "log"
	"encoding/json"
	"github.com/shishircipher/acmego/client"
	"github.com/shishircipher/acmego/challenges"
	"github.com/shishircipher/acmego/registration"
	"github.com/shishircipher/acmego/certificate"
	"github.com/shishircipher/acmego/config"
	"github.com/shishircipher/acmego/log"
)

func main () {
	fmt.Println("Welcome to acmego")
        email, domain, accountprivateKey, err := registration.GetUserAccount()
	if err != nil {
		logger.Fatalf("could not get user account %v ", err)
	}
	fmt.Printf("file content - email %s and domain %s \n", email, domain)
	
	client1 := client.CreateDefaultHTTPClient()
	ourUserAgent := "xenolf-acme/4.21.0"
	
	doer := client.NewDoer(client1, ourUserAgent)


	env := "staging"
	var directoryURL string

	if env == "staging" {
		directoryURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	} else {
		directoryURL = "https://acme-v02.api.letsencrypt.org/directory"
	}
	directory, err := config.FetchDirectory(directoryURL, doer)
	if err != nil {
		logger.Fatalf("Error fetching ACME directory: %v", err)
	}
//	fmt.Printf(" directory url : %v+", directory)


	nonceUrl := directory.NewNonce
	manager := client.NewManager(doer, nonceUrl)
	nonce, err := manager.Nonce()
	if err != nil {
		logger.Fatalf("no nounce %v", err)
	}
	logger.Info("%v",nonce)

		 // Define ACME create account URL
        acmeNewAccountUrl := directory.NewAccount
       	// Create payload for account creation
       	payloadAccount := map[string]interface{}{
               	"termsOfServiceAgreed": true,
               	"contact": []string{email,
			},
        	}
        	// Marshal payload
        payloadAccountBytes, err := json.Marshal(payloadAccount)
       	if err != nil {
                	logger.Fatalf("Failed to marshal payload: %v", err)
        }
        	//log.Printf(" %v ",payloadAccountBytes)
	location := ""
	resp, location1, manager, err := client.PostPayload(doer, acmeNewAccountUrl, payloadAccountBytes, accountprivateKey, location ,manager)
	if err != nil {
		logger.Fatalf("failed to create account %v", err)
	}
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
                logger.Fatalf("Failed to marshal payloadOrder: %v", err)
        }
	acmeNewOrderUrl := directory.NewOrder
        responseOrder, location1, manager , err := client.PostPayload(doer, acmeNewOrderUrl, payloadOrderBytes, accountprivateKey, location1 ,manager)
	if err != nil {
		logger.Fatalf("failed to create new order %v", err)
	}
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
	manager, dnsTxt, err := challenges.DNS01Challenges(domain, authURL, doer, accountprivateKey, accountLocation, manager)
	if err != nil {
		logger.Fatalf(" failed to dns challenge %v , Please retry", err)
	}
	logger.Info("dns txt : %v ", dnsTxt)

	// You can add your api here using dnsTxt.



	logger.Info("Manager: %+v\n", &manager)
//	time.Sleep(300 * time.Second)
        
	finalizeStr := responseOrder["finalize"].(string)


	if err := certificate.CSRRequest(finalizeStr, doer, accountLocation, domain ,accountprivateKey,  manager); err != nil {
		fmt.Printf("failed to create certficate %v", err)
	}	
	fmt.Println(" Done ")

}
