package mail

type mockMailer struct {
	mockedMessages *[]map[string]interface{}
}

func (s mockMailer) GetMockedMessages() *[]map[string]interface{} {
	return s.mockedMessages
}

func (s mockMailer) Send(msg Message) error {
	mail := map[string]interface{}{
		"to":      msg.To,
		"cc":      msg.Cc,
		"cco":     msg.Cco,
		"subject": msg.Subject,
		"body":    msg.HTML,
	}

	*s.mockedMessages = append(*s.mockedMessages, mail)

	return nil
}

// MockMailUtil is an implementation of IMailUtil that writes to disk
var MockMailUtil = mockMailer{
	mockedMessages: &[]map[string]interface{}{},
}
