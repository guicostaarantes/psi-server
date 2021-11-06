package models

// TreatmentPriceRange represents a range of prices that a psychologist will charge per session of a treatment.
type TreatmentPriceRange struct {
	Name         string `json:"name"`
	MinimumPrice int64  `json:"minimumPrice"`
	MaximumPrice int64  `json:"maximumPrice"`
	EligibleFor  string `json:"eligibleFor"`
}

// TreatmentPriceRangeOffering represents the will of a psychologist to allow a patient to initiate a treatment under such price conditions.
type TreatmentPriceRangeOffering struct {
	ID             string `json:"id"`
	PsychologistID string `json:"psychologistId"`
	PriceRangeName string `json:"priceRangeName"`
}
