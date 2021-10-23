package models

// Agreement is an evidence that a term was agreed or disagreed by a profile
type Agreement struct {
	ID          string `json:"id" bson:"id"`
	TermName    string `json:"termName" bson:"termName"`
	TermVersion int64  `json:"termVersion" bson:"termVersion"`
	ProfileID   string `json:"profileId" bson:"profileId"`
	SignedAt    int64  `json:"signedAt" bson:"signedAt"`
}
