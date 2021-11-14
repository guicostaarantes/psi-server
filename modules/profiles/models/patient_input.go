package profiles_models

import "github.com/99designs/gqlgen/graphql"

// UpsertPatientInput is the schema for information needed to create or update a patient
type UpsertPatientInput struct {
	UserID    string          `json:"userId"`
	FullName  string          `json:"fullName"`
	LikeName  string          `json:"likeName"`
	BirthDate int64           `json:"birthDate"`
	City      string          `json:"city"`
	Avatar    *graphql.Upload `json:"avatar"`
}
