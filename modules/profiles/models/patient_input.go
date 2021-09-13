package models

import "github.com/99designs/gqlgen/graphql"

// UpsertPatientInput is the schema for information needed to create or update a patient
type UpsertPatientInput struct {
	UserID    string          `json:"userId" bson:"userId"`
	FullName  string          `json:"fullName" bson:"fullName"`
	LikeName  string          `json:"likeName" bson:"likeName"`
	BirthDate int64           `json:"birthDate" bson:"birthDate"`
	City      string          `json:"city" bson:"city"`
	Avatar    *graphql.Upload `json:"avatar" bson:"avatar"`
}
