package mail

type FakeMailUtil struct {
	MockedMessages *[]map[string]interface{}
}

func (s FakeMailUtil) GetMockedMessages() (*[]map[string]interface{}, error) {
	return s.MockedMessages, nil
}

func (s FakeMailUtil) Send(msg Message) error {
	mail := map[string]interface{}{
		"to":      msg.To,
		"cc":      msg.Cc,
		"cco":     msg.Cco,
		"subject": msg.Subject,
		"body":    msg.HTML,
	}

	*s.MockedMessages = append(*s.MockedMessages, mail)

	return nil
}
