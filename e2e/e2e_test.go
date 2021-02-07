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

func gql(router *chi.Mux, query string, token string) *httptest.ResponseRecorder {

	body := fmt.Sprintf(`{"query": %q}`, query)

	request := httptest.NewRequest(http.MethodPost, "/gql", strings.NewReader(body))

	request.Header["Authorization"] = []string{token}
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

		response := gql(router, query, "")

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

		response := gql(router, query, "")

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

		response := gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		expiresAt := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "expiresAt")
		assert.NotEqual(t, time.Now().Unix()+res.SecondsToExpire, expiresAt)

		storedVariables["coordinator_token"] = token

	})

	t.Run("should create psychologist user", func(t *testing.T) {

		query := `mutation {
			createPsychologistUser(input: {
				email: "tom.brady@psi.com.br",
				firstName: "Thomas",
				lastName: "Brady"
			})
		}`

		response := gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createPsychologistUser\":null}}", response.Body.String())
	})

	t.Run("should not create psychologist user with same email", func(t *testing.T) {

		query := `mutation {
			createPsychologistUser(input: {
				email: "tom.brady@psi.com.br",
				firstName: "Thomas",
				lastName: "Brady"
			})
		}`

		response := gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"user with same email already exists\",\"path\":[\"createPsychologistUser\"]}],\"data\":{\"createPsychologistUser\":null}}", response.Body.String())
	})

	t.Run("should reset password with token sent via email", func(t *testing.T) {

		query := `mutation {
			processPendingMail
		}`

		response := gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"processPendingMail\":null}}", response.Body.String())

		var mailBody string
		mailbox, mailboxErr := res.MailUtil.GetMockedMessages()
		assert.Equal(t, mailboxErr, nil)

		for _, mail := range *mailbox {
			if reflect.DeepEqual(mail["to"], []string{"tom.brady@psi.com.br"}) && mail["subject"] == "Bem-vindo ao PSI" {
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

		response = gql(router, query, "")

		assert.Equal(t, "{\"data\":{\"resetPassword\":null}}", response.Body.String())

	})

	t.Run("should log in as psychologist user", func(t *testing.T) {

		query := `{
			authenticateUser(input: { 
				email: "tom.brady@psi.com.br", 
				password: "Def456$%^", 
				ipAddress: "100.100.100.100" 
			}) { 
				token 
				expiresAt 
			} 
		}`

		response := gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		expiresAt := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "expiresAt")
		assert.NotEqual(t, time.Now().Unix()+res.SecondsToExpire, expiresAt)

		storedVariables["psychologist_token"] = token

	})

	t.Run("should not create patient user with existing mail", func(t *testing.T) {

		query := `mutation {
			createPatientUser(input: {
				email: "coordinator@psi.com.br",
				firstName: "Patrick",
				lastName: "Mahomes"
			})
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"user with same email already exists\",\"path\":[\"createPatientUser\"]}],\"data\":{\"createPatientUser\":null}}", response.Body.String())
	})

	t.Run("should create patient user", func(t *testing.T) {

		query := `mutation {
			createPatientUser(input: {
				email: "patrick.mahomes@psi.com.br",
				firstName: "Patrick",
				lastName: "Mahomes"
			})
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"data\":{\"createPatientUser\":null}}", response.Body.String())
	})

	t.Run("should reset password with token sent via email", func(t *testing.T) {

		query := `mutation {
			processPendingMail
		}`

		response := gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"processPendingMail\":null}}", response.Body.String())

		var mailBody string
		mailbox, mailboxErr := res.MailUtil.GetMockedMessages()
		assert.Equal(t, mailboxErr, nil)

		for _, mail := range *mailbox {
			if reflect.DeepEqual(mail["to"], []string{"patrick.mahomes@psi.com.br"}) && mail["subject"] == "Bem-vindo ao PSI" {
				mailBody = mail["body"].(string)
				break
			}
		}

		regex := regexp.MustCompile("token=(?P<token>[A-Za-z0-9]{64})")
		match := regex.FindStringSubmatch(mailBody)

		query = fmt.Sprintf(`mutation {
			resetPassword(input: {
				token: %q,
				password: "Ghi789&*("
			})
		}`, match[1])

		response = gql(router, query, "")

		assert.Equal(t, "{\"data\":{\"resetPassword\":null}}", response.Body.String())

	})

	t.Run("should log in as patient user", func(t *testing.T) {

		query := `{
			authenticateUser(input: { 
				email: "patrick.mahomes@psi.com.br", 
				password: "Ghi789&*(", 
				ipAddress: "100.100.100.100" 
			}) { 
				token 
				expiresAt 
			} 
		}`

		response := gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		expiresAt := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "expiresAt")
		assert.NotEqual(t, time.Now().Unix()+res.SecondsToExpire, expiresAt)

		storedVariables["patient_token"] = token

	})

	t.Run("should get own user information", func(t *testing.T) {

		query := `{
			getOwnUser {
				id
				email
			}
		}`

		response := gql(router, query, storedVariables["coordinator_token"])

		email := fastjson.GetString(response.Body.Bytes(), "data", "getOwnUser", "email")
		assert.Equal(t, "coordinator@psi.com.br", email)
		id := fastjson.GetString(response.Body.Bytes(), "data", "getOwnUser", "id")
		storedVariables["coordinator_id"] = id

		response = gql(router, query, storedVariables["psychologist_token"])

		email = fastjson.GetString(response.Body.Bytes(), "data", "getOwnUser", "email")
		assert.Equal(t, "tom.brady@psi.com.br", email)
		id = fastjson.GetString(response.Body.Bytes(), "data", "getOwnUser", "id")
		storedVariables["psychologist_id"] = id

		response = gql(router, query, storedVariables["patient_token"])

		email = fastjson.GetString(response.Body.Bytes(), "data", "getOwnUser", "email")
		assert.Equal(t, "patrick.mahomes@psi.com.br", email)
		id = fastjson.GetString(response.Body.Bytes(), "data", "getOwnUser", "id")
		storedVariables["patient_id"] = id
	})

	t.Run("should not create psychologist user unless user is coordinator", func(t *testing.T) {

		query := `mutation {
			createPsychologistUser(input: {
				email: "tom.brady@psi.com.br",
				firstName: "Thomas",
				lastName: "Brady"
			})
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createPsychologistUser\"]}],\"data\":{\"createPsychologistUser\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createPsychologistUser\"]}],\"data\":{\"createPsychologistUser\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createPsychologistUser\"]}],\"data\":{\"createPsychologistUser\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"user with same email already exists\",\"path\":[\"createPsychologistUser\"]}],\"data\":{\"createPsychologistUser\":null}}", response.Body.String())

	})

	t.Run("should set psychologist characteristics only if user is coordinator", func(t *testing.T) {

		query := `mutation {
			setPsychologistCharacteristics(input: [
				{
					name: "skin-color",
					many: false,
					possibleValues: [
						"black",
						"white",
						"not-informed"
					]
				},
				{
					name: "techniques",
					many: true,
					possibleValues: [
						"technique-1",
						"technique-2",
						"technique-3",
					]
				}
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setPsychologistCharacteristics\"]}],\"data\":{\"setPsychologistCharacteristics\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setPsychologistCharacteristics\"]}],\"data\":{\"setPsychologistCharacteristics\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setPsychologistCharacteristics\"]}],\"data\":{\"setPsychologistCharacteristics\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setPsychologistCharacteristics\":null}}", response.Body.String())

	})

	t.Run("should get all psychologist available characteristics only if user is coordinator or psychologist", func(t *testing.T) {

		query := `{
			getPsychologistCharacteristics {
				name
				many
				possibleValues
			}
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getPsychologistCharacteristics\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getPsychologistCharacteristics\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		charName0 := fastjson.GetString(response.Body.Bytes(), "data", "getPsychologistCharacteristics", "0", "name")
		assert.Equal(t, "skin-color", charName0)

		response = gql(router, query, storedVariables["coordinator_token"])

		charName1 := fastjson.GetString(response.Body.Bytes(), "data", "getPsychologistCharacteristics", "1", "name")
		assert.Equal(t, "techniques", charName1)

		charID := fastjson.GetString(response.Body.Bytes(), "data", "getPsychologistCharacteristics", "0", "id")
		storedVariables["psychologist_characteristic_0_id"] = charID

	})

	t.Run("should create own psychologist profile only if user is coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			createOwnPsychologistProfile(input: {
				birthDate: 239414400,
				city: "Boston - MA"
			})
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createOwnPsychologistProfile\"]}],\"data\":{\"createOwnPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createOwnPsychologistProfile\"]}],\"data\":{\"createOwnPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createOwnPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createOwnPsychologistProfile\":null}}", response.Body.String())

	})

	t.Run("should update own psychologist profile only if user is coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			updateOwnPsychologistProfile(input: {
				birthDate: 239414400,
				city: "Tampa - FL"
			})
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"updateOwnPsychologistProfile\"]}],\"data\":{\"updateOwnPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"updateOwnPsychologistProfile\"]}],\"data\":{\"updateOwnPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"updateOwnPsychologistProfile\":null}}", response.Body.String())

		query = `mutation {
			updateOwnPsychologistProfile(input: {
				birthDate: 772502400,
				city: "Belo Horizonte - MG"
			})
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"updateOwnPsychologistProfile\":null}}", response.Body.String())

	})

	t.Run("should choose own psychologist characteristic only if user is coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			setOwnPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "skin-color",
					values: [
						"black"
					]
				},
				{
					characteristicName: "techniques",
					values: [
						"technique-1",
						"technique-3"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setOwnPsychologistCharacteristicChoices\"]}],\"data\":{\"setOwnPsychologistCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setOwnPsychologistCharacteristicChoices\"]}],\"data\":{\"setOwnPsychologistCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPsychologistCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPsychologistCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should not select multiple psychologist characteristics if many option is false", func(t *testing.T) {

		query := `mutation {
			setOwnPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "skin-color",
					values: [
						"white",
						"black"
					]
				},
				{
					characteristicName: "techniques",
					values: [
						"technique-1",
						"technique-3"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"characteristic 'skin-color' needs exactly one value\",\"path\":[\"setOwnPsychologistCharacteristicChoices\"]}],\"data\":{\"setOwnPsychologistCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should select one or more psychologist characteristics if many option is true", func(t *testing.T) {

		query := `mutation {
			setOwnPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "skin-color",
					values: [
						"white"
					]
				},
				{
					characteristicName: "techniques",
					values: [
						"technique-1",
						"technique-3"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPsychologistCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setOwnPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "skin-color",
					values: [
						"black"
					]
				},
				{
					characteristicName: "techniques",
					values: [
						"technique-2"
					]
				},
			])
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPsychologistCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should get own psychologist profile", func(t *testing.T) {

		query := `{
			getOwnPsychologistProfile {
				birthDate
				city
				characteristics {
					name
					many
					values
				}
			}
		}`

		response := gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"birthDate\":772502400,\"city\":\"Belo Horizonte - MG\",\"characteristics\":[{\"name\":\"skin-color\",\"many\":false,\"values\":[\"black\"]},{\"name\":\"techniques\",\"many\":true,\"values\":[\"technique-2\"]}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"birthDate\":239414400,\"city\":\"Tampa - FL\",\"characteristics\":[{\"name\":\"skin-color\",\"many\":false,\"values\":[\"white\"]},{\"name\":\"techniques\",\"many\":true,\"values\":[\"technique-1\",\"technique-3\"]}]}}}", response.Body.String())

	})

}
