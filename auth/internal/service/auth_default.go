package service

import (
	"auth/internal/domain"
	"auth/internal/repository"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

var (
	UserServiceURL = "http://user_service:5001/users"
)

type UserRequest struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UniqueID  string `json:"unique_id"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	ErrorMsg string `json:"error"`
}

// AuthDefault represents the default implementation of the AuthService
type AuthDefault struct {
	rp repository.AuthRepository
}

// NewAuthDefault creates a new instance of AuthDefault
func NewAuthDefault(rp repository.AuthRepository) *AuthDefault {
	return &AuthDefault{rp: rp}
}

func (ad *AuthDefault) Login(auth *domain.Auth) (err error) {
	err = ad.rp.Login(auth)
	if err != nil {
		switch err {
		case repository.ErrRepositoryEntryNotFound:
			err = ErrServiceInvalidCredentials
		default:

		}
		return
	}

	return
}

func (ad *AuthDefault) Register(user *domain.User, ctx context.Context) (err error) {
	err = createUserService(user, ctx)
	if err != nil {
		return err
	}

	var auth domain.Auth
	auth.Username = user.Username
	auth.Password = user.Password

	err = ad.rp.Register(&auth)
	if err != nil {
		switch err {
		case repository.ErrRepositoryDuplicateEntry:
			err = ErrServiceDuplicateEntry
		default:

		}
		return
	}

	return
}

func createUserService(user *domain.User, ctx context.Context) (err error) {

	// create unique ID for the user
	uniqueID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	user.UniqueID = uniqueID.String()

	userRequest := UserRequest{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UniqueID:  user.UniqueID,
	}

	jsonData, err := json.Marshal(userRequest)
	if err != nil {
		return err
	}

	// create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, UserServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	// get the request ID and forwarded for from the context
	requestID, ok := ctx.Value("request_id").(string)
	if !ok {
		return fmt.Errorf("failed to get request ID from context")
	}
	traceInfo, ok := ctx.Value("trace_info").(string)
	if !ok {
		return fmt.Errorf("failed to get trace info for from context")
	}

	// set the headers
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("X-Trace-Info", traceInfo)
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// parse the response body
	var userResponse UserResponse
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("%w: %s", ErrServiceFailedToCreateUser, userResponse.ErrorMsg)
	}

	// asign the user ID
	user.ID = userResponse.ID

	return
}
