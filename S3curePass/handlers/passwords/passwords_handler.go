package passwords

import (
	"S3curePass/checker"
	"S3curePass/database"
	"S3curePass/encryption"
	"S3curePass/messages"
	"S3curePass/models"
	"S3curePass/models/requests"
	"S3curePass/models/responses"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

//returns all password for a user
func GetAllPasswords(w http.ResponseWriter, r *http.Request, userID string) {
	var getPasswordRequest requests.GetPasswordRequest
	//decode request body
	err := json.NewDecoder(r.Body).Decode(&getPasswordRequest)

	//check if body was valid
	if err != nil {
		http.Error(w, messages.BAD_REQUEST, http.StatusBadRequest)
		return
	}

	//get password(-hash) of current user
	currentUserHash := database.GetPasswordHashForUserId(userID)

	// make sure that the (user-)password which the user provided is valid
	encryption.CheckPassword(getPasswordRequest.UserPassword, currentUserHash)

	//password was valid
	if !checker.Check(getPasswordRequest.UserPassword, checker.PASSWORD_REGEX) {
		// password was invalid --> return error
		http.Error(w, messages.UNAUTHORIZED, http.StatusBadRequest)
		return
	}

	//get all encrypted passwords from database
	passwords := database.GetAllPasswordsForUser(userID)

	//decrypt all passwords
	for i := range passwords {
		passwords[i].Password = encryption.ToPlainPassword(passwords[i].Password, passwords[i].Nonce, getPasswordRequest.UserPassword, passwords[i].CreatedOn)
	}

	//return passwords
	response := responses.PasswordResponse{}

	//map to response
	for _, password := range passwords {
		passwordEntry := responses.PasswordResponseEntry{
			Password:    password.Password,
			UseLocation: password.UseLocation,
		}

		response.Passwords = append(response.Passwords, passwordEntry)
	}

	//encode result to json and send back to user
	json.NewEncoder(w).Encode(response)
}

//returns all passwords where a given use location is contained in the passwords use location
func GetPasswordsByUseLocation(w http.ResponseWriter, r *http.Request, userID string) {
	var getPasswordRequest requests.GetPasswordRequest
	//decode request body
	err := json.NewDecoder(r.Body).Decode(&getPasswordRequest)

	//check if body was valid
	if err != nil {
		http.Error(w, messages.BAD_REQUEST, http.StatusBadRequest)
		return
	}

	//get password(-hash) of current user
	currentUserHash := database.GetPasswordHashForUserId(userID)

	// make sure that the (user-)password which the user provided is valid
	encryption.CheckPassword(getPasswordRequest.UserPassword, currentUserHash)

	//password was valid
	if !checker.Check(getPasswordRequest.UserPassword, checker.PASSWORD_REGEX) {
		// password was invalid --> return error
		http.Error(w, messages.PASSWORD_INVALID, http.StatusBadRequest)
		return
	}

	//get all encrypted passwords from database
	passwords := database.FindPasswordByLocation(userID, getPasswordRequest.UseLocation)

	//decrypt all passwords
	for i := range passwords {
		passwords[i].Password = encryption.ToPlainPassword(passwords[i].Password, passwords[i].Nonce, getPasswordRequest.UserPassword, passwords[i].CreatedOn)
	}

	//return passwords
	response := responses.PasswordResponse{}

	//map to response
	//only add passwords where the use location contains the given use location
	for _, password := range passwords {
		passwordEntry := responses.PasswordResponseEntry{
			Password:    password.Password,
			UseLocation: password.UseLocation,
		}

		response.Passwords = append(response.Passwords, passwordEntry)
	}

	//encode result to json and send back to user
	json.NewEncoder(w).Encode(response)
}

//save a new password entry for a given user into the database
func AddPassword(w http.ResponseWriter, r *http.Request, userID string) {
	var createPasswordRequest requests.AddPasswordRequest
	//decode request body
	err := json.NewDecoder(r.Body).Decode(&createPasswordRequest)

	//check if body was valid
	if err != nil {
		http.Error(w, messages.BAD_REQUEST, http.StatusBadRequest)
		return
	}

	//get password(-hash) of current user
	currentUserHash := database.GetPasswordHashForUserId(userID)

	// make sure that the (user-)password which the user provided is valid
	encryption.CheckPassword(createPasswordRequest.UserPassword, currentUserHash)

	//password was valid
	if !checker.Check(createPasswordRequest.PasswordToAdd, checker.PASSWORD_REGEX) {
		// password was invalid --> return error
		http.Error(w, messages.PASSWORD_INVALID, http.StatusBadRequest)
		return
	}

	//get current time
	currentTime := time.Now()

	//encrypt the new password
	encryptedPassword, nonce := encryption.ToSecuredPassword(createPasswordRequest.PasswordToAdd, createPasswordRequest.UserPassword, currentTime)

	newPassword := models.Password{
		Password:    encryptedPassword,
		Nonce:       nonce,
		UseLocation: createPasswordRequest.UseLocation,
		CreatedOn:   currentTime,
	}

	// ... and insert it in the db!
	database.CreatePassword(userID, &newPassword)

	//return ok
	w.WriteHeader(200)
}

//deletes a given password entry for a user
func DeletePassword(w http.ResponseWriter, r *http.Request, userID string) {
	var deleteRequest requests.DeletePasswordRequest
	//parse body
	err := json.NewDecoder(r.Body).Decode(&deleteRequest)

	//check if body was valid
	if err != nil {
		http.Error(w, messages.BAD_REQUEST, http.StatusBadRequest)
		return
	}

	//delete password from database
	deleted := database.DeletePassword(userID, deleteRequest.PasswordID)

	//check if password has been deleted
	if !deleted {
		//password couldn't be deleted because it either didn't exist or a database error occured
		http.Error(w, messages.PASSWORD_NOT_DELETED, http.StatusBadRequest)
		return
	}

	//password was successfully deleted -> return ok
	w.WriteHeader(200)
}

//edits a given password for a user
func UpdatePassword(w http.ResponseWriter, r *http.Request, userID string) {
	//parse body
	var updateRequest requests.UpdatePasswordRequest
	err := json.NewDecoder(r.Body).Decode(&updateRequest)

	//check if body was valid
	if err != nil {
		http.Error(w, messages.BAD_REQUEST, http.StatusBadRequest)
		return
	}

	//get password(-hash) of current user
	currentUserHash := database.GetPasswordHashForUserId(userID)

	// make sure that the (user-)password which the user provided is valid
	encryption.CheckPassword(updateRequest.UserPassword, currentUserHash)

	//password was valid
	if !checker.Check(updateRequest.Password, checker.PASSWORD_REGEX) {
		// password was invalid --> return error
		http.Error(w, messages.PASSWORD_INVALID, http.StatusBadRequest)
		return
	}

	//get current time
	currentTime := time.Now()

	//encrypt the new password
	encryptedPassword, nonce := encryption.ToSecuredPassword(updateRequest.Password, updateRequest.UserPassword, currentTime)

	password := models.Password{
		PasswordID:  uuid.MustParse(updateRequest.PasswordID),
		Password:    encryptedPassword,
		Nonce:       nonce,
		UseLocation: updateRequest.UseLocation,
		CreatedOn:   currentTime,
	}

	// ... and update password entry in database
	database.UpdatePassword(userID, &password)

	//return ok
	w.WriteHeader(200)
}
