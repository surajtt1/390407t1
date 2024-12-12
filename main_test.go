package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
	}{
		{
			name:         "Valid registration",
			requestBody:  `{"username":"validuser","email":"user@example.com","password":"validpass"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Short username",
			requestBody:  `{"username":"ab","email":"user@example.com","password":"validpass"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid email",
			requestBody:  `{"username":"validuser","email":"invalidemail","password":"validpass"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Empty password",
			requestBody:  `{"username":"validuser","email":"user@example.com","password":""}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Missing fields",
			requestBody:  `{"username":"","email":"","password":""}`,
			expectedCode: http.StatusBadRequest,
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
