package services

import (
	"context"
	"sort"
	"time"

	"github.com/guicostaarantes/psi-server/modules/characteristics/models"
	cooldowns_models "github.com/guicostaarantes/psi-server/modules/cooldowns/models"
	cooldowns_services "github.com/guicostaarantes/psi-server/modules/cooldowns/services"
	treatments_models "github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetTopAffinitiesForPatientService is a service that calculates the affinity between a given patient and all psychologists with pending treatments, and saves the most relevant ones to a table
type SetTopAffinitiesForPatientService struct {
	DatabaseUtil        database.IDatabaseUtil
	MaxAffinityNumber   int64
	SaveCooldownService cooldowns_services.SaveCooldownService
}

// Execute is the method that runs the business logic of the service
func (s SetTopAffinitiesForPatientService) Execute(patientID string) error {

	result := map[string]*models.AffinityScore{}

	// Get available psychologist IDs from PENDING treatments
	treatmentCursor, findErr := s.DatabaseUtil.FindMany("treatments", map[string]interface{}{"status": string(treatments_models.Pending)})
	if findErr != nil {
		return findErr
	}

	defer treatmentCursor.Close(context.Background())

	for treatmentCursor.Next(context.Background()) {

		treatment := treatments_models.Treatment{}

		decodeErr := treatmentCursor.Decode(&treatment)
		if decodeErr != nil {
			return decodeErr
		}

		_, ok := result[treatment.PsychologistID]
		if !ok {
			result[treatment.PsychologistID] = &models.AffinityScore{}
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
