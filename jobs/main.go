package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/guicostaarantes/psi-server/jobs/tasks"
)

func print(str string) {
	fmt.Println(str)
}

func main() {
	jobrunnerToken := ""
	url := os.Getenv("PSI_BACKEND_URL")
	jobrunnerUser := os.Getenv("PSI_JOBRUNNER_USERNAME")
	jobrunnerPass := os.Getenv("PSI_JOBRUNNER_PASSWORD")
	getNewTokenFrequency := os.Getenv("PSI_GET_NEW_TOKEN_FREQUENCY")
	processPendingMailFrequency := os.Getenv("PSI_PROCESS_PENDING_MAIL_FREQUENCY")
	createPendingAppointmentsFrequency := os.Getenv("PSI_CREATE_PENDING_APPOINTMENTS_FREQUENCY")

	s := gocron.NewScheduler(time.UTC)
	phase := time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)

	s.Every(getNewTokenFrequency).StartAt(phase).SingletonMode().Do(tasks.GetNewTokenIfNecessary, &jobrunnerToken, url, jobrunnerUser, jobrunnerPass)
	s.Every(processPendingMailFrequency).StartAt(phase).SingletonMode().Do(tasks.ProcessPendingMail, &jobrunnerToken, url)
	s.Every(createPendingAppointmentsFrequency).StartAt(phase).SingletonMode().Do(tasks.CreatePendingAppointments, &jobrunnerToken, url)

	s.StartBlocking()
}
