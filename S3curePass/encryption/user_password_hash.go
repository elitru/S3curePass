package encryption

import (
	"S3curePass/errors"
	"golang.org/x/crypto/bcrypt"
)

//hashes a given password string
func GetHash(password string) string {
	//hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	//check for fatal error
	errors.CheckFatal(err);
	return string(hash)
}

//checks whether two passwords are the same
func CheckPassword(password , hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
