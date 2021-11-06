package services

import (
	"time"

	"github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetCooldownService is a service that retrieves a cooldown in a database
type GetCooldownService struct {
	OrmUtil orm.IOrmUtil
}

func (s GetCooldownService) Execute(profileID string, profileType models.CooldownProfileType, cooldownType models.CooldownType) (*models.Cooldown, error) {
	cooldown := &models.Cooldown{}

	result := s.OrmUtil.Db().Where(
		"profile_id = ? AND profile_type = ? AND cooldown_type = ? AND valid_until > ?",
		profileID,
		profileType,
		cooldownType,
		time.Now().Unix(),
	).Limit(1).Find(&cooldown)
	if result.Error != nil {
		return nil, result.Error
	}

	if cooldown.ID == "" {
		return nil, nil
	}

	return cooldown, nil
}
