package users_models

// ResetPasswordInput is the schema for information needed to reset a password
type ResetPasswordInput struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}
