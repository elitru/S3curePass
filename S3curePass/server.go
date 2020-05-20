package main

import (
	"S3curePass/database"
	"S3curePass/handlers/auth"
	"S3curePass/handlers/middleware"
	"S3curePass/handlers/passwords"
	"S3curePass/logger"
	"net/http"

	"github.com/gorilla/mux"
)

//method which is called to start the web service
func main() {
	router := Router()

	logger.Log("Web service running on port " + PORT)
	//start web service
	err := http.ListenAndServe(":"+PORT, router)

	//check if an error occured during starting proccess
	if err != nil {
		logger.Error("Error starting web servcie -> \n" + err.Error())
	}
}

const (
	//path to the public folder
	PUBLIC = "./public/"
	//path prefix for public folder
	PUBLIC_PREFIX = "/"
	//port on which the web service will be running
	PORT = "80"
)

//returns a new mux router, with all routes defined
func Router() *mux.Router {
	router := mux.NewRouter()

	//define routes
	router.HandleFunc("/auth/login", auth.Login).Methods("POST")
	router.HandleFunc("/auth/register", auth.Register).Methods("POST")
	router.HandleFunc("/passwords/all", middleware.WithAuthentication(passwords.GetAllPasswords).ServeHTTP).Methods("POST")
	router.HandleFunc("/passwords/find", middleware.WithAuthentication(passwords.GetPasswordsByUseLocation).ServeHTTP).Methods("POST")
	router.HandleFunc("/passwords/add", middleware.WithAuthentication(passwords.AddPassword).ServeHTTP).Methods("POST")
	router.HandleFunc("/passwords/delete", middleware.WithAuthentication(passwords.DeletePassword).ServeHTTP).Methods("POST")
	router.HandleFunc("/passwords/update", middleware.WithAuthentication(passwords.UpdatePassword).ServeHTTP).Methods("POST")

	//connect to database
	database.Connect()

	//define static directory
	router.
		PathPrefix(PUBLIC_PREFIX).
		Handler(http.StripPrefix(PUBLIC_PREFIX, http.FileServer(http.Dir(PUBLIC))))

	return router
}

/*
	-> Test user
	{
		"firstname": "Max",
		"lastname": "Mustermann",
		"email": "max.mustermann@mail.com",
		"username": "admin1",
		"password": "Admin#1999"
	}
*/
