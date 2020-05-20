package middleware

import (
	"S3curePass/config"
	"S3curePass/handlers/auth"
	"S3curePass/messages"
	"net/http"
)

//middleware for checking whether a user is validly logged in and has permission to access a given route
func WithAuthentication(next func(w http.ResponseWriter, r *http.Request, userID string)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//read token from header
		token := r.Header.Get(config.GetConfig().Token.Header)

		if token == "" {
			http.Error(w, messages.UNAUTHORIZED, http.StatusUnauthorized)
			return
		}

		//check if token is valid
		userID, err := auth.DecodeToken(token)

		//check if token is valid
		if err != nil {
			http.Error(w, messages.UNAUTHORIZED, http.StatusUnauthorized)
			return
		}

		//call next handler
		next(w, r, userID)
	})
}
