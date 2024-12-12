package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const largeStringLength = 1000

func TestRegisterHandler(t *testing.T) {
	// Helper to generate large strings
	generateLargeString := func(length int) string {
		return strings.Repeat("a", length)
	}

	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid registration",
			requestBody:  `{"username":"validuser","email":"user@example.com","password":"validpass"}`,
			expectedCode: http.StatusCreated,
			expectedBody: `{"message":"User registered successfully"}`,
		},
		{
			name:         "Invalid email",
			requestBody:  `{"username":"validuser","email":"invalidemail","password":"validpass"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid email address",
		},
		{
			name:         "Very long username",
			requestBody:  `{"username":"` + generateLargeString(256) + `","email":"user@example.com","password":"validpass"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Username must be between 3 and 255 characters",
		},
		{
			name:         "Very long email",
			requestBody:  `{"username":"validuser","email":"` + generateLargeString(321) + `@example.com","password":"validpass"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid email address",
		},
		{
			name:         "Very long password",
			requestBody:  `{"username":"validuser","email":"user@example.com","password":"` + generateLargeString(129) + `"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Password must be between 6 and 128 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(registerHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.expectedCode)
			}

			if !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("handler returned wrong response body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
