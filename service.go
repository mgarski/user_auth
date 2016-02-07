package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"os"
)

var config Configuration

type Configuration struct {
	DbConnection string `json:"dbConnection"`
}

func main() {
	// load up the database config
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("configuration error:", err)
		return
	}

	http.HandleFunc("/login", LogIn)
	http.HandleFunc("/logout", LogOut)
	http.HandleFunc("/user", UserManagement)
	http.HandleFunc("/validate", Validate)
	http.ListenAndServe(":8080", nil)
}

func UserManagement(w http.ResponseWriter, r *http.Request) {
	var response BasicResponse

	switch {
	// create
	case r.Method == "PUT":
		response = registerUser(r)
		break
	// update
	case r.Method == "POST":
		response = updateUser(r)
		break
	// delete
	case r.Method == "DELETE":
		response = deleteUser(r)
		break;
	}

	w.WriteHeader(response.Code)
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

// Logs a user in
func LogIn(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeAuthenticationRequest(r.Body)
	if err != nil {
		handleRequestError(w, err)
		return
	}

	user := User{Email:req.Email, Password:req.Password}
	userId := user.ValidateCredentials()
	fmt.Printf("user id: %d\n", userId)
	var response AuthenticationResponse
	if(userId > 0) {
		response.Token = GenerateToken(userId)
		response.Code = 200
		response.Message = "OK"
	} else {
		response.Code = 401
		response.Message = "Authentication Failed"
		w.WriteHeader(401)
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

// logs a user out by invalidating any tokens
func LogOut(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeDeleteRequest(r.Body)
	if err != nil {
		handleRequestError(w, err)
		return
	}

	if !FlushToken(req.Id) {
		w.WriteHeader(400)
	}
}

// validates a token
func Validate(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeTokenRequest(r.Body)
	if err != nil {

	}
	if !ValidateToken(req.Token) {
		w.WriteHeader(404)
	}
}

// creates a new user account
func registerUser(r *http.Request) BasicResponse {
	req, err := DecodeRegistrationRequest(r.Body)
	if err != nil {
		return BasicResponse{Message:err.Error(), Code:400}
	}

	u := User{Name:req.Name, Email:req.Email, Password:req.Password}
	code, message := u.Register()
	return BasicResponse{code, message}
}

// performs work of updating a user account
func updateUser(r *http.Request) BasicResponse {
	req, err := DecodeUpdateRequest(r.Body)
	if err != nil {
		return BasicResponse{Message:err.Error(), Code:400}
	}

	u := User{Id:req.Id, Name:req.Name, Email:req.Email, Password:req.Password}
	code, message := u.Update()
	return BasicResponse{code, message}
}

// performs work of deleting a user account
func deleteUser(r *http.Request) BasicResponse {
	req, err := DecodeDeleteRequest(r.Body)
	if err != nil {
		return BasicResponse{Message:err.Error(), Code:400}
	}

	u := User{Id:req.Id}
	code, message := u.Delete()
	return BasicResponse{code, message}
}

// helper for handling a request error
func handleRequestError(w http.ResponseWriter, err error) {
	encoder := json.NewEncoder(w)
	response := BasicResponse{Message:err.Error(), Code:400}
	w.WriteHeader(400)
	encoder.Encode(response)
}
