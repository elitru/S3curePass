package auth

import (
	"S3curePass/checker"
	"S3curePass/config"
	"S3curePass/database"
	"S3curePass/encryption"
	"S3curePass/errors"
	"S3curePass/messages"
	"S3curePass/models/requests"
	"S3curePass/models/responses"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//logs a user in by sending him his authentication token
func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest requests.LoginRequest

	//decode request body
	err := json.NewDecoder(r.Body).Decode(&loginRequest)

	//check if body was valid
	if err != nil || loginRequest.Password == "" || loginRequest.Username == "" {
		http.Error(w, messages.BAD_REQUEST, http.StatusBadRequest)
		return
	}

	//get user from database to check password
	found, user := database.GetUserByUsername(loginRequest.Username)

	//check whether a user with the given username even exists
	if !found {
		http.Error(w, messages.INVALID_LOGIN, http.StatusUnauthorized)
		return
	}

	//check is password is correct
	if !encryption.CheckPassword(loginRequest.Password, user.Password) {
		http.Error(w, messages.INVALID_LOGIN, http.StatusUnauthorized)
		return
	}

	//generate the authentication token
	token, validUntil := generateToken(user.UserID.String())

	//return result
	json.NewEncoder(w).Encode(responses.LoginResponse{Token: token, ValidUntil: validUntil})
}

//creates a new account
func Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest requests.RegisterRequest

	//decode request body
	err := json.NewDecoder(r.Body).Decode(&registerRequest)

	//check if body was valid
	if err != nil {
		http.Error(w, messages.BAD_REQUEST, http.StatusBadRequest)
		return
	}

	//get user from database to whether a user with that username already exists
	found, _ := database.GetUserByUsername(registerRequest.Username)

	if found {
		//a user with that username already exists
		http.Error(w, messages.USERNAME_ALREADY_TAKEN, http.StatusBadRequest)
		return
	}

	//check for duplicate email
	//get user from database to whether a user with that email already exists
	found, _ = database.GetUserByEmail(registerRequest.Email)

	if found {
		//a user with that email already exists
		http.Error(w, messages.EMAIL_ALREADY_TAKEN, http.StatusBadRequest)
		return
	}

	//check remaining user data
	if !checker.Check(registerRequest.Firstname, checker.NAME_REGEX) || !checker.Check(registerRequest.Lastname, checker.NAME_REGEX) {
		http.Error(w, messages.NAME_INVALID, http.StatusBadRequest)
		return
	}

	if !checker.Check(registerRequest.Email, checker.EMAIL_REGEX) {
		http.Error(w, messages.EMAIL_INVALID, http.StatusBadRequest)
		return
	}

	if !checker.Check(registerRequest.Password, checker.PASSWORD_REGEX) {
		http.Error(w, messages.PASSWORD_INVALID, http.StatusBadRequest)
		return
	}

	if !checker.Check(registerRequest.Username, checker.USERNAME_REGEX) {
		http.Error(w, messages.USERNAME_ALREADY_TAKEN, http.StatusBadRequest)
		return
	}

	//everything ok -> create user in database
	userID := database.CreateUser(&registerRequest)

	//generate the authentication token
	token, validUntil := generateToken(userID)

	//return result
	json.NewEncoder(w).Encode(responses.LoginResponse{Token: token, ValidUntil: validUntil})
}

//returns the user id to a given request header containing the jwt authentication token
func GetUserID(header *http.Header) string {
	userID, _ := DecodeToken(header.Get(config.GetConfig().Token.Header))
	return userID
}

//generates the authentication jwt token for a user
func generateToken(userID string) (string, time.Time) {
	os.Setenv("ACCESS_SECRET", config.GetConfig().Token.Secret)

	exp := time.Now().Add(time.Hour * 24)

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = exp.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	errors.CheckFatal(err)

	return token, exp
}

//decodes the jwt token into it parts (userID)
func DecodeToken(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Token.Secret), nil
	})

	if err != nil {
		return "", err
	}

	return claims["user_id"].(string), nil
}
