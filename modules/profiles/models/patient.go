package profiles_models

import (
	"time"

	"gorm.io/gorm"
)

// Patient is the schema for the profile of a patient
type Patient struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt`
	UpdatedAt time.Time      `json:"updatedAt`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    string         `json:"userId" gorm:"index"`
	FullName  string         `json:"fullName"`
	LikeName  string         `json:"likeName"`
	BirthDate time.Time      `json:"birthDate"`
	City      string         `json:"city"`
	Avatar    string         `json:"avatar"`
}
