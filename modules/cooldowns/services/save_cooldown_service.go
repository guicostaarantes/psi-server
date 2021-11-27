package cooldowns_services

import (
	"errors"
	"time"

	cooldowns_models "github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SaveCooldownService is a service that stores a cooldown in a database
type SaveCooldownService struct {
	IdentifierUtil                     identifier.IIdentifierUtil
	OrmUtil                            orm.IOrmUtil
	InterruptTreatmentCooldownDuration time.Duration
	TopAffinitiesCooldownDuration      time.Duration
}

func (s SaveCooldownService) Execute(profileID string, profileType cooldowns_models.CooldownProfileType, cooldownType cooldowns_models.CooldownType) error {
	_, cooldownID, cooldownIDErr := s.IdentifierUtil.GenerateIdentifier()
	if cooldownIDErr != nil {
		return cooldownIDErr
	}

	var duration time.Duration
	switch cooldownType {
	case cooldowns_models.TreatmentInterrupted:
		duration = s.InterruptTreatmentCooldownDuration
	case cooldowns_models.TopAffinitiesSet:
		duration = s.TopAffinitiesCooldownDuration
	default:
		return errors.New("cooldownType does not have a duration")
	}

	validUntil := time.Now().Add(duration)

	cooldown := cooldowns_models.Cooldown{
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
