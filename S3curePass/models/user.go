package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserID       uuid.UUID `json:"user_id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredOn time.Time `json:"registered_on"`
}
