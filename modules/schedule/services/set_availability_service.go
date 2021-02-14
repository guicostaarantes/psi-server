package services

import (
	"fmt"
	"sort"
	"time"

	"github.com/guicostaarantes/psi-server/modules/schedule/models"
	"github.com/guicostaarantes/psi-server/utils/database"
)

// SetAvailabilityService is a service that sets the availabities for a psychologist
type SetAvailabilityService struct {
	DatabaseUtil               database.IDatabaseUtil
	SecondsLimitAvailability   int64
	SecondsMinimumAvailability int64
}

// Execute is the method that runs the business logic of the service
func (s SetAvailabilityService) Execute(id string, input []*models.SetAvailabilityInput) error {

	newAvailabilities := []interface{}{}

	sort.SliceStable(input, func(i, j int) bool {
		return input[i].Start < input[j].Start
	})

	for key, aval := range input {

		now := time.Now()

		if key == 0 && aval.Start < now.Unix() {
			return fmt.Errorf("availabilities must not start in the past: one availability starting at %d, current time is %d", aval.Start, now.Unix())
		}

		if key > 0 && aval.Start <= input[key-1].End {
			return fmt.Errorf("availabilities must not clash: two or more availabilities include the timestamp %d", aval.Start)
		}

		if aval.End > now.Add(time.Second*time.Duration(s.SecondsLimitAvailability)).Unix() {
			return fmt.Errorf("availabilities must not finish later than %d seconds from now: one availability ending at %d, limit time is %d", s.SecondsLimitAvailability, aval.End, now.Add(time.Second*time.Duration(s.SecondsLimitAvailability)).Unix())
		}

		if aval.Start+s.SecondsMinimumAvailability > aval.End {
			return fmt.Errorf("availabilities must last at least %d seconds: one availability starting at %d, ending at %d", s.SecondsMinimumAvailability, aval.Start, aval.End)
		}

		newAval := models.Availability{
			PsychologistID: id,
			Start:          aval.Start,
			End:            aval.End,
		}

		newAvailabilities = append(newAvailabilities, newAval)

	}

	deleteErr := s.DatabaseUtil.DeleteMany("psi_db", "availabilities", map[string]interface{}{"psychologistId": id})
	if deleteErr != nil {
		return deleteErr
	}

	writeErr := s.DatabaseUtil.InsertMany("psi_db", "availabilities", newAvailabilities)
	if writeErr != nil {
		return writeErr
	}

	return nil

}
