package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterHandlerEdgeCases(t *testing.T) {
	// Preparing a large string helper function
	generateLargeString := func(length int) string {
		return strings.Repeat("a", length)
	}

	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
	}{
		{
			name:         "Large Username",
			requestBody:  `{"username":"` + generateLargeString(256) + `","email":"user@example.com","password":"validpass"}`,
			expectedCode: http.StatusBadRequest, // We expect validation to fail
		},
		{
			name:         "Large Email",
			requestBody:  `{"username":"validuser","email":"` + generateLargeString(320) + `","password":"validpass"}`,
			expectedCode: http.StatusBadRequest, // We expect validation to fail
		},
		{
			name:         "Large Password",
			requestBody:  `{"username":"validuser","email":"user@example.com","password":"` + generateLargeString(128) + `"}`,
			expectedCode: http.StatusBadRequest, // We expect validation to fail
		},
		{
			name:         "Maximum Acceptable Username",
			requestBody:  `{"username":"aaa","email":"user@example.com","password":"validpass"}`,
			expectedCode: http.StatusCreated, // This should pass
		},
		{
			name:         "Maximum Acceptable Email",
			requestBody:  `{"username":"validuser","email":"a@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com","password":"validpass"}`,
			expectedCode: http.StatusCreated, // This should pass
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(registerHandler)

			// Serve the HTTP request
			handler.ServeHTTP(rr, req)

			// Check if the status code is what we expect
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}
		})
	}
}
