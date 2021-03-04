package services

import (
	models "github.com/guicostaarantes/psi-server/modules/messages/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetMessagesService is a service that sets translated messages
type SetMessagesService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s SetMessagesService) Execute(lang string, input []*models.MessageInput) error {

	messages := []interface{}{}

	for _, msg := range input {
		newMessage := models.Message{
			Lang:  lang,
			Key:   msg.Key,
			Value: msg.Value,
		}
		messages = append(messages, newMessage)
	}

	writeErr := s.DatabaseUtil.InsertMany("psi_db", "messages", messages)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
