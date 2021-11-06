package services

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	cooldowns_models "github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	cooldowns_services "github.com/guicostaarantes/psi-server/modules/cooldowns/services"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// SetTopAffinitiesForPatientService is a service that calculates the affinity between a given patient and all psychologists with pending treatments, and saves the most relevant ones to a table
type SetTopAffinitiesForPatientService struct {
	DatabaseUtil        database.IDatabaseUtil
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

	possiblePriceRanges := []string{}

	priceRangesCursor, findErr := s.DatabaseUtil.FindMany("treatment_price_ranges", map[string]interface{}{})
	if findErr != nil {
		return findErr
	}

	defer priceRangesCursor.Close(context.Background())

	for priceRangesCursor.Next(context.Background()) {

		priceRange := treatments_models.TreatmentPriceRange{}

		decodeErr := priceRangesCursor.Decode(&priceRange)
		if decodeErr != nil {
			return decodeErr
		}

		for _, pr := range strings.Split(priceRange.EligibleFor, ",") {
			if _, exists := patientChoices["income"][pr]; exists {
				possiblePriceRanges = append(possiblePriceRanges, priceRange.Name)
			}
		}

	}

	// Check if psychologist has at least one treatment price range offering with a possible price range
	priceRangeOfferingsCursor, findErr := s.DatabaseUtil.FindMany("treatment_price_range_offerings", map[string]interface{}{})
	if findErr != nil {
		return findErr
	}

	defer priceRangeOfferingsCursor.Close(context.Background())

	for priceRangeOfferingsCursor.Next(context.Background()) {

		priceRangeOffering := treatments_models.TreatmentPriceRangeOffering{}

		decodeErr := priceRangeOfferingsCursor.Decode(&priceRangeOffering)
		if decodeErr != nil {
			return decodeErr
		}

		for _, v := range possiblePriceRanges {
			if v == priceRangeOffering.PriceRangeName {
				// If there is, add psychologist to list
				_, ok := affinityResult[priceRangeOffering.PsychologistID]
				if !ok {
					affinityResult[priceRangeOffering.PsychologistID] = &models.AffinityScore{}
				}
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

	resultSlice := []models.Affinity{}

	// Transform result map in result slice
	for psychologistID, re := range affinityResult {

		if re.ScoreForPatient >= 0 && re.ScoreForPsychologist >= 0 {
			resultSlice = append(resultSlice, models.Affinity{
				PatientID:            patientID,
				PsychologistID:       psychologistID,
				CreatedAt:            time.Now().Unix(),
				ScoreForPatient:      re.ScoreForPatient,
				ScoreForPsychologist: re.ScoreForPsychologist,
			})
		}

	}

	// Sort based on sum of points and cut only the most relevant limited to s.MaxAffinityNumber
	sort.SliceStable(resultSlice, func(i int, j int) bool {
		return resultSlice[i].ScoreForPatient+resultSlice[i].ScoreForPsychologist > resultSlice[j].ScoreForPatient+resultSlice[j].ScoreForPsychologist
	})
	if len(resultSlice) > int(s.MaxAffinityNumber) {
		resultSlice = resultSlice[:s.MaxAffinityNumber]
	}

	// Transferring to []interface{} in order to add to database
	topAffinities := []interface{}{}
	for _, slice := range resultSlice {
		topAffinities = append(topAffinities, slice)
	}

	deleteErr := s.DatabaseUtil.DeleteMany("top_affinities", map[string]interface{}{"patientId": patientID})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("top_affinities", topAffinities)
	if writeErr != nil {
		return writeErr
	}

	saveErr := s.SaveCooldownService.Execute(patientID, cooldowns_models.Patient, cooldowns_models.TopAffinitiesSet)
	if saveErr != nil {
		return saveErr
	}

	return nil

}
