package requests

type GetPasswordRequest struct {
	UserPassword string `json:"user_password"`
	UseLocation  string `json:"use_location"`
}
