package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRegisterHandler_ConcurrentWithMetrics(t *testing.T) {
	// Number of concurrent requests to simulate
	const numRequests = 100
	const simultaneousGoroutines = 10

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numRequests)

	var totalTime int64
	var failureCount atomic.Int32
	requestBody := `{"username":"testuser","email":"testuser@example.com","password":"validpass"}`

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(requestBody))
			if err != nil {
				atomic.AddInt32(&failureCount, 1)
				t.Errorf("could not create request: %v", err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(registerHandler)

			startTime := time.Now()
			// Serve the HTTP request
			handler.ServeHTTP(rr, req)
			endTime := time.Since(startTime)

			totalTime += int64(endTime.Nanoseconds()) / 1_000_000 // in milliseconds

			// You may want to add checks for the response status code if needed
			if status := rr.Code; status != http.StatusCreated {
				atomic.AddInt32(&failureCount, 1)
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	elapsed := time.Since(start)

	// Calculate average response time
	avgResponseTime := int64(0)
	if numRequests > 0 {
		avgResponseTime = totalTime / int64(numRequests)
	}

	t.Logf("Completed %d requests in %s", numRequests, elapsed)
	t.Logf("Average response time per request: %d ms", avgResponseTime)
	t.Logf("Number of failed requests: %d", failureCount.Load())
}
