package models

// AvailabilityResponse is the schema for information to be sent to the user about a psychologist availability
type AvailabilityResponse struct {
	Start int64 `json:"start" bson:"start"`
	End   int64 `json:"end" bson:"end"`
}
