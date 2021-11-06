package services

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	cooldowns_models "github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	cooldowns_services "github.com/guicostaarantes/psi-server/modules/cooldowns/services"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SetTopAffinitiesForPatientService is a service that calculates the affinity between a given patient and all psychologists with pending treatments, and saves the most relevant ones to a table
type SetTopAffinitiesForPatientService struct {
	OrmUtil             orm.IOrmUtil
	MaxAffinityNumber   int64
	SaveCooldownService *cooldowns_services.SaveCooldownService
}

// Execute is the method that runs the business logic of the service
func (s SetTopAffinitiesForPatientService) Execute(patientID string) error {

	affinityResult := map[string]*models.AffinityScore{}

	// Get patient characteristics choices
	patientCharacteristicChoices := []*models.CharacteristicChoice{}

	result := s.OrmUtil.Db().Where("profile_id = ?", patientID).Find(&patientCharacteristicChoices)
	if result.Error != nil {
		return result.Error
	}

	// patientChoices[characteristicName][selectedValue] = true if exists, undefined otherwise
	patientChoices := map[string]map[string]bool{}

	for _, choice := range patientCharacteristicChoices {
		if _, exists := patientChoices[choice.CharacteristicName]; !exists {
			patientChoices[choice.CharacteristicName] = map[string]bool{}
		}
		patientChoices[choice.CharacteristicName][choice.SelectedValue] = true
	}

	// Get possible price ranges
	if len(patientChoices["income"]) == 0 {
		return errors.New("missing income for patient")
	}

	possiblePriceRanges := map[string]bool{}
	priceRanges := []*treatments_models.TreatmentPriceRange{}

	result = s.OrmUtil.Db().Find(&priceRanges)
	if result.Error != nil {
		return result.Error
	}

	for _, priceRange := range priceRanges {
		for _, pr := range strings.Split(priceRange.EligibleFor, ",") {
			if _, exists := patientChoices["income"][pr]; exists {
				possiblePriceRanges[priceRange.Name] = true
			}
		}
	}

	// Check if psychologist has at least one treatment price range offering with a possible price range
	priceRangesOfferings := []*treatments_models.TreatmentPriceRangeOffering{}

	result = s.OrmUtil.Db().Find(&priceRangesOfferings)
	if result.Error != nil {
		return result.Error
	}

	for _, priceRangeOffering := range priceRangesOfferings {
		if _, exists := possiblePriceRanges[priceRangeOffering.PriceRangeName]; exists {
			if _, exists := affinityResult[priceRangeOffering.PsychologistID]; !exists {
				affinityResult[priceRangeOffering.PsychologistID] = &models.AffinityScore{}
			}
		}
	}

	// Get all psychologists preferences and calculate score for psychologist
	psychologistPreferences := []*models.Preference{}

	result = s.OrmUtil.Db().Where("target = ?", models.PsychologistTarget).Find(&psychologistPreferences)
	if result.Error != nil {
		return result.Error
	}

	for _, preference := range psychologistPreferences {
		if _, exists := affinityResult[preference.ProfileID]; exists {
			if _, exists := patientChoices[preference.CharacteristicName][preference.SelectedValue]; exists {
				affinityResult[preference.ProfileID].ScoreForPsychologist += preference.Weight
			}
		}
	}

	// Get patient preferences
	patientPreferences := []*models.Preference{}

	result = s.OrmUtil.Db().Where("profile_id = ?", patientID).Find(&patientPreferences)
	if result.Error != nil {
		return result.Error
	}

	// patientPrefs[characteristicName][selectedValue] = weight
	patientPrefs := map[string]map[string]int64{}

	for _, pref := range patientPreferences {
		if _, exists := patientPrefs[pref.CharacteristicName]; !exists {
			patientPrefs[pref.CharacteristicName] = map[string]int64{}
		}
		patientPrefs[pref.CharacteristicName][pref.SelectedValue] = pref.Weight
	}

	// Get all psychologists characteristic choices and calculate score for patient
	psychologistChoices := []*models.CharacteristicChoice{}

	result = s.OrmUtil.Db().Where("target = ?", models.PsychologistTarget).Find(&psychologistChoices)
	if result.Error != nil {
		return result.Error
	}

	for _, choice := range psychologistChoices {
		if _, exists := affinityResult[choice.ProfileID]; exists {
			if weight, exists := patientPrefs[choice.CharacteristicName][choice.SelectedValue]; exists {
				affinityResult[choice.ProfileID].ScoreForPatient += weight
			}
		}
	}

	topAffinities := []*models.Affinity{}

	// Transform result map in result slice
	for psychologistID, re := range affinityResult {

		if re.ScoreForPatient >= 0 && re.ScoreForPsychologist >= 0 {
			topAffinities = append(topAffinities, &models.Affinity{
				PatientID:            patientID,
				PsychologistID:       psychologistID,
				CreatedAt:            time.Now().Unix(),
				ScoreForPatient:      re.ScoreForPatient,
				ScoreForPsychologist: re.ScoreForPsychologist,
			})
		}

	}

	// Sort based on sum of points and cut only the most relevant limited to s.MaxAffinityNumber
	sort.SliceStable(topAffinities, func(i int, j int) bool {
		return topAffinities[i].ScoreForPatient+topAffinities[i].ScoreForPsychologist > topAffinities[j].ScoreForPatient+topAffinities[j].ScoreForPsychologist
	})
	if len(topAffinities) > int(s.MaxAffinityNumber) {
		topAffinities = topAffinities[:s.MaxAffinityNumber]
	}

	result = s.OrmUtil.Db().Delete(&models.Affinity{}, "patient_id = ?", patientID)
	if result.Error != nil {
		return result.Error
	}

	result = s.OrmUtil.Db().Create(&topAffinities)
	if result.Error != nil {
		return result.Error
	}

	saveErr := s.SaveCooldownService.Execute(patientID, cooldowns_models.Patient, cooldowns_models.TopAffinitiesSet)
	if saveErr != nil {
		return saveErr
	}

	return nil

}
