package responses

import "time"

type LoginResponse struct {
	Token      string    `json:"token"`
	ValidUntil time.Time `json:"valid_until"`
}