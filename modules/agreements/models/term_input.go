package models

// UpsertTermInput is the schema for information needed to upsert a term
type UpsertTermInput struct {
	Name        string          `json:"name" bson:"name"`
	Version     int64           `json:"version" bson:"version"`
	ProfileType TermProfileType `json:"profileType" bson:"profileType"`
	Active      bool            `json:"active" bson:"active"`
}
