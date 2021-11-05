package services

import (
	"errors"

	"github.com/guicostaarantes/psi-server/modules/treatments/models"
	"github.com/guicostaarantes/psi-server/utils/orm"
)

// CheckTreatmentCollisionService is a service that checks if a treatment period collides with others from the same psychologist
type CheckTreatmentCollisionService struct {
	OrmUtil                 orm.IOrmUtil
	ScheduleIntervalSeconds int64
}

func LCM(a, b int64) int64 {
	c := a
	d := b
	for d != 0 {
		t := d
		d = c % d
		c = t
	}
	return a * b / c
}

// Execute is the method that runs the business logic of the service
func (s CheckTreatmentCollisionService) Execute(psychologistID string, frequency int64, phase int64, duration int64, updatingID string) error {

	if phase >= s.ScheduleIntervalSeconds*frequency {
		return errors.New("phase cannot be bigger than the schedule interval")
	}

	psychologistTreatments := []*models.Treatment{}

	result := s.OrmUtil.Db().Where("psychologist_id = ?", psychologistID).Find(&psychologistTreatments)
	if result.Error != nil {
		return result.Error
	}

	for _, treatment := range psychologistTreatments {
		if treatment.Status != models.Pending && treatment.Status != models.Active {
			continue
		}

		lcm := LCM(frequency, treatment.Frequency)

		candidateDates := [][]int64{}
		treatmentDates := [][]int64{}

		for counter := int64(0); counter < lcm; counter++ {
			if counter%frequency == 0 {
				candidateStart := counter*s.ScheduleIntervalSeconds + phase
				candidateEnd := (candidateStart + duration) % (lcm * s.ScheduleIntervalSeconds)
				candidateDates = append(candidateDates, []int64{candidateStart, candidateEnd})
			}
			if counter%treatment.Frequency == 0 {
				treatmentStart := counter*s.ScheduleIntervalSeconds + treatment.Phase
				treatmentEnd := (treatmentStart + treatment.Duration) % (lcm * s.ScheduleIntervalSeconds)
				treatmentDates = append(treatmentDates, []int64{treatmentStart, treatmentEnd})
			}
		}

		for _, candidateValues := range candidateDates {
			for _, treatmentValues := range treatmentDates {
				// If 3 of the 4 conditions below are true, it means there is no clash between treatments
				noClash := 0
				if candidateValues[0] < candidateValues[1] {
					noClash++
				}
				if candidateValues[1] <= treatmentValues[0] {
					noClash++
				}
				if treatmentValues[0] < treatmentValues[1] {
					noClash++
				}
				if treatmentValues[1] <= candidateValues[0] {
					noClash++
				}

				if noClash < 3 && treatment.ID != updatingID {
					return errors.New("there is another treatment in the same period")
				}
			}
		}
	}

	return nil

}
