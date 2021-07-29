package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type authenticateUserResponseBody struct {
	Data struct {
		AuthenticateUser struct {
			Token string `json:"token"`
		} `json:"authenticateUser"`
	} `json:"data"`
}

func GetNewTokenIfNecessary(token *string, url string, user string, pass string) {
	if *token == "" {
		bodyTpl := `{"query":"{ authenticateUser( input: { email: \"%s\", password: \"%s\" } ) { token } }"}`
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintf(bodyTpl, user, pass))))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonBody, _ := ioutil.ReadAll(resp.Body)
		body := authenticateUserResponseBody{}
		json.Unmarshal(jsonBody, &body)
		*token = body.Data.AuthenticateUser.Token
	}
}
