package registration

import (
	"os"
	"log"
	"github.com/shishircipher/acmego/certcrypto"
	"encoding/json"
	"crypto/ecdsa"
	"github.com/shishircipher/acmego/log"

)
func GetUserAccount() ( string, string, *ecdsa.PrivateKey) {
	file, err := os.ReadFile("data.json") // For read access.
        if err != nil {
        log.Fatalf("file could not read %v ",err)
        }
        //var content interface{}
        var content map[string]string
        err1 := json.Unmarshal(file, &content)
        if err1 != nil {
        log.Fatalf("file could not read content %v ",err)
        }
        domain := content["domain"]
        email := content["email"]
        logger.Info("file content - email %s and domain %s \n", email, domain)
        // Create the private key for new account

        privateKey, errprivateKey := os.ReadFile("./.acmego/account/account.key")
//	log.Printf( "private key is %s",privateKey) //Do not print in production environment
        if errprivateKey != nil {
                privateKey1, err:= certcrypto.GeneratePrivateKey("P256")
                if err != nil {
                log.Fatalf("could not create privateKey %v ",err)
                }
                privateKey = certcrypto.PEMEncode(privateKey1)
                errprivateKeyWrite := os.WriteFile("./.acmego/account/account.key", privateKey, 0600)
                if errprivateKeyWrite != nil {
                log.Fatalf("could not create file %v ", errprivateKeyWrite)
                }
                log.Fatalf("privateKey could not read %v ", errprivateKey)
        }
	accountprivateKey := certcrypto.ReadECKey("./.acmego/account/account.key")
	return  email, domain, accountprivateKey
}
