package models

// Availability represents a time period in which a psychologist is available to receive appointments
type Availability struct {
	PsychologistID string `json:"psychologistId" bson:"psychologistId"`
	Start          int64  `json:"start" bson:"start"`
	End            int64  `json:"end" bson:"end"`
}
