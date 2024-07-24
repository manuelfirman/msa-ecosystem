package handler

import (
	"auth/internal/crypt"
	"auth/internal/domain"
	"auth/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Headers struct {
	RequestID    string
	ForwardedFor string
}

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

func (ad *AuthDefault) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve and check X-Request-ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			http.Error(w, "X-Request-ID header is required", http.StatusBadRequest)
			return
		}

		// Retrieve and check X-Trace-Info
		traceInfo := r.Header.Get("X-Trace-Info")
		if traceInfo == "" {
			http.Error(w, "X-Trace-Info header is required", http.StatusBadRequest)
			return
		}

		// Append service-specific trace info
		traceInfo += fmt.Sprintf(", http://%s%s%s", os.Getenv("SERVICE_NAME"), os.Getenv("PORT"), r.URL.Path)

		// Add headers data to the context
		ctx := context.WithValue(context.Background(), "request_id", requestID)
		ctx = context.WithValue(ctx, "trace_info", traceInfo)

		// Process the request body
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

		// Deserialize and process the user
		user := deserializeUser(userJson)
		user.Password, err = crypt.HashPassword(user.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		err = ad.sv.Register(&user, ctx)
		if err != nil {
			switch err {
			case service.ErrServiceDuplicateEntry:
				http.Error(w, "Duplicate entry", http.StatusConflict)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Serialize and send the response
		data := map[string]any{
			"message": "user created",
			"data":    serializeUser(user),
		}

		response, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Error serializing response", http.StatusInternalServerError)
			return
		}

		userTraceInfo := "http://user_service:5001/users"

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", requestID)
		traceInfo += fmt.Sprintf(", %s, http://%s%s%s", userTraceInfo, os.Getenv("SERVICE_NAME"), os.Getenv("PORT"), r.URL.Path)
		w.Header().Set("x-trace-info", traceInfo)
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

func getPortFromRequest(r *http.Request) string {
	// Extract the port from the Host header
	host := r.Host
	// Split the Host header into host and port
	parts := strings.Split(host, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	// Return default port for HTTP and HTTPS if not specified
	return "80" // Default for HTTP
}
