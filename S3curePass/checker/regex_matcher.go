package checker

import (
	"S3curePass/logger"
	"regexp"
)

//international name regex
const NAME_REGEX = `^[a-zA-Z]+(([',. -][a-zA-Z ])?[a-zA-Z]*)*$`
//Minimum eight characters, at least one letter, one number and one special character
const PASSWORD_REGEX = `^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]{8,}$`
//Alphanumeric string that may include _ and – having a length of 3 to 16 characters
const USERNAME_REGEX = `^[a-z0-9_-]{3,16}$`
//global email regex
const EMAIL_REGEX = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

//matches a string with a given regex
func Check(param string, regex string) bool {
	match, err := regexp.MatchString(regex, param)

	if err != nil {
		logger.Error(err.Error())
		return false
	}

	return match
}