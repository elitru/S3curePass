package database

import (
	"S3curePass/encryption"
	"S3curePass/errors"
	"S3curePass/models"
	"S3curePass/models/requests"
	"database/sql"
)

const (
	GET_USER_BY_USERNAME       = "get_user_by_username"
	GET_USER_BY_EMAIL          = "get_user_by_email"
	INSERT_USER                = "insert_user"
	GET_ALL_PASSWORDS_FOR_USER = "get_all_passwords_for_user"
	GET_HASH_FOR_USER_ID       = "get_hash_for_user_id"
	FIND_PASSWORDS_BY_LOCATION = "find_password_by_location"
	INSERT_PASSWORD            = "insert_password"
	DELETE_PASSWORD            = "delete_password"
	UPDATE_PASSWORD            = "update_password_for_user"
)

//get a user from the database according to a given username
//returns -> bool: whether an entry was found
//			 models.User: is the user object if a result was found
func GetUserByUsername(username string) (bool, models.User) {
	//get query from file
	sqlQuery := getQuery(GET_USER_BY_USERNAME)
	//execute statement on database
	result := connection.QueryRow(sqlQuery, username)
	user := models.User{}
	//parse result
	err := result.Scan(
		&user.UserID,
		&user.Firstname,
		&user.Lastname,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.RegisteredOn)

	//handle errors
	switch err {
	case sql.ErrNoRows:
		return false, models.User{}
	default:
		errors.CheckFatal(err)
	}

	return true, user
}

//get a user from the database according to a given email
//returns -> bool: whether an entry was found
//			 models.User: is the user object if a result was found
func GetUserByEmail(email string) (bool, models.User) {
	//get query from file
	sqlQuery := getQuery(GET_USER_BY_EMAIL)
	//execute statement on database
	result := connection.QueryRow(sqlQuery, email)
	user := models.User{}
	//parse result
	err := result.Scan(
		&user.UserID,
		&user.Firstname,
		&user.Lastname,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.RegisteredOn)

	//handle errors
	switch err {
	case sql.ErrNoRows:
		return false, models.User{}
	default:
		errors.CheckFatal(err)
	}

	return true, user
}

//create user in database
//returns userID of created user
func CreateUser(req *requests.RegisterRequest) string {
	//get query from file
	sqlQuery := getQuery(INSERT_USER)

	//execute insert statement
	var userID string
	err := connection.QueryRow(sqlQuery,
		req.Firstname,
		req.Lastname,
		req.Username,
		req.Email,
		encryption.GetHash(req.Password)).Scan(&userID)

	//check errors
	errors.CheckFatal(err)

	return userID
}

//retrieves all passwords according to a given user
func GetAllPasswordsForUser(userID string) []models.Password {
	//get query from file
	sqlQuery := getQuery(GET_ALL_PASSWORDS_FOR_USER)

	rows, err := connection.Query(sqlQuery, userID)

	//handle errors
	switch err {
	case sql.ErrNoRows:
		return []models.Password{}
	default:
		errors.CheckFatal(err)
	}

	//parse to structs
	var passwords []models.Password

	for rows.Next() {
		//create "empty" password...
		var password models.Password = models.Password{}

		//... and populate object with values
		err = rows.Scan(
			&password.PasswordID,
			&password.Password,
			&password.Nonce,
			&password.UseLocation,
			&password.CreatedOn,
			&password.UserID)

		//check for errors
		errors.CheckFatal(err)

		//add parsed password to list
		passwords = append(passwords, password)
	}

	return passwords
}

//get the password-hash for a given user-id
func GetPasswordHashForUserId(userID string) string {
	//get query from file
	sqlQuery := getQuery(GET_HASH_FOR_USER_ID)

	var passwordHash string

	//gets one row and parses hit
	err := connection.QueryRow(sqlQuery, userID).Scan(&passwordHash)

	errors.CheckFatal(err)

	return passwordHash
}

//finds a password by a use location
func FindPasswordByLocation(userID, useLocation string) []models.Password {
	//get query from file
	sqlQuery := getQuery(FIND_PASSWORDS_BY_LOCATION)

	rows, err := connection.Query(sqlQuery, userID, useLocation)

	//handle errors
	switch err {
	case sql.ErrNoRows:
		return []models.Password{}
	default:
		errors.CheckFatal(err)
	}

	//parse to structs
	var passwords []models.Password

	for rows.Next() {
		//create "empty" password...
		var password models.Password = models.Password{}

		//... and populate object with values
		err = rows.Scan(
			&password.PasswordID,
			&password.Password,
			&password.Nonce,
			&password.UseLocation,
			&password.CreatedOn,
			&password.UserID)

		//check for errors
		errors.CheckFatal(err)

		//add parsed password to list
		passwords = append(passwords, password)
	}

	return passwords
}

//deletes a password for a user if it exists
//returns true if a password was found and delete otherwise it returns false
func DeletePassword(userID, passwordID string) bool {
	//get query from file
	sqlQuery := getQuery(DELETE_PASSWORD)
	//execute query
	res, err := connection.Exec(sqlQuery, passwordID, userID)

	//handle errors
	errors.CheckFatal(err)
	//check if somethings was deleted
	count, err := res.RowsAffected()
	errors.CheckFatal(err)

	if count == 0 {
		return false
	}

	return true
}

// Inserts a new password in the database, given a password struct
func CreatePassword(userID string, password *models.Password) {
	// get query from file
	sqlQuery := getQuery(INSERT_PASSWORD)

	//execute query
	_, err := connection.Query(sqlQuery, password.Password, password.Nonce, password.UseLocation, password.CreatedOn, userID)

	// handle errors
	errors.CheckFatal(err)
}

//changes the details of a password for a user
func UpdatePassword(userID string, password *models.Password) {
	// get query from file
	sqlQuery := getQuery(UPDATE_PASSWORD)

	//execute query and save password to database
	_, err := connection.Exec(sqlQuery, password.Password, password.Nonce, password.UseLocation, password.CreatedOn, userID, password.PasswordID.String())

	//handle errors
	errors.CheckFatal(err)
}
