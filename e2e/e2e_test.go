package e2e

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/guicostaarantes/psi-server/graph"
	"github.com/guicostaarantes/psi-server/graph/resolvers"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/mail"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/serializing"
	"github.com/guicostaarantes/psi-server/utils/token"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
)

func gql(router *chi.Mux, query string, headers map[string][]string) *httptest.ResponseRecorder {

	body := fmt.Sprintf(`{"query": %q}`, query)

	request := httptest.NewRequest(http.MethodPost, "/gql", strings.NewReader(body))

	request.Header = headers
	request.Header["Content-Type"] = []string{"application/json"}

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	return response

}

func TestEnd2End(t *testing.T) {

	storedVariables := map[string]string{}

	res := &resolvers.Resolver{
		DatabaseUtil:           database.MockDatabaseUtil,
		HashUtil:               hash.BcryptHashUtil,
		IdentifierUtil:         identifier.UUIDIdentifierUtil,
		MailUtil:               mail.MockMailUtil,
		MatchUtil:              match.RegexpMatchUtil,
		SerializingUtil:        serializing.JSONSerializingUtil,
		TokenUtil:              token.RngTokenUtil,
		SecondsToCooldownReset: int64(86400),
		SecondsToExpire:        int64(1800),
		SecondsToExpireReset:   int64(86400),
	}

	os.Setenv("PSI_BOOTSTRAP_USER", "coordinator@psi.com.br|Abc123!@#")

	router := graph.CreateServer(res)

	t.Run("should not log in with incorrect email", func(t *testing.T) {

		query := `{
			authenticateUser(input: { 
				email: "coordinator@psi.com", 
				password: "Abc123!@#", 
				ipAddress: "100.100.100.100" 
			}) { 
				token 
				expiresAt 
			} 
		}`

		response := gql(router, query, map[string][]string{})

		assert.Equal(t, "{\"errors\":[{\"message\":\"incorrect credentials\",\"path\":[\"authenticateUser\"]}],\"data\":null}", response.Body.String())

	})

	t.Run("should not log in with incorrect password", func(t *testing.T) {

		query := `{
			authenticateUser(input: { 
				email: "coordinator@psi.com.br", 
				password: "123!@#Abc", 
				ipAddress: "100.100.100.100" 
			}) { 
				token 
				expiresAt 
			} 
		}`

		response := gql(router, query, map[string][]string{})

		assert.Equal(t, "{\"errors\":[{\"message\":\"incorrect credentials\",\"path\":[\"authenticateUser\"]}],\"data\":null}", response.Body.String())

	})

	t.Run("should log in as bootstrap coordinator", func(t *testing.T) {

		query := `{
			authenticateUser(input: { 
				email: "coordinator@psi.com.br", 
				password: "Abc123!@#", 
				ipAddress: "100.100.100.100" 
			}) { 
				token 
				expiresAt 
			} 
		}`

		response := gql(router, query, map[string][]string{})

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		expiresAt := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "expiresAt")
		assert.NotEqual(t, time.Now().Unix()+res.SecondsToExpire, expiresAt)

		storedVariables["coordinator_token"] = token

	})

	t.Run("should get own user data", func(t *testing.T) {

		query := `{
			getOwnUser {
				email
				firstName
				lastName
			}
		}`

		response := gql(router, query, map[string][]string{
			"Authorization": {storedVariables["coordinator_token"]},
		})

		assert.Equal(t, "{\"data\":{\"getOwnUser\":{\"email\":\"coordinator@psi.com.br\",\"firstName\":\"Bootstrap\",\"lastName\":\"User\"}}}", response.Body.String())
	})

	t.Run("should create psychologist user", func(t *testing.T) {

		query := `mutation {
			createPsychologistUser(input: {
				email: "alice.johnson@psi.com.br",
				firstName: "Alice",
				lastName: "Johnson"
			})
		}`

		response := gql(router, query, map[string][]string{
			"Authorization": {storedVariables["coordinator_token"]},
		})

		assert.Equal(t, "{\"data\":{\"createPsychologistUser\":null}}", response.Body.String())
	})

	t.Run("should not create psychologist user with same email", func(t *testing.T) {

		query := `mutation {
			createPsychologistUser(input: {
				email: "alice.johnson@psi.com.br",
				firstName: "Alice",
				lastName: "Johnson"
			})
		}`

		response := gql(router, query, map[string][]string{
			"Authorization": {storedVariables["coordinator_token"]},
		})

		assert.Equal(t, "{\"errors\":[{\"message\":\"user with same email already exists\",\"path\":[\"createPsychologistUser\"]}],\"data\":{\"createPsychologistUser\":null}}", response.Body.String())
	})

	t.Run("should reset password with token sent via email", func(t *testing.T) {

		query := `mutation {
			processPendingMail
		}`

		response := gql(router, query, map[string][]string{
			"Authorization": {storedVariables["coordinator_token"]},
		})

		assert.Equal(t, "{\"data\":{\"processPendingMail\":null}}", response.Body.String())

		mailbox := res.MailUtil.GetMockedMessages()
		var mailBody string

		for _, mail := range *mailbox {
			fmt.Printf("%v\n", mail)
			if reflect.DeepEqual(mail["to"], []string{"alice.johnson@psi.com.br"}) && mail["subject"] == "Bem-vindo ao PSI" {
				mailBody = mail["body"].(string)
				break
			}
		}

		regex := regexp.MustCompile("token=(?P<token>[A-Za-z0-9]{64})")
		match := regex.FindStringSubmatch(mailBody)

		query = fmt.Sprintf(`mutation {
			resetPassword(input: {
				token: %q,
				password: "Def456$%%^"
			})
		}`, match[1])

		fmt.Println(query)

		response = gql(router, query, map[string][]string{})

		assert.Equal(t, "{\"data\":{\"resetPassword\":null}}", response.Body.String())

	})

	t.Run("should log in as psychologist user", func(t *testing.T) {

		query := `{
			authenticateUser(input: { 
				email: "alice.johnson@psi.com.br", 
				password: "Def456$%^", 
				ipAddress: "100.100.100.100" 
			}) { 
				token 
				expiresAt 
			} 
		}`

		response := gql(router, query, map[string][]string{})

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		expiresAt := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "expiresAt")
		assert.NotEqual(t, time.Now().Unix()+res.SecondsToExpire, expiresAt)

		storedVariables["psychologist_token"] = token

	})

}
