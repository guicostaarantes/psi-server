package treatments_services

import (
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SetTreatmentPriceRangesService is a service that sets all possible patient characteristics
type SetTreatmentPriceRangesService struct {
	IdentifierUtil identifier.IIdentifierUtil
	OrmUtil        orm.IOrmUtil
}

// Execute is the method that runs the business logic of the service
func (s SetTreatmentPriceRangesService) Execute(input []*treatments_models.TreatmentPriceRange) error {

	currentTreatmentPriceRanges := []*treatments_models.TreatmentPriceRange{}

	result := s.OrmUtil.Db().Find(&currentTreatmentPriceRanges)
	if result.Error != nil {
		return result.Error
	}

	currentPriceRanges := map[string]*treatments_models.TreatmentPriceRange{}

	for _, pr := range currentTreatmentPriceRanges {
		currentPriceRanges[pr.Name] = pr
	}

	for _, pr := range input {
		if _, exists := currentPriceRanges[pr.Name]; exists {

			currentPriceRanges[pr.Name].MinimumPrice = pr.MinimumPrice
			currentPriceRanges[pr.Name].MaximumPrice = pr.MaximumPrice
			currentPriceRanges[pr.Name].EligibleFor = pr.EligibleFor

			result := s.OrmUtil.Db().Save(currentPriceRanges[pr.Name])
			if result.Error != nil {
				return result.Error
			}

			delete(currentPriceRanges, pr.Name)

		} else {

			_, prID, prIDErr := s.IdentifierUtil.GenerateIdentifier()
			if prIDErr != nil {
				return prIDErr
			}

			result := s.OrmUtil.Db().Create(&treatments_models.TreatmentPriceRange{
				ID:           prID,
				Name:         pr.Name,
				MinimumPrice: pr.MinimumPrice,
				MaximumPrice: pr.MaximumPrice,
				EligibleFor:  pr.EligibleFor,
			})
			if result.Error != nil {
				return result.Error
			}

		}

		// Deleting remaining prs
		for _, pr := range currentPriceRanges {
			result := s.OrmUtil.Db().Delete(&pr)
			if result.Error != nil {
				return result.Error
			}
		}
	}

	return nil

}
