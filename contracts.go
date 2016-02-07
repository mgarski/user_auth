package main
import (
	"io"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

// used when registering a new user
type RegistrationRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func DecodeRegistrationRequest(reader io.Reader) (RegistrationRequest, error) {
	decoder := json.NewDecoder(reader)
	var req RegistrationRequest
	err := decoder.Decode(&req)
	if(err == nil) {
		empty := make([]string, 0, 3)
		if len(req.Name) == 0 {
			empty = append(empty, "name")
		}
		if len(req.Password) == 0 {
			empty = append(empty, "password")
		}
		if len(req.Email) == 0 {
			empty = append(empty, "email")
		}
		if len(empty) > 0 {
			err = errors.New(fmt.Sprintf("Fields missing from request: %v", empty))
		} else if !validateEmail(req.Email) {
			err = errors.New(fmt.Sprintf("Email [%s] is not valid", req.Email))
		}
	}
	return req, err
}

// used when updating an existing user
type UpdateRequest struct {
	RegistrationRequest
	Id int `json:"id"`
}

func DecodeUpdateRequest(reader io.Reader) (UpdateRequest, error) {
	decoder := json.NewDecoder(reader)
	var req UpdateRequest
	err := decoder.Decode(&req)
	if(err == nil) {
		empty := make([]string, 0, 4)
		if req.Id < 1 {
			empty = append(empty, "id")
		}
		if len(req.Name) == 0 {
			empty = append(empty, "name")
		}
		if len(req.Password) == 0 {
			empty = append(empty, "password")
		}
		if len(req.Email) == 0 {
			empty = append(empty, "email")
		}
		if len(empty) > 0 {
			err = errors.New(fmt.Sprintf("Fields missing from request: %v", empty))
		} else if !validateEmail(req.Email) {
			err = errors.New(fmt.Sprintf("Email [%s] is not valid", req.Email))
		}
	}
	return req, err
}

// used when authenticating a user
type AuthenticationRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func DecodeAuthenticationRequest(reader io.Reader) (AuthenticationRequest, error) {
	decoder := json.NewDecoder(reader)
	var req AuthenticationRequest
	err := decoder.Decode(&req)
	if(err == nil) {
		empty := make([]string, 0, 2)
		if len(req.Email) == 0 {
			empty = append(empty, "email")
		}
		if len(req.Password) == 0 {
			empty = append(empty, "password")
		}
		if len(empty) > 0 {
			err = errors.New(fmt.Sprintf("Fields missing from request: %v", empty))
		} else if !validateEmail(req.Email) {
			err = errors.New(fmt.Sprintf("Email [%s] is not valid", req.Email))
		}
	}
	return req, err
}

// used when validating the token
type TokenRequest struct {
	Token string `json:"token"`
}

func DecodeTokenRequest(reader io.Reader) (TokenRequest, error) {
	decoder := json.NewDecoder(reader)
	var req TokenRequest
	err := decoder.Decode(&req)
	if(err == nil) {
		if len(req.Token) == 0 {
			err = errors.New("Token field missing from request")
		}
	}
	return req, err
}

// used when validating the token or logging out the user
type DeleteRequest struct {
	Id int `json:"id"`
}

func DecodeDeleteRequest(reader io.Reader) (DeleteRequest, error) {
	decoder := json.NewDecoder(reader)
	var req DeleteRequest
	err := decoder.Decode(&req)
	if(err == nil) {
		if req.Id < 1 {
			err = errors.New("ID field missing from request")
		}
	}
	return req, err
}

// basic response with an error code and message
// returned for registration, update, delete, log out & validation
type BasicResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

// response used for log in
type AuthenticationResponse struct {
	Token string `json:"token,omitempty"`
	BasicResponse
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
