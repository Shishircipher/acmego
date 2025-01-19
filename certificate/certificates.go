package certificate
import (

        "github.com/shishircipher/acmego/client"
	"github.com/shishircipher/acmego/certcrypto"
        "log"
        "io"
	"net/http"
	"fmt"
	"os"
        "encoding/json"
	"encoding/base64"
	"crypto/ecdsa"
	"github.com/shishircipher/acmego/log"
)
func CSRRequest(url string, doer *client.Doer, location string, domain string, privateKey *ecdsa.PrivateKey, manager *client.Manager) error {


	privateKeyDomain, err:= certcrypto.GeneratePrivateKey("2048")
        if err != nil {
                  log.Fatalf(" could not create privateKeyDomain %v ",err)
        }
        // Construct the file path dynamically using the domain name
        filePath := "./.acmego/" + domain + "/privateKey.key"

        // Ensure the directory exists
        errdomaindir := os.MkdirAll("./.acmego/"+domain, 0700)
        if errdomaindir != nil {
                log.Fatalf("Failed to create directory: %v", errdomaindir)
        }
        privateKeyDomainPEM := certcrypto.PEMEncode(privateKeyDomain)
        // Write the private key to the file
        errprivateKeyDomainWrite := os.WriteFile(filePath, privateKeyDomainPEM, 0600)
        if errprivateKeyDomainWrite != nil {
                log.Fatalf("Failed to write private key to file: %v", errprivateKeyDomainWrite)
        }

        logger.Info("Private key written successfully to %v", filePath)
	san := []string{domain}
        mustStaple := true
        // Generate the DER certificate
        csrbytes, err := certcrypto.GenerateCSR(privateKeyDomain,  domain, san, mustStaple)
        if err != nil {
                log.Fatalf("Failed to generate DER certificate: %v\n", err)
        }


        csrEncoded := base64.RawURLEncoding.EncodeToString(csrbytes)
        // finalizeStr
        csrPem := map[string]interface{}{
        "csr": csrEncoded, // Assuming csrPem is a string containing the CSR in PEM format
        }
        payloadCSRBytes, err := json.Marshal(csrPem)
        if err != nil {
                log.Fatalf("Failed to marshal payload: %v", err)
        }
        responseCSR, location1, manager := client.PostPayload(doer, url, payloadCSRBytes, privateKey, location ,manager)

        logger.Info("ResponseCSR: %+v\n", responseCSR)
        logger.Info("Manager: %+v\n", &manager)
	logger.Spinner(60)
	var finalOrder interface{} // The response will be unmarshalled here

// Perform the GET request
        response1, err := doer.Get(location1, finalOrder)
        if err != nil {
                log.Fatalf("HTTP GET request failed: %v", err)
        }
	
	//log.Println(response1.Header)
        if response1.StatusCode != http.StatusCreated {
                log.Printf("status code: %d", response1.StatusCode)
        }
        bodyBytes, err := io.ReadAll(response1.Body)
        if err != nil {
             log.Printf("failed to bodybytes: %s", err)
        }
        // Log the raw body (optional, useful for debugging)
        logger.Info("Raw Body: %v", string(bodyBytes))
        // Parse JSON response
        var response map[string]interface{}
        err = json.Unmarshal(bodyBytes, &response)
        if err != nil {
                log.Fatalf("Failed to parse JSON: %v", err)
        }
// Extract the `certificate` field from the parsed response
        certificateUrl, ok := response["certificate"].(string)
        if !ok || certificateUrl == "" {
                log.Println("Certificate field is missing or not a string")
		fmt.Println("dns-01 challenge failed")
        } else {
         log.Printf("Certificate URL: %s", certificateUrl)
        }



	// Download the certificate if the URL is valid
        if err := downloadCertificate(certificateUrl, domain, doer); err != nil {
        log.Printf("Failed to download certificate: %v", err)
        }
	return nil
}

// Helper function to download the certificate
func downloadCertificate(url string, domain string, doer *client.Doer) error {
    resp, err := doer.GetResponse(url)
    if err != nil {
        return fmt.Errorf("failed to fetch certificate: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code while downloading certificate: %d", resp.StatusCode)
    }
    certBytes, err := io.ReadAll(resp.Body)
    // Write the certificate content to a file
    certFilePath := "./.acmego/" + domain + "/certificate.pem"
    err = os.WriteFile(certFilePath, certBytes, 0600)
    if err != nil {
        return fmt.Errorf("failed to create certificate file: %w", err)
    }

    fmt.Println("Certificate downloaded successfully as 'certificate.pem'")
    return nil
}
