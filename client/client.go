package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"io"
	"bytes"
	"crypto"
//	"log"
	"strings"
	"time"
	"runtime"
	"encoding/json"
	"github.com/shishircipher/acmego/acme"
	"github.com/shishircipher/acmego/log"
)
const (
	// ourUserAgent is the User-Agent of this underlying library package.
	ourUserAgent = "xenolf-acme/4.21.0"

	// ourUserAgentComment is part of the UA comment linked to the version status of this underlying library package.
	// values: detach|release
	// NOTE: Update this with each tagged release.
	ourUserAgentComment = "detach"
)
type RequestOption func(*http.Request) error

type Doer struct {
	httpClient *http.Client
	userAgent string
}

//NewDoer Creates a New Doer 
func NewDoer(client *http.Client, userAgent string) *Doer {
	return &Doer {
		httpClient : client,
		userAgent : userAgent,
	}
}




//https://pkg.go.dev/net/http@go1.23.4#hdr-Clients_and_Transports

func CreateDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout : 2 * time.Minute,
		Transport : &http.Transport {
			Proxy : http.ProxyFromEnvironment,
			DialContext : (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   30 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			TLSClientConfig: &tls.Config{
			//	ServerName: os.Getenv(caServerNameEnvVar),
				RootCAs:    initCertPool(),
			},
		},
	}
}

func initCertPool() *x509.CertPool {
	return nil
}

func contentType(ct string) RequestOption {
	return func(req *http.Request) error {
		req.Header.Set("Content-Type", ct)
		return nil
	}
}
// Get performs a GET request with a proper User-Agent string.
// If "response" is not provided, callers should close resp.Body when done reading from it.
func (d *Doer) Get(url string, response interface{}) (*http.Response, error) {
	req, err := d.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return d.do(req, response)
}
// Head performs a HEAD request with a proper User-Agent string.
// The response body (resp.Body) is already closed when this function returns.
func (d *Doer) Head(url string) (*http.Response, error) {
	req, err := d.newRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}

	return d.do(req, nil)
}
// Post performs a POST request with a proper User-Agent string.
// If "response" is not provided, callers should close resp.Body when done reading from it.
func (d *Doer) Post(url string, body io.Reader, bodyType string, response interface{}) (*http.Response, error) {
	req, err := d.newRequest(http.MethodPost, url, body, contentType(bodyType))
	if err != nil {
		return nil, err
	}

	return d.do(req, response)
}
func (d *Doer) newRequest(method, uri string, body io.Reader, opts ...RequestOption) (*http.Request, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", d.formatUserAgent())

	for _, opt := range opts {
		err = opt(req)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
	}

	return req, nil
}
func (d *Doer) do(req *http.Request, response interface{}) (*http.Response, error) {
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err = checkError(req, resp); err != nil {
		return resp, err
	}
	if response != nil {
		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}
		defer resp.Body.Close()
		err = json.Unmarshal(raw, response)
		if err != nil {
			return resp, fmt.Errorf("failed to unmarshal %q to type %T: %w", raw, response, err)
		}
	}
	return resp, nil
}
// formatUserAgent builds and returns the User-Agent string to use in requests.
func (d *Doer) formatUserAgent() string {
	ua := fmt.Sprintf("%s %s (%s; %s; %s)", d.userAgent, ourUserAgent, ourUserAgentComment, runtime.GOOS, runtime.GOARCH)
	return strings.TrimSpace(ua)
}

func checkError(req *http.Request, resp *http.Response) error {
	if resp.StatusCode >= http.StatusBadRequest {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%d :: %s :: %s :: %w", resp.StatusCode, req.Method, req.URL, err)
		}

		var errorDetails *acme.ProblemDetails
		err = json.Unmarshal(body, &errorDetails)
		if err != nil {
			return fmt.Errorf("%d ::%s :: %s :: %w :: %s", resp.StatusCode, req.Method, req.URL, err, string(body))
		}

		errorDetails.Method = req.Method
		errorDetails.URL = req.URL.String()

		if errorDetails.HTTPStatus == 0 {
			errorDetails.HTTPStatus = resp.StatusCode
		}

		// Check for errors we handle specifically
		if errorDetails.HTTPStatus == http.StatusBadRequest && errorDetails.Type == acme.BadNonceErr {
			return &acme.NonceError{ProblemDetails: errorDetails}
		}

		return errorDetails
	}
	return nil

}
//client/client.go:184:84: syntax error: unexpected map after top level declaration
//client/client.go:191:2: syntax error: unexpected name response, expected {
//func PostPayload(doer *Doer, url string, payload bytes, manager *Manager) response map[string]interface , manager *Manager {
//	jwscontent, err := SignContent(acmeURL, payload)
//	if err != nil {
//		log.Printf("failed to create jwscontent",err)
//	}
//	signedBody := bytes.NewBufferString(jwscontent.FullSerialize())
//	var response map[string]interface
//	response, err := doer.Post(url, signedBody, "application/jose+json", response)
//	if err2 != nil {
//		log.Printf("failed to send POST request: %s", err2)
//	}
//	nonce := GetFromResponse(response)
//	manager.Push(nonce)
//	return response, manager

//}
// PostPayload sends a payload to the specified URL using the provided Doer and Manager.
func PostPayload(doer *Doer, url string, payload []byte, privateKey crypto.PrivateKey, location string, manager *Manager) (map[string]interface{}, string, *Manager, error) {
	// Sign the content
	jws := NewJWS(privateKey, location, manager)
	jwsContent, err := jws.SignContent(url, payload)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to create JWS content: %v", err)
	}

	// Prepare the signed body for the POST request
	signedBody := bytes.NewBufferString(jwsContent.FullSerialize())
	var post interface {}
	// Make the POST request
	response, err := doer.Post(url, signedBody, "application/jose+json", post)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to send POST request: %v", err)
	}
	// Handle the response
        if response.StatusCode != http.StatusCreated {
                logger.Info(" status code: %v", response.StatusCode)
        }

	// Extract nonce from the response and push it to the Manager
	nonce, err := GetFromResponse(response)
	if err != nil {
                return nil, "", nil, fmt.Errorf("failed to get nounce: %v", err)
        }
	manager.Push(nonce)
	// Read the body content
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
            return nil, "", nil, fmt.Errorf("failed to bodybytes: %v", err)
	}
	location = response.Header.Get("Location")
	for key, values := range response.Header {
		for _, value := range values {
			logger.Info("  %s: %s\n", key, value)
		}
	}
	// Declare a variable to store the response
        var responseData map[string]interface{}

        if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
	        return nil, "", nil, fmt.Errorf("failednd orderrequest: %s", err)
        }
	// Return the response and the updated Manager
	return responseData, location, manager, nil
}

func (d *Doer) GetResponse(url string) (*http.Response, error) {
    // Create a new HTTP request
    req, err := d.newRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }

    // Perform the HTTP request
    resp, err := d.httpClient.Do(req)
    if err != nil {
        return nil, err
    }

    // Check for HTTP or API-specific errors
    if err = checkError(req, resp); err != nil {
        return resp, err
    }

    // Return the response
    return resp, nil
}

