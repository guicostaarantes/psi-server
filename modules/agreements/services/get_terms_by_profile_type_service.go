package services

import (
	"github.com/guicostaarantes/psi-server/modules/agreements/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// GetTermsByProfileTypeService is a service that gets the terms for a specific profile type
type GetTermsByProfileTypeService struct {
	OrmUtil orm.IOrmUtil
}

func (s GetTermsByProfileTypeService) Execute(profileType models.TermProfileType) ([]*models.Term, error) {

	terms := []*models.Term{}

	result := s.OrmUtil.Db().Where("profile_type = ?", profileType).Find(&terms)
	if result.Error != nil {
		return nil, result.Error
	}

	return terms, nil

}
