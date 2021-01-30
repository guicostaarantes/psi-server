package users_models

type ResetPassword struct {
	ID        string `json:"id" bson:"id"`
	UserID    string `json:"userId" bson:"userId"`
	Token     string `json:"token" bson:"token"`
	IssuedAt  int64  `json:"issuedAt" bson:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt" bson:"expiresAt"`
}
