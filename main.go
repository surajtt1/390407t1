package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

// User represents a user object
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// registerHandler handles user registration requests.
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate input fields
	if len(user.Username) < 3 || len(user.Username) > 255 {
		http.Error(w, "Username must be between 3 and 255 characters", http.StatusBadRequest)
		return
	}
	if len(user.Email) == 0 || len(user.Email) > 320 || !emailRegex.MatchString(user.Email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}
	if len(user.Password) < 6 || len(user.Password) > 128 {
		http.Error(w, "Password must be between 6 and 128 characters", http.StatusBadRequest)
		return
	}

	// Registration successful
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message":"User registered successfully"}`)
}
