package encryption

import (
	"S3curePass/config"
	"S3curePass/errors"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

//returns the aes-256 encrypted password which will be stored
//the userPassword and the creation time will be needed for the secret key
//as the user password word must have at least 8 characters, we can use them
func ToSecuredPassword(passwordToBeSecured, userPassword string, created time.Time) (passwordEncrypted string, nonce string) {
	//generate secret key
	key := userPassword[0:8] + config.GetConfig().Encryption.Secret + created.Format("20060102150405")
	secretKey := []byte(key)

	block, err := aes.NewCipher(secretKey)
	//check for errors
	errors.CheckFatal(err)

	//create nonce
	nonceBytes := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonceBytes); err != nil {
		errors.CheckFatal(err)
	}

	aesgcm, err := cipher.NewGCM(block)
	//check for errors
	errors.CheckFatal(err)

	//encrypt password
	encryptedPassword := aesgcm.Seal(nil, nonceBytes, []byte(passwordToBeSecured), nil)
	return fmt.Sprintf("%x", encryptedPassword), fmt.Sprintf("%x", nonceBytes)
}

//returns the plain password from a given aes-256 encrypted password
func ToPlainPassword(encryptedPassword, nonce, userPassword string, created time.Time) string {
	//generate secret key
	key := userPassword[0:8] + config.GetConfig().Encryption.Secret + created.Format("20060102150405")
	secretKey := []byte(key)

	cipherText, err := hex.DecodeString(encryptedPassword)
	//check for errors
	errors.CheckFatal(err)

	nonceBytes, err := hex.DecodeString(nonce)
	//check for errors
	errors.CheckFatal(err)

	block, err := aes.NewCipher(secretKey)
	//check for errors
	errors.CheckFatal(err)

	aesgcm, err := cipher.NewGCM(block)
	//check for errors
	errors.CheckFatal(err)

	//decrypt password
	passwordPlain, err := aesgcm.Open(nil, nonceBytes, cipherText, nil)
	//check for errors
	errors.CheckFatal(err)

	return string(passwordPlain)
}