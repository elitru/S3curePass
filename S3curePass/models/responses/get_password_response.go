package responses

type PasswordResponse struct {
	Passwords []PasswordResponseEntry `json:"passwords"`
}

type PasswordResponseEntry struct {
	Password    string    `json:"password"`
	UseLocation string    `json:"use_location"`
}