package agreements_models

import (
	"time"

	"gorm.io/gorm"
)

// Agreement is an evidence that a term was agreed or disagreed by a profile
type Agreement struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"createdAt`
	UpdatedAt   time.Time      `json:"updatedAt`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	TermName    string         `json:"termName"`
	TermVersion int64          `json:"termVersion"`
	ProfileID   string         `json:"profileId"`
}
