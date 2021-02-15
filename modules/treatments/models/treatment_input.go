package models

// CreateTreatmentInput is the schema for information needed to create a new treatment
type CreateTreatmentInput struct {
	Duration int64 `json:"duration" bson:"duration"`
	Price    int64 `json:"price" bson:"price"`
	Interval int64 `json:"interval" bson:"interval"`
}

// UpdateTreatmentInput is the schema for information needed to update a treatment
type UpdateTreatmentInput struct {
	Duration int64 `json:"duration" bson:"duration"`
	Price    int64 `json:"price" bson:"price"`
	Interval int64 `json:"interval" bson:"interval"`
}
