package models

// Authentication is the schema for saving successful authentication attempts for a user
type Authentication struct {
	UserID    string `json:"userId" bson:"userId"`
	IPAddress string `json:"ipAddress" bson:"ipAddress"`
	IssuedAt  int64  `json:"issuedAt" bson:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt" bson:"expiresAt"`
	Token     string `json:"token" bson:"token"`
}

// BadAuthentication is the schema for saving UNsuccessful authentication attempts for a user
type BadAuthentication struct {
	UserID    string `json:"userId" bson:"userId"`
	IPAddress string `json:"ipAddress" bson:"ipAddress"`
	IssuedAt  int64  `json:"issuedAt" bson:"issuedAt"`
}

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
