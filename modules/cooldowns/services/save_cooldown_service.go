package services

import (
	"errors"
	"time"

	"github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/identifier"
)

// SaveCooldownService is a service that stores a cooldown in a database
type SaveCooldownService struct {
	DatabaseUtil                 database.IDatabaseUtil
	IdentifierUtil               identifier.IIdentifierUtil
	TopAffinitiesCooldownSeconds int64
}

func (s SaveCooldownService) Execute(profileID string, profileType models.CooldownProfileType, cooldownType models.CooldownType) error {
	_, cooldownID, cooldownIDErr := s.IdentifierUtil.GenerateIdentifier()
	if cooldownIDErr != nil {
		return cooldownIDErr
	}

	var duration int64
	switch cooldownType {
	case models.TopAffinitiesSet:
		duration = s.TopAffinitiesCooldownSeconds
	default:
		return errors.New("cooldownType does not have a duration")
	}

	createdAt := time.Now().Unix()
	validUntil := createdAt + duration

	cooldown := models.Cooldown{
		ID:           cooldownID,
		ProfileID:    profileID,
		ProfileType:  profileType,
		CooldownType: cooldownType,
		CreatedAt:    createdAt,
		ValidUntil:   validUntil,
	}

	insertErr := s.DatabaseUtil.InsertOne("cooldowns", cooldown)
	if insertErr != nil {
		return insertErr
	}

	return nil
}
