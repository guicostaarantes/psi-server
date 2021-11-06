package models

// Agreement is an evidence that a term was agreed or disagreed by a profile
type Agreement struct {
	ID          string `json:"id"`
	TermName    string `json:"termName"`
	TermVersion int64  `json:"termVersion"`
	ProfileID   string `json:"profileId"`
	SignedAt    int64  `json:"signedAt"`
}
