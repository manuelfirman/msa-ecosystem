package handler

import (
	"auth/internal/crypt"
	"auth/internal/domain"
	"auth/internal/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserJSON struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UniqueID  string `json:"unique_id"`
}

type AuthJSON struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	UniqueID string `json:"unique_id"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// AuthDefault is the default implementation of the AuthHandler interface
type AuthDefault struct {
	sv service.AuthService
}

// NewAuthDefault creates a new AuthDefault
func NewAuthDefault(sv service.AuthService) *AuthDefault {
	return &AuthDefault{sv: sv}
}

// Login handles the login request
func (ad *AuthDefault) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var authMap map[string]any
		if err = json.Unmarshal(bytes, &authMap); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err = ValidateKeyExistance(authMap, "username", "password"); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var creds AuthJSON
		if err = json.Unmarshal(bytes, &creds); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		token, err := crypt.NewToken(creds.Username, 5*time.Minute)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			HttpOnly: true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"logged in"}`))
	}
}

// Register handles the register request
func (ad *AuthDefault) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the user from the request body
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var userMap map[string]any
		if err = json.Unmarshal(bytes, &userMap); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err = ValidateKeyExistance(userMap, "username", "email", "password", "first_name", "last_name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var userJson UserJSON
		if err = json.Unmarshal(bytes, &userJson); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// deserialize the user
		user := deserializeUser(userJson)
		// hash the password
		user.Password, err = crypt.HashPassword(user.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// register the user with the service
		err = ad.sv.Register(&user)
		if err != nil {
			switch err {
			case service.ErrServiceDuplicateEntry:
				http.Error(w, "Duplicate entry", http.StatusConflict)
			default:
				http.Error(w, "Error creating user", http.StatusInternalServerError)
			}
			return
		}

		// serialize the user and send the response
		data := map[string]any{
			"message": "user created",
			"data":    serializeUser(user),
		}

		response, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Error serializing response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	}
}

// Verify handles the verify request
func (ad *AuthDefault) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from the request header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusBadRequest)
			return
		}

		// remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// validate the token
		_, err := crypt.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// send a simple response
		response := map[string]string{
			"status":  "success",
			"message": "Token is valid",
		}

		bytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error serializing response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}

// deserializeUser converts a UserJSON struct to a User struct
func deserializeUser(userReq UserJSON) domain.User {
	return domain.User{
		ID:        userReq.ID,
		Username:  userReq.Username,
		Email:     userReq.Email,
		Password:  userReq.Password,
		FirstName: userReq.FirstName,
		LastName:  userReq.LastName,
		UniqueID:  userReq.UniqueID,
	}
}

// serializeUser converts a User struct to a JSON string
func serializeUser(user domain.User) UserJSON {
	return UserJSON{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UniqueID:  user.UniqueID,
	}
}

func ValidateKeyExistance(mp map[string]any, keys ...string) error {
	for _, key := range keys {
		if _, ok := mp[key]; !ok {
			return fmt.Errorf("key %s does not exist", key)
		}
	}
	return nil
}
