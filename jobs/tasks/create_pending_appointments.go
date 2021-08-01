package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type createPendingAppointmentsResponseBody struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func CreatePendingAppointments(token *string, url string) {
	if *token != "" {
		bodyTpl := `{"query":"mutation { createPendingAppointments }"}`
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyTpl)))
		req.Header.Set("Authorization", *token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonBody, _ := ioutil.ReadAll(resp.Body)
		body := createPendingAppointmentsResponseBody{}
		json.Unmarshal(jsonBody, &body)
		if len(body.Errors) > 0 {
			if body.Errors[0].Message == "forbidden" {
				*token = ""
			} else {
				log.Fatalf(`CreatePendingAppointments returned error %s`, body.Errors[0].Message)
			}
		}
	}
}
