package models

// ResetPassword is the schema for a reset password request in the database
type ResetPassword struct {
	UserID    string `json:"userId" bson:"userId"`
	Token     string `json:"token" bson:"token"`
	IssuedAt  int64  `json:"issuedAt" bson:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt" bson:"expiresAt"`
	Redeemed  bool   `json:"redeemed" bson:"redeemed"`
}

// ResetPasswordInput is the schema for information needed to reset a password
type ResetPasswordInput struct {
	Token    string `json:"token" bson:"token"`
	Password string `json:"password" bson:"password"`
}
