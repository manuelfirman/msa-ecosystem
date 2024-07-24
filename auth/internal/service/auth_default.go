package service

import (
	"auth/internal/domain"
	"auth/internal/repository"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

var (
	userServiceURL = "http://user_service:5001/users"
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
	ID int `json:"id"`
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

func (ad *AuthDefault) Register(user *domain.User) (err error) {

	err = createUserService(user)
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

func createUserService(user *domain.User) (err error) {
	// Crear un ID único para el usuario
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

	resp, err := http.Post(userServiceURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return ErrServiceFailedToCreateUser
	}

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

	// Asignar el ID al usuario
	user.ID = userResponse.ID

	return
}
