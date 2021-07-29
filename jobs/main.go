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
	url := os.Getenv("PSI_BACKEND_URL")
	jobrunnerUser := os.Getenv("PSI_JOBRUNNER_USERNAME")
	jobrunnerPass := os.Getenv("PSI_JOBRUNNER_PASSWORD")
	jobrunnerToken := ""

	s := gocron.NewScheduler(time.UTC)

	s.Every(10).Seconds().SingletonMode().Do(tasks.GetNewTokenIfNecessary, &jobrunnerToken, url, jobrunnerUser, jobrunnerPass)
	s.Every(10).Seconds().SingletonMode().Do(tasks.ProcessPendingMail, &jobrunnerToken, url)
	s.Every(1).Day().At("12:00").SingletonMode().Do(tasks.CreatePendingAppointments, &jobrunnerToken, url)

	s.StartBlocking()
}
