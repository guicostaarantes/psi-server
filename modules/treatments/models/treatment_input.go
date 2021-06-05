package models

// CreateTreatmentInput is the schema for information needed to create a new treatment
type CreateTreatmentInput struct {
	WeeklyStart int64 `json:"weeklyStart" bson:"weeklyStart"`
	Duration    int64 `json:"duration" bson:"duration"`
	Price       int64 `json:"price" bson:"price"`
}

// UpdateTreatmentInput is the schema for information needed to update a treatment
type UpdateTreatmentInput struct {
	WeeklyStart int64 `json:"weeklyStart" bson:"weeklyStart"`
	Duration    int64 `json:"duration" bson:"duration"`
	Price       int64 `json:"price" bson:"price"`
}
