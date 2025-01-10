package nonces

import (
	"net/http"
//	"net/http/httptest"
	"testing"
//	"time"
	"log"
	"github.com/shishircipher/acmego/client"
)


func TestNonceManager(t *testing.T) {
	doer := client.NewDoer(http.DefaultClient, "lego-test")
        nonce_url := "https://acme-staging-v02.api.letsencrypt.org/acme/new-nonce"
	j := NewManager(doer, nonce_url)

        nonce, err := j.Nonce()
        if err != nil {
            t.Fatalf("Failed to fetch nonce: %v", err)
         }

        log.Printf("Nonce successfully fetched: %s", nonce)
//	j, err := j.Nonce()
//	k := j.Pop()
//	log.Println(k)
//        nonce, ok := j.Pop()
//      if ok {
//        log.Printf("Nonce successfully popped: %s", nonce)
 //  } else {
     //  log.Println("No nonce available to pop")
   //  }




}
