package models

// UpsertAgreementInput is the schema for information needed to upsert an agreement
type UpsertAgreementInput struct {
	TermName    string `json:"termName" bson:"termName"`
	TermVersion int64  `json:"termVersion" bson:"termVersion"`
	Agreed      bool   `json:"agreed" bson:"agreed"`
}
