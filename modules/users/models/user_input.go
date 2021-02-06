package models

// CreateUserInput is the schema for information needed to create a user
type CreateUserInput struct {
	Email     string `json:"email" bson:"email"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Role      Role   `json:"role" bson:"role"`
}

// CreateUserWithPasswordInput is the schema for information needed to create a user
type CreateUserWithPasswordInput struct {
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Role      Role   `json:"role" bson:"role"`
}

// UpdateUserInput is the schema for information needed to update a user
type UpdateUserInput struct {
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Role      Role   `json:"role" bson:"role"`
}
