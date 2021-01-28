package mail

// IMailUtil is an abstraction for a email sender
type IMailUtil interface {
	Send(msg MailMessage) error
}

// MailMessage holds the payload of a mail message
type MailMessage struct {
	FromAddress string   `json:"fromAddress" bson:"fromAddress"`
	FromName    string   `json:"fromName" bson:"fromName"`
	To          []string `json:"to" bson:"to"`
	Cc          []string `json:"cc" bson:"cc"`
	Cco         []string `json:"cco" bson:"cco"`
	Subject     string   `json:"subject" bson:"subject"`
	Html        string   `json:"html" bson:"html"`
}
