package requests

type UpdatePasswordRequest struct {
	UserPassword string `json:"user_password"`
	PasswordID   string `json:"password_id"`
	Password     string `json:"password"`
	UseLocation  string `json:"use_location"`
}
