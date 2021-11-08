package models

// Authentication is the schema for saving successful authentication attempts for a user
type Authentication struct {
	UserID    string `json:"userId" gorm:"primaryKey"`
	IssuedAt  int64  `json:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt"`
	Token     string `json:"token" gorm:"index"`
}
