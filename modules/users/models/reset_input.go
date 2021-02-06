package models

// ResetPasswordInput is the schema for information needed to reset a password
type ResetPasswordInput struct {
	Token    string `json:"token" bson:"token"`
	Password string `json:"password" bson:"password"`
}
