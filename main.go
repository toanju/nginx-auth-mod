package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func validMAC(message []byte, messageMAC []byte, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func auth_request_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusUnauthorized)
		log.Printf("Invalid request from %s", r.RemoteAddr)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		log.Printf("Body not readable from %s", r.RemoteAddr)
		return
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	signature = strings.Replace(signature, "sha256=", "", 1)
	if signature == "" {
		http.Error(w, "", http.StatusUnauthorized)
		log.Printf("Signature missing %s", r.RemoteAddr)
		return
	}

	key := os.Getenv("SECRET_TOKEN")
	if key == "" {
		http.Error(w, "", http.StatusUnauthorized)
		log.Printf("Key note found")
		return
	}

	if !validMAC(body, []byte(signature), []byte(key)) {
		http.Error(w, "", http.StatusUnauthorized)
		log.Printf("invalid request from %s", r.RemoteAddr)
		return
	}

	// OK
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", auth_request_handler)
	http.ListenAndServe(":8080", nil)
}
