package models

import (
	"github.com/google/uuid"
	"time"
)

type Password struct {
	PasswordID  uuid.UUID `json:"password_id"`
	Password    string    `json:"password"`
	Nonce       string    `json:"nonce"`
	UseLocation string    `json:"use_location"`
	CreatedOn   time.Time `json:"created_on"`
	UserID      uuid.UUID `json:"user_id"`
}
