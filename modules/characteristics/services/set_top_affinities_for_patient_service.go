package services

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	characteristic_models "github.com/guicostaarantes/psi-server/modules/characteristics/models"
	cooldowns_models "github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	cooldowns_services "github.com/guicostaarantes/psi-server/modules/cooldowns/services"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetTopAffinitiesForPatientService is a service that calculates the affinity between a given patient and all psychologists with pending treatments, and saves the most relevant ones to a table
type SetTopAffinitiesForPatientService struct {
	DatabaseUtil        database.IDatabaseUtil
	MaxAffinityNumber   int64
	SaveCooldownService *cooldowns_services.SaveCooldownService
}

// Execute is the method that runs the business logic of the service
func (s SetTopAffinitiesForPatientService) Execute(patientID string) error {

	result := map[string]*models.AffinityScore{}

	// Get possible price ranges
	incomeChar := characteristic_models.CharacteristicChoice{}

	findErr := s.DatabaseUtil.FindOne("characteristic_choices", map[string]interface{}{"profileId": patientID, "characteristicName": "income"}, &incomeChar)
	if findErr != nil {
		return findErr
	}

	if incomeChar.SelectedValue == "" {
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

		for _, v := range strings.Split(priceRange.EligibleFor, ",") {
			if v == incomeChar.SelectedValue {
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
				_, ok := result[priceRangeOffering.PsychologistID]
				if !ok {
					result[priceRangeOffering.PsychologistID] = &models.AffinityScore{}
				}
			}
		}

	}

	// Get patient characteristics choices
	patientCharacteristicChoices := []models.CharacteristicChoice{}

	patientChoicesCursor, findErr := s.DatabaseUtil.FindMany("characteristic_choices", map[string]interface{}{"profileId": patientID})
	if findErr != nil {
		return findErr
	}

	defer patientChoicesCursor.Close(context.Background())

	for patientChoicesCursor.Next(context.Background()) {

		choice := models.CharacteristicChoice{}

		decodeErr := patientChoicesCursor.Decode(&choice)
		if decodeErr != nil {
			return decodeErr
		}

		patientCharacteristicChoices = append(patientCharacteristicChoices, choice)

	}

	// Get all psychologists preferences and calculate score for psychologist
	preferencesCursor, findErr := s.DatabaseUtil.FindMany("preferences", map[string]interface{}{"target": string(models.PsychologistTarget)})
	if findErr != nil {
		return findErr
	}

	defer preferencesCursor.Close(context.Background())

	for preferencesCursor.Next(context.Background()) {

		preference := models.Preference{}

		decodeErr := preferencesCursor.Decode(&preference)
		if decodeErr != nil {
			return decodeErr
		}

		_, ok := result[preference.ProfileID]
		if ok {
			for _, char := range patientCharacteristicChoices {
				if char.CharacteristicName == preference.CharacteristicName && char.SelectedValue == preference.SelectedValue {
					result[preference.ProfileID].ScoreForPsychologist += preference.Weight
				}
			}
		}

	}

	// Get patient preferences
	patientPreferences := []models.Preference{}

	patientPreferencesCursor, findErr := s.DatabaseUtil.FindMany("preferences", map[string]interface{}{"profileId": patientID})
	if findErr != nil {
		return findErr
	}

	defer patientPreferencesCursor.Close(context.Background())

	for patientPreferencesCursor.Next(context.Background()) {

		preference := models.Preference{}

		decodeErr := patientPreferencesCursor.Decode(&preference)
		if decodeErr != nil {
			return decodeErr
		}

		patientPreferences = append(patientPreferences, preference)

	}

	// Get all psychologists characteristic choices and calculate score for patient
	choicesCursor, findErr := s.DatabaseUtil.FindMany("characteristic_choices", map[string]interface{}{"target": string(models.PsychologistTarget)})
	if findErr != nil {
		return findErr
	}

	defer choicesCursor.Close(context.Background())

	for choicesCursor.Next(context.Background()) {

		characteristicChoice := models.CharacteristicChoice{}

		decodeErr := choicesCursor.Decode(&characteristicChoice)
		if decodeErr != nil {
			return decodeErr
		}

		_, ok := result[characteristicChoice.ProfileID]
		if ok {
			for _, pref := range patientPreferences {
				if pref.CharacteristicName == characteristicChoice.CharacteristicName && pref.SelectedValue == characteristicChoice.SelectedValue {
					result[characteristicChoice.ProfileID].ScoreForPatient += pref.Weight
				}
			}
		}

	}

	resultSlice := []models.Affinity{}

	// Transform result map in result slice
	for psychologistID, re := range result {

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
