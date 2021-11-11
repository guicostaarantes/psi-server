package services

import (
	"github.com/guicostaarantes/psi-server/modules/agreements/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetAgreementsByProfileIdService is a service that gets the agreements for a specific profileId
type GetAgreementsByProfileIdService struct {
	OrmUtil orm.IOrmUtil
}

func (s GetAgreementsByProfileIdService) Execute(profileID string, profileType models.TermProfileType) ([]*models.Agreement, error) {

	agreements := []*models.Agreement{}

	result := s.OrmUtil.Db().Where("profile_id = ?", profileID).Find(&agreements)
	if result.Error != nil {
		return nil, result.Error
	}

	return agreements, nil

}
