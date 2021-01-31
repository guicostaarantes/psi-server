package users_models

// User is the schema for a user in the database
type User struct {
	ID        string `json:"id" bson:"id"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	Active    bool   `json:"active" bson:"active"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Role      string `json:"role" bson:"role"`
}

// CreateUserInput is the schema for information needed to create a user
type CreateUserInput struct {
	Email     string `json:"email" bson:"email"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Role      string `json:"role" bson:"role"`
}

// CreateUserWithPasswordInput is the schema for information needed to create a user
type CreateUserWithPasswordInput struct {
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Role      string `json:"role" bson:"role"`
}

// UpdateUserInput is the schema for information needed to update a user
type UpdateUserInput struct {
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Role      string `json:"role" bson:"role"`
}
