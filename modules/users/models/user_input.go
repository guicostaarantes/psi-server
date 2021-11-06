package models

// CreateUserInput is the schema for information needed to create a user
type CreateUserInput struct {
	Email string `json:"email"`
	Role  Role   `json:"role"`
}

// CreateUserWithPasswordInput is the schema for information needed to create a user
type CreateUserWithPasswordInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

// UpdateUserInput is the schema for information needed to update a user
type UpdateUserInput struct {
	Active bool `json:"active"`
	Role   Role `json:"role"`
}
