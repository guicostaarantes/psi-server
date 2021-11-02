package models

// Psychologist is the schema for the profile of a psychologist
type Psychologist struct {
	ID        string `json:"id"`
	UserID    string `json:"userId" gorm:"index"`
	FullName  string `json:"fullName"`
	LikeName  string `json:"likeName"`
	BirthDate int64  `json:"birthDate"`
	City      string `json:"city"`
	Crp       string `json:"crp"`
	Whatsapp  string `json:"whatsapp"`
	Instagram string `json:"instagram"`
	Bio       string `json:"bio"`
	Avatar    string `json:"avatar"`
}
