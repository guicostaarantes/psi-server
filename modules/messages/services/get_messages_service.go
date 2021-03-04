package services

import (
	"context"

	models "github.com/guicostaarantes/psi-server/modules/messages/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetMessagesService is a service that gets translated messages
type GetMessagesService struct {
	DatabaseUtil database.IDatabaseUtil
}

// Execute is the method that runs the business logic of the service
func (s GetMessagesService) Execute(lang string, keys []string) ([]*models.Message, error) {

	messages := []*models.Message{}

	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "messages", map[string]interface{}{"lang": lang})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		msg := models.Message{}

		decodeErr := cursor.Decode(&msg)
		if decodeErr != nil {
			return nil, decodeErr
		}

		for _, key := range keys {
			if key == msg.Key {
				messages = append(messages, &msg)
			}
		}
	}

	return messages, nil

}
