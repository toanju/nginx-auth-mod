// SPDX-License-Identifier: MIT

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

func create_request(body string, method string) *http.Request {
	if body == "" {
		return httptest.NewRequest(method, "/", nil)
	} else {
		return httptest.NewRequest(method, "/", strings.NewReader(body))
	}
}

func TestAuth(t *testing.T) {
	key := "Super5ecr3t!"
	data := []struct {
		body       string
		method     string
		add_mac    bool
		valid_mac  bool
		set_token  bool
		statusCode int
	}{
		{
			// 0
			body:       "this could be a super interesting message",
			method:     http.MethodPost,
			valid_mac:  true,
			add_mac:    true,
			set_token:  true,
			statusCode: http.StatusOK,
		},
		{
			// 1
			body:       "oh no this uses wrong methods",
			method:     http.MethodGet,
			valid_mac:  true,
			add_mac:    true,
			set_token:  true,
			statusCode: http.StatusUnauthorized,
		},
		{
			// 3
			body:       "this could be a super interesting message",
			method:     http.MethodPost,
			valid_mac:  false,
			add_mac:    true,
			set_token:  true,
			statusCode: http.StatusUnauthorized,
		},
		{
			// 4
			body:       "this could be a super interesting message",
			method:     http.MethodPost,
			valid_mac:  true,
			add_mac:    false,
			set_token:  true,
			statusCode: http.StatusUnauthorized,
		},
		{
			// 5
			body:       "this could be a super interesting message",
			method:     http.MethodPost,
			valid_mac:  true,
			add_mac:    true,
			set_token:  false,
			statusCode: http.StatusUnauthorized,
		},
	}

	for cnt, item := range data {
		//request := nil
		if item.set_token == true {
			os.Setenv("SECRET_TOKEN", key)
		}

		request := httptest.NewRequest(item.method, "/", strings.NewReader(item.body))
		if item.valid_mac == true && item.add_mac == true {
			request.Header.Set("X-Hub-Signature-256", "sha256="+string(create_valid_hmac([]byte(item.body), []byte(key))))
		} else if item.valid_mac == false && item.add_mac == true {
			request.Header.Set("X-Hub-Signature-256", "sha256="+string(create_valid_hmac([]byte("invalid hmac"), []byte(key))))
		} else {
			// add_mac == false
		}
		responseRec := httptest.NewRecorder()

		auth_request_handler(responseRec, request)

		if responseRec.Code != item.statusCode {
			t.Errorf("Test %d: Invalid status %d expected %d", cnt, responseRec.Code, item.statusCode)
		}

		if item.set_token == true {
			os.Unsetenv("SECRET_TOKEN")
		}
	}
}
