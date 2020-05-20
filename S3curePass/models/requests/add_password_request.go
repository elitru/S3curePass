package requests

type AddPasswordRequest struct {
	UserPassword  string `json:"user_password"`
	PasswordToAdd string `json:"password_to_add"`
	UseLocation   string `json:"use_location"`
}
