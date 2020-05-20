package errors

import "S3curePass/logger"

func CheckFatal(err error) {
	if err != nil {
		logger.Error(err.Error())
		panic("")
	}
}

func CheckError(err error) {
	if err != nil {
		logger.Fatal(err.Error())
	}
}