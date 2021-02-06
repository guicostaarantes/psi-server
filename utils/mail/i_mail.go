package mail

// IMailUtil is an abstraction for a email sender
type IMailUtil interface {
	GetMockedMessages() *[]map[string]interface{}
	Send(msg Message) error
}

// Message holds the payload of a mail message
type Message struct {
	FromAddress string   `json:"fromAddress" bson:"fromAddress"`
	FromName    string   `json:"fromName" bson:"fromName"`
	To          []string `json:"to" bson:"to"`
	Cc          []string `json:"cc" bson:"cc"`
	Cco         []string `json:"cco" bson:"cco"`
	Subject     string   `json:"subject" bson:"subject"`
	HTML        string   `json:"html" bson:"html"`
}
