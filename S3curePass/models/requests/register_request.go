package requests

type RegisterRequest struct {
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
}