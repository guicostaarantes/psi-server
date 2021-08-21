package models

// CreateTreatmentInput is the schema for information needed to create a new treatment
type CreateTreatmentInput struct {
	Frequency int64 `json:"frequency" bson:"frequency"`
	Phase     int64 `json:"phase" bson:"phase"`
	Duration  int64 `json:"duration" bson:"duration"`
	Price     int64 `json:"price" bson:"price"`
}

// UpdateTreatmentInput is the schema for information needed to update a treatment
type UpdateTreatmentInput struct {
	Frequency int64 `json:"frequency" bson:"frequency"`
	Phase     int64 `json:"phase" bson:"phase"`
	Duration  int64 `json:"duration" bson:"duration"`
	Price     int64 `json:"price" bson:"price"`
}
