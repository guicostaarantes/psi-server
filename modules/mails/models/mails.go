package mails_models

import "time"

// TransientMailMessage holds the MailMessage and a flag to know if it was handled
type TransientMailMessage struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"createdAt`
	UpdatedAt   time.Time `json:"updatedAt`
	FromAddress string    `json:"fromAddress"`
	FromName    string    `json:"fromName"`
	To          string    `json:"to"`
	Cc          string    `json:"cc"`
	Cco         string    `json:"cco"`
	Subject     string    `json:"subject"`
	Html        string    `json:"html"`
	Processed   bool      `json:"processed" gorm:"index"`
}
