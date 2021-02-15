package models

// SetAvailabilityInput is the schema for information needed to create an availability
type SetAvailabilityInput struct {
	Start int64 `json:"start" bson:"start"`
	End   int64 `json:"end" bson:"end"`
}
