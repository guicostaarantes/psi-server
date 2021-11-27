package profiles_models

import (
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// UpsertPsychologistInput is the schema for information needed to create or update a psychologist
type UpsertPsychologistInput struct {
	UserID    string          `json:"userId"`
	FullName  string          `json:"fullName"`
	LikeName  string          `json:"likeName"`
	BirthDate time.Time       `json:"birthDate"`
	City      string          `json:"city"`
	Crp       string          `json:"crp"`
	Whatsapp  string          `json:"whatsapp"`
	Instagram string          `json:"instagram"`
	Bio       string          `json:"bio"`
	Avatar    *graphql.Upload `json:"avatar"`
}
