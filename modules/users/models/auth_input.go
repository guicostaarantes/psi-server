package users_models

// AuthenticateUserInput is the schema for information needed to authenticate a user
type AuthenticateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ValidateUserTokenInput is the schema for information needed to validate a user's token
type ValidateUserTokenInput struct {
	Token string `json:"token"`
}
