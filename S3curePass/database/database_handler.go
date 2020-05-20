package database

import (
	"S3curePass/config"
	"S3curePass/errors"
	"S3curePass/logger"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var connection *sql.DB = nil

//opens a connection to the database
func Connect() {
	if connection != nil {
		return
	}

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetConfig().Database.Host,
		config.GetConfig().Database.User,
		config.GetConfig().Database.Password,
		config.GetConfig().Database.Database,
		config.GetConfig().Database.Port)

	conn, err := sql.Open("postgres", connectionString)
	errors.CheckFatal(err)
	connection = conn
	logger.Log("Database connection established successfully")
}

//closes the connection to the database
func Disconnect() {
	if connection != nil {
		err := connection.Close()
		errors.CheckFatal(err)
		logger.Log("Database connection closed successfully")
	}
}
