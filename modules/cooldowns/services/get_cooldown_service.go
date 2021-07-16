package services

import (
	"context"
	"time"

	"github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// GetCooldownService is a service that retrieves a cooldown in a database
type GetCooldownService struct {
	DatabaseUtil database.IDatabaseUtil
}

func (s GetCooldownService) Execute(profileID string, profileType models.CooldownProfileType, cooldownType models.CooldownType) (*models.Cooldown, error) {
	cursor, findErr := s.DatabaseUtil.FindMany("psi_db", "cooldowns", map[string]interface{}{"profileId": profileID, "profileType": profileType, "cooldownType": cooldownType})
	if findErr != nil {
		return nil, findErr
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		cooldown := models.Cooldown{}

		decodeErr := cursor.Decode(&cooldown)
		if decodeErr != nil {
			return nil, decodeErr
		}

		if cooldown.ValidUntil > time.Now().Unix() {
			return &cooldown, nil
		}
	}

	return nil, nil
}
