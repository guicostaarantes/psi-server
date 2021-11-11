package models

// CreateTreatmentInput is the schema for information needed to create a new treatment
type CreateTreatmentInput struct {
	Frequency      int64  `json:"frequency"`
	Phase          int64  `json:"phase"`
	Duration       int64  `json:"duration"`
	PriceRangeName string `json:"priceRangeName"`
}

// UpdateTreatmentInput is the schema for information needed to update a treatment
type UpdateTreatmentInput struct {
	Frequency      int64  `json:"frequency"`
	Phase          int64  `json:"phase"`
	Duration       int64  `json:"duration"`
	PriceRangeName string `json:"priceRangeName"`
}
