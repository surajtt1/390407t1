package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

const largeStringLength = 10000 // Arbitrarily large length for testing

func TestRegisterHandler(t *testing.T) {
    tests := []struct {
        name           string
        requestBody   string
        expectedCode  int
    }{
        {
            name: "Valid registration",
            requestBody: `{"username":"validuser","email":"user@example.com","password":"validpass"}`,
            expectedCode: http.StatusCreated,
        },
        {
            name: "Short username",
            requestBody: `{"username":"ab","email":"user@example.com","password":"validpass"}`,
            expectedCode: http.StatusBadRequest,
        },
        {
            name: "Invalid email",
            requestBody: `{"username":"validuser","email":"invalidemail","password":"validpass"}`,
            expectedCode: http.StatusBadRequest,
        },
        {
            name: "Empty password",
            requestBody: `{"username":"validuser","email":"user@example.com","password":""}`,
            expectedCode: http.StatusBadRequest,
        },
        {
            name: "Missing fields",
            requestBody: `{"username":"","email":"","password":""}`,
            expectedCode: http.StatusBadRequest,
        },
        {
            name: "Very long username",
            requestBody: `{"username":"` + string(make([]byte, largeStringLength)) + `","email":"user@example.com","password":"validpass"}`,
            expectedCode: http.StatusBadRequest,
        },
        {
            name: "Very long email",
            requestBody: `{"username":"validuser","email":"` + string(make([]byte, largeStringLength)) + `@example.com","password":"validpass"}`,
            expectedCode: http.StatusBadRequest,
        },
        {
            name: "Very long password",
            requestBody: `{"username":"validuser","email":"user@example.com","password":"` + string(make([]byte, largeStringLength)) + `"}`,
            expectedCode: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tt.requestBody))
            if err != nil {
                t.Fatalf("could not create request: %v", err)
            }

            rr := httptest.NewRecorder()
            handler := http.HandlerFunc(registerHandler)
            handler.ServeHTTP(rr, req)

            if status := rr.Code; status != tt.expectedCode {
                t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
            }
        })
    }
}
