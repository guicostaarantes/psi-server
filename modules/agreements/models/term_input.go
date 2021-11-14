package agreements_models

// UpsertTermInput is the schema for information needed to upsert a term
type UpsertTermInput struct {
	Name        string          `json:"name"`
	Version     int64           `json:"version"`
	ProfileType TermProfileType `json:"profileType"`
	Active      bool            `json:"active"`
}
