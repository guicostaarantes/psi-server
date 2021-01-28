package mails_models

// TransientMailMessage holds the MailMessage and a flag to know if it was handled
type TransientMailMessage struct {
	ID          string   `json:"id" bson:"id"`
	FromAddress string   `json:"fromAddress" bson:"fromAddress"`
	FromName    string   `json:"fromName" bson:"fromName"`
	To          []string `json:"to" bson:"to"`
	Cc          []string `json:"cc" bson:"cc"`
	Cco         []string `json:"cco" bson:"cco"`
	Subject     string   `json:"subject" bson:"subject"`
	Html        string   `json:"html" bson:"html"`
	Processed   bool     `json:"processed" bson:"processed"`
}
