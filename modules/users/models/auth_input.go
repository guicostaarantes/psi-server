package models

// AuthenticateUserInput is the schema for information needed to authenticate a user
type AuthenticateUserInput struct {
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	IPAddress string `json:"ipAddress" bson:"ipAddress"`
}

// ValidateUserTokenInput is the schema for information needed to validate a user's token
type ValidateUserTokenInput struct {
	Token string `json:"token" bson:"token"`
}
