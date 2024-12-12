package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkRegisterHandler(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(registerHandler))
	defer ts.Close()

	client := &http.Client{}
	url := ts.URL + "/register"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(`{"username":"validuser","email":"user@example.com","password":"validpass"}`))
			if err != nil {
				b.Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			_ = res.Body.Close()

			if res.StatusCode != http.StatusCreated {
				b.Errorf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
			}
		}
	})
}
