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

	// WaitGroup to wait for all goroutines
	var wg sync.WaitGroup
	wg.Add(numRequests)
	requestBody := `{"username":"testuser","email":"testuser@example.com","password":"validpass"}`

	start := time.Now()

	// Variables to store response times and failure counts
	var mu sync.Mutex
	totalResponseTime := 0.0
	failedRequests := 0

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

			// Start timing the request
			reqStart := time.Now()
			handler.ServeHTTP(rr, req)
			reqDuration := time.Since(reqStart).Seconds()

			mu.Lock() // Lock to ensure thread-safe access to shared data
			totalResponseTime += reqDuration
			if status := rr.Code; status != http.StatusCreated {
				failedRequests++
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
			}
			mu.Unlock()
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	elapsed := time.Since(start)

	// Calculate average response time
	averageResponseTime := totalResponseTime / float64(numRequests)

	// Log the results
	t.Logf("Completed %d requests in %s", numRequests, elapsed)
	t.Logf("Average response time per request: %.4f seconds", averageResponseTime)
	t.Logf("Number of failed requests: %d", failedRequests)
}
