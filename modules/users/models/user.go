package users_models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents the level of authorization of the user
type Role string

const (
	// JobRunner is reserved for robot users that will run jobs in the server
	JobRunner Role = "JOBRUNNER"
	// Coordinator is responsible for creating and managing psychologists
	Coordinator Role = "COORDINATOR"
	// Psychologist is responsible for treating patients
	Psychologist Role = "PSYCHOLOGIST"
	// Patient is seeking psychological assistance
	Patient Role = "PATIENT"
)

// User is the schema for a user in the database
type User struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt`
	UpdatedAt time.Time      `json:"updatedAt`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Email     string         `json:"email" gorm:"index"`
	Password  string         `json:"password"`
	Active    bool           `json:"active"`
	Role      Role           `json:"role"`
}
