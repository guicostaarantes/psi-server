package models

// ResetPassword is the schema for a reset password request in the database
type ResetPassword struct {
	UserID    string `json:"userId"`
	Token     string `json:"token" gorm:"primaryKey"`
	IssuedAt  int64  `json:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt"`
	Redeemed  bool   `json:"redeemed"`
}
