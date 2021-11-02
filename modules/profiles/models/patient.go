package models

// Patient is the schema for the profile of a patient
type Patient struct {
	ID        string `json:"id"`
	UserID    string `json:"userId" gorm:"index"`
	FullName  string `json:"fullName"`
	LikeName  string `json:"likeName"`
	BirthDate int64  `json:"birthDate"`
	City      string `json:"city"`
	Avatar    string `json:"avatar"`
}
