package mail

// IMailUtil is an abstraction for a email sender
type IMailUtil interface {
	GetMockedMessages() (*[]map[string]interface{}, error)
	Send(msg Message) error
}

// Message holds the payload of a mail message
type Message struct {
	FromAddress string   `json:"fromAddress"`
	FromName    string   `json:"fromName"`
	To          []string `json:"to"`
	Cc          []string `json:"cc"`
	Cco         []string `json:"cco"`
	Subject     string   `json:"subject"`
	HTML        string   `json:"html"`
}
