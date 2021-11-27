package users_models

import "time"

// Authentication is the schema for saving successful authentication attempts for a user
type Authentication struct {
	UserID    string    `json:"userId" gorm:"primaryKey"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	Token     string    `json:"token" gorm:"index"`
}
