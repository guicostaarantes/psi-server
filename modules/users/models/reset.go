package users_models

// ResetPassword is the schema for a reset password request in the database
type ResetPassword struct {
	ID        string `json:"id" bson:"id"`
	UserID    string `json:"userId" bson:"userId"`
	Token     string `json:"token" bson:"token"`
	IssuedAt  int64  `json:"issuedAt" bson:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt" bson:"expiresAt"`
}

// ResetPasswordInput is the schema for information needed to reset a password
type ResetPasswordInput struct {
	Token    string `json:"token" bson:"token"`
	Password string `json:"password" bson:"password"`
}
