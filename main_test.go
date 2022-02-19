package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func create_valid_hmac(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	sig := mac.Sum(nil)

	return sig
}

func TestAuth(t *testing.T) {
	key := "Super5ecr3t!"
	os.Setenv("SECRET_TOKEN", key)
	data := []struct {
		body       string
		method     string
		statusCode int
	}{
		{
			body:       "this could be a super interesting message",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
		},
		{
			body:       "oh no this uses wrong methods",
			method:     http.MethodGet,
			statusCode: http.StatusUnauthorized,
		},
	}

	for _, item := range data {
		request := httptest.NewRequest(item.method, "/", strings.NewReader(item.body))
		request.Header.Set("X-Hub-Signature-256", "sha256="+string(create_valid_hmac([]byte(item.body), []byte(key))))
		responseRec := httptest.NewRecorder()

		auth_request_handler(responseRec, request)

		if responseRec.Code != item.statusCode {
			t.Errorf("Invalid status %d expected %d", responseRec.Code, item.statusCode)
		}
	}
}
