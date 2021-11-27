package users_models

import "time"

// ResetPassword is the schema for a reset password request in the database
type ResetPassword struct {
	UserID    string    `json:"userId" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"index"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	Redeemed  bool      `json:"redeemed"`
}
