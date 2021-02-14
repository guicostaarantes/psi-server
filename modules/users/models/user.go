package models

// Role represents the level of authorization of the user
type Role string

const (
	// Coordinator is responsible for creating and managing psychologists
	Coordinator Role = "COORDINATOR"
	// Psychologist is responsible for treating patients
	Psychologist Role = "PSYCHOLOGIST"
	// Patient is seeking psychological assistance
	Patient Role = "PATIENT"
)

// User is the schema for a user in the database
type User struct {
	ID       string `json:"id" bson:"id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Active   bool   `json:"active" bson:"active"`
	Role     Role   `json:"role" bson:"role"`
}
