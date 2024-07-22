// auth/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if creds.Username == "admin" && creds.Password == "password" {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "example_token",
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func main() {
	http.HandleFunc("/login", login)
	log.Println("Auth service listening on :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
