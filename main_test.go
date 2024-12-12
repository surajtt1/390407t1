package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestRegisterHandler_Concurrent(t *testing.T) {
	// Number of concurrent requests to simulate
	const numRequests = 100

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numRequests)

	requestBody := `{"username":"testuser","email":"testuser@example.com","password":"validpass"}`

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(requestBody))
			if err != nil {
				t.Errorf("could not create request: %v", err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(registerHandler)

			// Serve the HTTP request
			handler.ServeHTTP(rr, req)

			// Check the response status code
			if status := rr.Code; status != http.StatusCreated {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	elapsed := time.Since(start)

	t.Logf("Completed %d requests in %s", numRequests, elapsed)
}
