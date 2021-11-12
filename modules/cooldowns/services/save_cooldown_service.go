package services

import (
	"errors"
	"time"

	"github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SaveCooldownService is a service that stores a cooldown in a database
type SaveCooldownService struct {
	IdentifierUtil                    identifier.IIdentifierUtil
	OrmUtil                           orm.IOrmUtil
	InterruptTreatmentCooldownSeconds int64
	TopAffinitiesCooldownSeconds      int64
}

func (s SaveCooldownService) Execute(profileID string, profileType models.CooldownProfileType, cooldownType models.CooldownType) error {
	_, cooldownID, cooldownIDErr := s.IdentifierUtil.GenerateIdentifier()
	if cooldownIDErr != nil {
		return cooldownIDErr
	}

	var duration int64
	switch cooldownType {
	case models.TreatmentInterrupted:
		duration = s.InterruptTreatmentCooldownSeconds
	case models.TopAffinitiesSet:
		duration = s.TopAffinitiesCooldownSeconds
	default:
		return errors.New("cooldownType does not have a duration")
	}

	validUntil := time.Now().Unix() + duration

	cooldown := models.Cooldown{
		ID:           cooldownID,
		ProfileID:    profileID,
		ProfileType:  profileType,
		CooldownType: cooldownType,
		ValidUntil:   validUntil,
	}

	result := s.OrmUtil.Db().Create(&cooldown)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
