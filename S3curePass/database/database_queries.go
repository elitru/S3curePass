package database

import (
	"S3curePass/errors"
	"io/ioutil"
)

const PATH_TO_QUERIES = "./database/queries/"

func getQuery(queryFilename string) string {
	query, err := ioutil.ReadFile(PATH_TO_QUERIES + queryFilename + ".sql")

	errors.CheckFatal(err)

	return string(query)
}