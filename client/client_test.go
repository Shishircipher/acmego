package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCreateDefaultHTTPClient(t *testing.T) {
	client := CreateDefaultHTTPClient()
	resp, err := client.Get("https://google.com")
	if err != nil {
		t.Errorf("Error making GET request: %v", err)
		return
	}
	defer resp.Body.Close() // Ensure response body is closed

	// Read response body (optional for test purposes)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
		return
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", string(body))
}

func TestCreateDefaultHTTPClient2(t *testing.T) {
	client := CreateDefaultHTTPClient()
	req, err := http.NewRequest("GET", "http://google.com", nil)
	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	resp, err := client.Do(req)
	 if err != nil {
                t.Errorf("Error making GET request: %v", err)
                return
        }
        defer resp.Body.Close() // Ensure response body is closed

        // Read response body (optional for test purposes)
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                t.Errorf("Error reading response body: %v", err)
                return
        }

        fmt.Printf("Response Status: %s\n", resp.Status)
        fmt.Printf("Response Body: %s\n", string(body))
}
