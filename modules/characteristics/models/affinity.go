package models

import (
	"time"

	"gorm.io/gorm"
)

// AffinityScore represents in the result field how much likely it is for a treatment to be succesful between psychologist and patient, based on their characteristics and preferences
type AffinityScore struct {
	ScoreForPatient      int64 `json:"scoreForPatient"`
	ScoreForPsychologist int64 `json:"scoreForPsychologist"`
}

// Affinity is the representation in the database of a calculation of affinity between psychologist and patient
type Affinity struct {
	ID                   string         `json:"id" gorm:"primaryKey"`
	CreatedAt            time.Time      `json:"createdAt`
	UpdatedAt            time.Time      `json:"updatedAt`
	DeletedAt            gorm.DeletedAt `gorm:"index"`
	PatientID            string         `json:"patientId" gorm:"index"`
	PsychologistID       string         `json:"psychologistId"`
	ScoreForPatient      int64          `json:"scoreForPatient"`
	ScoreForPsychologist int64          `json:"scoreForPsychologist"`
}
