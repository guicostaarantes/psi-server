package models

import "github.com/99designs/gqlgen/graphql"

// UpsertPsychologistInput is the schema for information needed to create or update a psychologist
type UpsertPsychologistInput struct {
	UserID    string          `json:"userId" bson:"userId"`
	FullName  string          `json:"fullName" bson:"fullName"`
	LikeName  string          `json:"likeName" bson:"likeName"`
	BirthDate int64           `json:"birthDate" bson:"birthDate"`
	City      string          `json:"city" bson:"city"`
	Crp       string          `json:"crp" bson:"crp"`
	Whatsapp  string          `json:"whatsapp" bson:"whatsapp"`
	Instagram string          `json:"instagram" bson:"instagram"`
	Bio       string          `json:"bio" bson:"bio"`
	Avatar    *graphql.Upload `json:"avatar" bson:"avatar"`
}
