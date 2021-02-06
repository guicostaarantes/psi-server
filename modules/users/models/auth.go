package models

// Authentication is the schema for saving successful authentication attempts for a user
type Authentication struct {
	UserID    string `json:"userId" bson:"userId"`
	IPAddress string `json:"ipAddress" bson:"ipAddress"`
	IssuedAt  int64  `json:"issuedAt" bson:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt" bson:"expiresAt"`
	Token     string `json:"token" bson:"token"`
}
