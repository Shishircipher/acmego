package registration

import (
	"os"
//	"log"
	"fmt"
	"github.com/shishircipher/acmego/certcrypto"
	"encoding/json"
	"crypto/ecdsa"
	"github.com/shishircipher/acmego/log"

)
func GetUserAccount() ( string, string, *ecdsa.PrivateKey, error) {
	file, err := os.ReadFile("data.json") // For read access.
        if err != nil {
        return "", "", nil, fmt.Errorf("file could not read %w ",err)
        }
        //var content interface{}
        var content map[string]string
        err1 := json.Unmarshal(file, &content)
        if err1 != nil {
        return "", "", nil, fmt.Errorf("file could not read content %w ",err)
        }
        domain := content["domain"]
        email := content["email"]
        logger.Info("file content - email %s and domain %s \n", email, domain)
        // Create the private key for new account
	erraccountdir := os.MkdirAll("./.acmego/account/", 0700)
        if erraccountdir != nil {
                return "", "", nil, fmt.Errorf("Failed to create directory: %w", erraccountdir)
        }
        privateKey, err := os.ReadFile("./.acmego/account/account.key")
//	log.Printf( "private key is %s",privateKey) //Do not print in production environment
        if err != nil {
                privateKey1, err:= certcrypto.GeneratePrivateKey("P256")
                if err != nil {
                return "", "", nil, fmt.Errorf("could not create privateKey %w ",err)
                }
                privateKey = certcrypto.PEMEncode(privateKey1)
                errprivateKeyWrite := os.WriteFile("./.acmego/account/account.key", privateKey, 0600)
                if errprivateKeyWrite != nil {
                return "", "", nil, fmt.Errorf("could not create file %w ", errprivateKeyWrite)
                }
         //       log.Fatalf("privateKey could not read %v ", errprivateKey)
        }
	accountprivateKey, err := certcrypto.ReadECKey("./.acmego/account/account.key")
	if err != nil {
        return "", "", nil, fmt.Errorf("failed to read ec key %w ",err)
        }
	return  email, domain, accountprivateKey, nil
}
