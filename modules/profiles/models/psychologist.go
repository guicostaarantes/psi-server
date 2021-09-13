package models

// Psychologist is the schema for the profile of a psychologist
type Psychologist struct {
	ID        string `json:"id" bson:"id"`
	UserID    string `json:"userId" bson:"userId"`
	FullName  string `json:"fullName" bson:"fullName"`
	LikeName  string `json:"likeName" bson:"likeName"`
	BirthDate int64  `json:"birthDate" bson:"birthDate"`
	City      string `json:"city" bson:"city"`
	Bio       string `json:"bio" bson:"bio"`
	Avatar    string `json:"avatar" bson:"avatar"`
}
