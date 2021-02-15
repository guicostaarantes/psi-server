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
		DatabaseUtil:               database.MockDatabaseUtil,
		HashUtil:                   hash.WeakBcryptHashUtil,
		IdentifierUtil:             identifier.UUIDIdentifierUtil,
		MailUtil:                   mail.MockMailUtil,
		MatchUtil:                  match.RegexpMatchUtil,
		SerializingUtil:            serializing.JSONSerializingUtil,
		TokenUtil:                  token.RngTokenUtil,
		SecondsLimitAvailability:   int64(2419200),
		SecondsMinimumAvailability: int64(1800),
		SecondsToCooldownReset:     int64(86400),
		SecondsToExpire:            int64(1800),
		SecondsToExpireReset:       int64(86400),
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

	t.Run("should create psychologist user only if coordinator", func(t *testing.T) {

		query := `mutation {
			createPsychologistUser(input: {
				email: "tom.brady@psi.com.br"
			})
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createPsychologistUser\"]}],\"data\":{\"createPsychologistUser\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createPsychologistUser\"]}],\"data\":{\"createPsychologistUser\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createPsychologistUser\"]}],\"data\":{\"createPsychologistUser\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createPsychologistUser\":null}}", response.Body.String())
	})

	t.Run("should not create psychologist user with same email", func(t *testing.T) {

		query := `mutation {
			createPsychologistUser(input: {
				email: "tom.brady@psi.com.br"
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

		regex := regexp.MustCompile("token=(?P<token>[^\"]+)")
		match := regex.FindStringSubmatch(mailBody)

		query = fmt.Sprintf(`mutation {
			resetPassword(input: {
				token: %q,
				password: "Def456$$$"
			})
		}`, match[1])

		response = gql(router, query, "")

		assert.Equal(t, "{\"data\":{\"resetPassword\":null}}", response.Body.String())

	})

	t.Run("should log in as psychologist user", func(t *testing.T) {

		query := `{
			authenticateUser(input: { 
				email: "tom.brady@psi.com.br", 
				password: "Def456$$$", 
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

	})

	t.Run("should reset password via link sent by email", func(t *testing.T) {

		query := `mutation {
			askResetPassword(email: "tom.brady@psi.com.br")
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"data\":{\"askResetPassword\":null}}", response.Body.String())

		query = `mutation {
			processPendingMail
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"processPendingMail\":null}}", response.Body.String())

		var mailBody string
		mailbox, mailboxErr := res.MailUtil.GetMockedMessages()
		assert.Equal(t, mailboxErr, nil)

		for _, mail := range *mailbox {
			if reflect.DeepEqual(mail["to"], []string{"tom.brady@psi.com.br"}) && mail["subject"] == "Redfinir senha do PSI" {
				mailBody = mail["body"].(string)
				break
			}
		}

		regex := regexp.MustCompile("token=(?P<token>[^\"]+)")
		match := regex.FindStringSubmatch(mailBody)

		query = fmt.Sprintf(`mutation {
			resetPassword(input: {
				token: %q,
				password: "Def456$%%^"
			})
		}`, match[1])

		response = gql(router, query, "")

		assert.Equal(t, "{\"data\":{\"resetPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: { 
				email: "tom.brady@psi.com.br", 
				password: "Def456$%^", 
				ipAddress: "100.100.100.100" 
			}) { 
				token 
				expiresAt 
			} 
		}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		expiresAt := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "expiresAt")
		assert.NotEqual(t, time.Now().Unix()+res.SecondsToExpire, expiresAt)

		storedVariables["psychologist_token"] = token

	})

	t.Run("should not create patient user with existing mail", func(t *testing.T) {

		query := `mutation {
			createPatientUser(input: {
				email: "coordinator@psi.com.br"
			})
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"user with same email already exists\",\"path\":[\"createPatientUser\"]}],\"data\":{\"createPatientUser\":null}}", response.Body.String())
	})

	t.Run("should create patient user", func(t *testing.T) {

		query := `mutation {
			createPatientUser(input: {
				email: "patrick.mahomes@psi.com.br"
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

		regex := regexp.MustCompile("token=(?P<token>[^\"]+)")
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

	t.Run("should create user and not assign profiles to them", func(t *testing.T) {

		query := `mutation {
			createPatientUser(input: {
				email: "aaron.rodgers@psi.com.br"
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
			if reflect.DeepEqual(mail["to"], []string{"aaron.rodgers@psi.com.br"}) && mail["subject"] == "Bem-vindo ao PSI" {
				mailBody = mail["body"].(string)
				break
			}
		}

		regex := regexp.MustCompile("token=(?P<token>[^\"]+)")
		match := regex.FindStringSubmatch(mailBody)

		query = fmt.Sprintf(`mutation {
			resetPassword(input: {
				token: %q,
				password: "Jkl012)!@"
			})
		}`, match[1])

		response = gql(router, query, "")

		assert.Equal(t, "{\"data\":{\"resetPassword\":null}}", response.Body.String())

	})

	t.Run("should log in as no-profile user", func(t *testing.T) {

		query := `{
			authenticateUser(input: { 
				email: "aaron.rodgers@psi.com.br", 
				password: "Jkl012)!@", 
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

		storedVariables["no_profile_user_token"] = token

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

		response = gql(router, query, storedVariables["no_profile_user_token"])

		email = fastjson.GetString(response.Body.Bytes(), "data", "getOwnUser", "email")
		assert.Equal(t, "aaron.rodgers@psi.com.br", email)
		id = fastjson.GetString(response.Body.Bytes(), "data", "getOwnUser", "id")
		storedVariables["no_profile_user_id"] = id
	})

	t.Run("should set psychologist characteristics only if user is coordinator", func(t *testing.T) {

		query := `mutation {
			setPsychologistCharacteristics(input: [
				{
					name: "black",
					type: BOOLEAN,
					possibleValues: []
				},
				{
					name: "gender",
					type: SINGLE,
					possibleValues: [
						"male",
						"female",
						"non-binary"
					]
				},
				{
					name: "techniques",
					type: MULTIPLE,
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
				type
				possibleValues
			}
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getPsychologistCharacteristics\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getPsychologistCharacteristics\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getPsychologistCharacteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"techniques\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"technique-1\",\"technique-2\",\"technique-3\"]}]}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getPsychologistCharacteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"techniques\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"technique-1\",\"technique-2\",\"technique-3\"]}]}}", response.Body.String())

	})

	t.Run("should create own psychologist profile only if user is coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			createOwnPsychologistProfile(input: {
				fullName: "Thomas Edward Patrick Brady, Jr."
				likeName: "Tom Brady",
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

		query = `mutation {
			createOwnPsychologistProfile(input: {
				fullName: "Peyton Williams Manning",
				likeName: "Peyton Manning",
				birthDate: 196484400,
				city: "Indianapolis - IN"
			})
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createOwnPsychologistProfile\":null}}", response.Body.String())

		query = `{
			getOwnPsychologistProfile {
				birthDate
				city
			}
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getOwnPsychologistProfile\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getOwnPsychologistProfile\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"birthDate\":239414400,\"city\":\"Boston - MA\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"birthDate\":196484400,\"city\":\"Indianapolis - IN\"}}}", response.Body.String())

	})

	t.Run("should update own psychologist profile only if user is coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			updateOwnPsychologistProfile(input: {
				fullName: "Thomas Edward Patrick Brady, Jr."
				likeName: "Tom Brady",
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
				fullName: "Peyton Williams Manning",
				likeName: "Peyton Manning",
				birthDate: 196484400,
				city: "Denver - CO"
			})
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"updateOwnPsychologistProfile\":null}}", response.Body.String())

		query = `{
			getOwnPsychologistProfile {
				fullName
				likeName
				birthDate
				city
			}
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getOwnPsychologistProfile\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getOwnPsychologistProfile\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":239414400,\"city\":\"Tampa - FL\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":196484400,\"city\":\"Denver - CO\"}}}", response.Body.String())

	})

	t.Run("should choose own psychologist characteristic only if user is coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			setOwnPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "black",
					selectedValues: [
						"false"
					]
				},
				{
					characteristicName: "gender",
					selectedValues: [
						"male"
					]
				},
				{
					characteristicName: "techniques",
					selectedValues: [
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

	t.Run("should not select multiple psychologist characteristics if characteristic options are true or false", func(t *testing.T) {

		query := `mutation {
			setOwnPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "black",
					selectedValues: [
						"true",
						"false"
					]
				},
				{
					characteristicName: "gender",
					selectedValues: [
						"male",
						"female"
					]
				},
				{
					characteristicName: "techniques",
					selectedValues: [
						"technique-1",
						"technique-3"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"characteristic 'black' must be either true or false\",\"path\":[\"setOwnPsychologistCharacteristicChoices\"]}],\"data\":{\"setOwnPsychologistCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should not select multiple psychologist characteristics if many option is false", func(t *testing.T) {

		query := `mutation {
			setOwnPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "black",
					selectedValues: [
						"false"
					]
				},
				{
					characteristicName: "gender",
					selectedValues: [
						"male",
						"female"
					]
				},
				{
					characteristicName: "techniques",
					selectedValues: [
						"technique-1",
						"technique-3"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"characteristic 'gender' needs exactly one value\",\"path\":[\"setOwnPsychologistCharacteristicChoices\"]}],\"data\":{\"setOwnPsychologistCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should select one or more psychologist characteristics if many option is true", func(t *testing.T) {

		query := `mutation {
			setOwnPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "black",
					selectedValues: [
						"false"
					]
				},
				{
					characteristicName: "gender",
					selectedValues: [
						"non-binary"
					]
				},
				{
					characteristicName: "techniques",
					selectedValues: [
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
					characteristicName: "black",
					selectedValues: [
						"true"
					]
				},
				{
					characteristicName: "gender",
					selectedValues: [
						"female"
					]
				},
				{
					characteristicName: "techniques",
					selectedValues: [
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
					type
					selectedValues
				}
			}
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"birthDate\":239414400,\"city\":\"Tampa - FL\",\"characteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"non-binary\"]},{\"name\":\"techniques\",\"type\":\"MULTIPLE\",\"selectedValues\":[\"technique-1\",\"technique-3\"]}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"birthDate\":196484400,\"city\":\"Denver - CO\",\"characteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"true\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"female\"]},{\"name\":\"techniques\",\"type\":\"MULTIPLE\",\"selectedValues\":[\"technique-2\"]}]}}}", response.Body.String())

	})

	t.Run("should create own patient profile if user is logged in", func(t *testing.T) {

		query := `mutation {
			createOwnPatientProfile(input: {
				fullName: "Patrick Lavon Mahomes II",
				likeName: "Patrick Mahomes",
				birthDate: 811296000,
				city: "Tyler - TX"
			})
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"createOwnPatientProfile\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createOwnPatientProfile\"]}],\"data\":{\"createOwnPatientProfile\":null}}", response.Body.String())

		query = `mutation {
			createOwnPatientProfile(input: {
				fullName: "Thomas Edward Patrick Brady, Jr."
				likeName: "Tom Brady",
				birthDate: 239414400,
				city: "Boston - MA"
			})
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createOwnPatientProfile\":null}}", response.Body.String())

		query = `mutation {
			createOwnPatientProfile(input: {
				fullName: "Peyton Williams Manning",
				likeName: "Peyton Manning",
				birthDate: 196484400,
				city: "Indianapolis - IN"
			})
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createOwnPatientProfile\":null}}", response.Body.String())

		query = `{
			getOwnPatientProfile {
				fullName
				likeName
				birthDate
				city
			}
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getOwnPatientProfile\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPatientProfile\":{\"fullName\":\"Patrick Lavon Mahomes II\",\"likeName\":\"Patrick Mahomes\",\"birthDate\":811296000,\"city\":\"Tyler - TX\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPatientProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":239414400,\"city\":\"Boston - MA\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPatientProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":196484400,\"city\":\"Indianapolis - IN\"}}}", response.Body.String())

	})

	t.Run("should update own patient profile if user is logged in", func(t *testing.T) {

		query := `mutation {
		updateOwnPatientProfile(input: {
			fullName: "Patrick Lavon Mahomes II",
			likeName: "Patrick Mahomes",
			birthDate: 811296000,
			city: "Kansas City - MS"
		})
	}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"updateOwnPatientProfile\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"updateOwnPatientProfile\"]}],\"data\":{\"updateOwnPatientProfile\":null}}", response.Body.String())

		query = `mutation {
		updateOwnPatientProfile(input: {
			fullName: "Thomas Edward Patrick Brady, Jr."
			likeName: "Tom Brady",
			birthDate: 239414400,
			city: "Tampa - FL"
		})
	}`

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"updateOwnPatientProfile\":null}}", response.Body.String())

		query = `mutation {
		updateOwnPatientProfile(input: {
			fullName: "Peyton Williams Manning",
			likeName: "Peyton Manning",
			birthDate: 196484400,
			city: "Denver - CO"
		})
	}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"updateOwnPatientProfile\":null}}", response.Body.String())

		query = `{
		getOwnPatientProfile {
			fullName
			likeName
			birthDate
			city
		}
	}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getOwnPatientProfile\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPatientProfile\":{\"fullName\":\"Patrick Lavon Mahomes II\",\"likeName\":\"Patrick Mahomes\",\"birthDate\":811296000,\"city\":\"Kansas City - MS\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPatientProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":239414400,\"city\":\"Tampa - FL\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPatientProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":196484400,\"city\":\"Denver - CO\"}}}", response.Body.String())

	})

	t.Run("should set patient characteristics only if user is coordinator", func(t *testing.T) {

		query := `mutation {
			setPatientCharacteristics(input: [
				{
					name: "has-consulted-before",
					type: BOOLEAN,
					possibleValues: [
						"true",
						"false"
					]
				},
				{
					name: "gender",
					type: SINGLE,
					possibleValues: [
						"male",
						"female",
						"non-binary"
					]
				},
				{
					name: "disabilities",
					type: MULTIPLE,
					possibleValues: [
						"vision",
						"hearing",
						"locomotion",
					]
				}
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setPatientCharacteristics\"]}],\"data\":{\"setPatientCharacteristics\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setPatientCharacteristics\"]}],\"data\":{\"setPatientCharacteristics\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setPatientCharacteristics\"]}],\"data\":{\"setPatientCharacteristics\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setPatientCharacteristics\":null}}", response.Body.String())

	})

	t.Run("should get all patient available characteristics if user is logged in", func(t *testing.T) {

		query := `{
			getPatientCharacteristics {
				name
				type
				possibleValues
			}
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getPatientCharacteristics\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"getPatientCharacteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]}]}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getPatientCharacteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]}]}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getPatientCharacteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]}]}}", response.Body.String())

	})

	t.Run("should choose own patient characteristic only if user is coordinator or patient", func(t *testing.T) {

		query := `mutation {
			setOwnPatientCharacteristicChoices(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValues: [
						"false"
					]
				},
				{
					characteristicName: "gender",
					selectedValues: [
						"male"
					]
				},
				{
					characteristicName: "disabilities",
					selectedValues: [
						"locomotion"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setOwnPatientCharacteristicChoices\"]}],\"data\":{\"setOwnPatientCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setOwnPatientCharacteristicChoices\"]}],\"data\":{\"setOwnPatientCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPatientCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPatientCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should not select multiple patient characteristics if characteristic options are true or false", func(t *testing.T) {

		query := `mutation {
			setOwnPatientCharacteristicChoices(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValues: [
						"true",
						"false"
					]
				},
				{
					characteristicName: "gender",
					selectedValues: [
						"male",
						"female"
					]
				},
				{
					characteristicName: "disabilities",
					selectedValues: [
						"vision"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"characteristic 'has-consulted-before' must be either true or false\",\"path\":[\"setOwnPatientCharacteristicChoices\"]}],\"data\":{\"setOwnPatientCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should not select multiple patient characteristics if many option is false", func(t *testing.T) {

		query := `mutation {
			setOwnPatientCharacteristicChoices(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValues: [
						"false"
					]
				},
				{
					characteristicName: "gender",
					selectedValues: [
						"male",
						"female"
					]
				},
				{
					characteristicName: "disabilities",
					selectedValues: [
						"vision"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"characteristic 'gender' needs exactly one value\",\"path\":[\"setOwnPatientCharacteristicChoices\"]}],\"data\":{\"setOwnPatientCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should select one or more patient characteristics if many option is true", func(t *testing.T) {

		query := `mutation {
			setOwnPatientCharacteristicChoices(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValues: [
						"false"
					]
				},
				{
					characteristicName: "gender",
					selectedValues: [
						"non-binary"
					]
				},
				{
					characteristicName: "disabilities",
					selectedValues: []
				},
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPatientCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setOwnPatientCharacteristicChoices(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValues: [
						"true"
					]
				},
				{
					characteristicName: "gender",
					selectedValues: [
						"female"
					]
				},
				{
					characteristicName: "disabilities",
					selectedValues: [
						"hearing",
						"vision"
					]
				},
			])
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPatientCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should get own patient profile", func(t *testing.T) {

		query := `{
			getOwnPatientProfile {
				birthDate
				city
				characteristics {
					name
					type
					selectedValues
				}
			}
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPatientProfile\":{\"birthDate\":239414400,\"city\":\"Tampa - FL\",\"characteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"selectedValues\":[]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"non-binary\"]},{\"name\":\"techniques\",\"type\":\"MULTIPLE\",\"selectedValues\":[]},{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"selectedValues\":[]}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPatientProfile\":{\"birthDate\":196484400,\"city\":\"Denver - CO\",\"characteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"selectedValues\":[]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"female\"]},{\"name\":\"techniques\",\"type\":\"MULTIPLE\",\"selectedValues\":[]},{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"true\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"female\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"selectedValues\":[\"hearing\",\"vision\"]}]}}}", response.Body.String())

	})

	t.Run("should set own psychologist preferences if coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			setOwnPsychologistPreferences(input: [
				{
					characteristicName: "disabilities",
					selectedValue: "locomotion",
					weight: 5
				},
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 4
				}
			])
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setOwnPsychologistPreferences\"]}],\"data\":{\"setOwnPsychologistPreferences\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setOwnPsychologistPreferences\"]}],\"data\":{\"setOwnPsychologistPreferences\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPsychologistPreferences\":null}}", response.Body.String())

		query = `mutation {
			setOwnPsychologistPreferences(input: [
				{
					characteristicName: "disabilities",
					selectedValue: "vision",
					weight: 5
				},
				{
					characteristicName: "gender",
					selectedValue: "male",
					weight: 2
				}
			])
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPsychologistPreferences\":null}}", response.Body.String())

	})

	t.Run("should get own psychologist preferences if psychologist or coordinator", func(t *testing.T) {

		query := `{
			getOwnPsychologistProfile {
				birthDate
				city
				preferences {
					characteristicName
					selectedValue
					weight
				}
			}
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"birthDate\":239414400,\"city\":\"Tampa - FL\",\"preferences\":[{\"characteristicName\":\"gender\",\"selectedValue\":\"female\",\"weight\":4},{\"characteristicName\":\"disabilities\",\"selectedValue\":\"locomotion\",\"weight\":5}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"birthDate\":196484400,\"city\":\"Denver - CO\",\"preferences\":[{\"characteristicName\":\"gender\",\"selectedValue\":\"male\",\"weight\":2},{\"characteristicName\":\"disabilities\",\"selectedValue\":\"vision\",\"weight\":5}]}}}", response.Body.String())

	})

	t.Run("should set own patient preferences if logged in", func(t *testing.T) {

		query := `mutation {
			setOwnPatientPreferences(input: [
				{
					characteristicName: "black",
					selectedValue: "true",
					weight: 5
				},
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 5
				}
			])
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setOwnPatientPreferences\"]}],\"data\":{\"setOwnPatientPreferences\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPatientPreferences\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPatientPreferences\":null}}", response.Body.String())

		query = `mutation {
			setOwnPatientPreferences(input: [
				{
					characteristicName: "black",
					selectedValue: "true",
					weight: 5
				},
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 3
				}
			])
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setOwnPatientPreferences\":null}}", response.Body.String())

	})

	t.Run("should get own psychologist preferences if logged in", func(t *testing.T) {

		query := `{
			getOwnPatientProfile {
				fullName
				likeName
				birthDate
				city
				preferences {
					characteristicName
					selectedValue
					weight
				}
			}
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getOwnPatientProfile\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPatientProfile\":{\"fullName\":\"Patrick Lavon Mahomes II\",\"likeName\":\"Patrick Mahomes\",\"birthDate\":811296000,\"city\":\"Kansas City - MS\",\"preferences\":[{\"characteristicName\":\"gender\",\"selectedValue\":\"female\",\"weight\":5}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPatientProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":239414400,\"city\":\"Tampa - FL\",\"preferences\":[{\"characteristicName\":\"gender\",\"selectedValue\":\"female\",\"weight\":5}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPatientProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":196484400,\"city\":\"Denver - CO\",\"preferences\":[{\"characteristicName\":\"gender\",\"selectedValue\":\"female\",\"weight\":3}]}}}", response.Body.String())

	})

	t.Run("should set own availability only if coordinator or psychologist", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := fmt.Sprintf(
			`mutation {
				setOwnAvailability(input: [
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d }
				])
			}`,
			tomorrow+8*3600,
			tomorrow+16*3600,
			tomorrow+32*3600,
			tomorrow+40*3600,
			tomorrow+56*3600,
			tomorrow+64*3600,
			tomorrow+80*3600,
			tomorrow+88*3600,
			tomorrow+104*3600,
			tomorrow+112*3600,
		)

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setOwnAvailability\"]}],\"data\":{\"setOwnAvailability\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setOwnAvailability\"]}],\"data\":{\"setOwnAvailability\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setOwnAvailability\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setOwnAvailability\":null}}", response.Body.String())

	})

	t.Run("should not set availability if span is less than SecondsMinimumAvailability", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := fmt.Sprintf(
			`mutation {
				setOwnAvailability(input: [
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d }
				])
			}`,
			tomorrow+8*3600,
			tomorrow+16*3600,
			tomorrow+32*3600,
			tomorrow+40*3600,
			tomorrow+56*3600,
			tomorrow+64*3600,
			tomorrow+80*3600,
			tomorrow+88*3600,
			tomorrow+104*3600,
			tomorrow+104*3600+res.SecondsMinimumAvailability/2,
		)

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, fmt.Sprintf("{\"errors\":[{\"message\":\"availabilities must last at least %d seconds: one availability starting at %d, ending at %d\",\"path\":[\"setOwnAvailability\"]}],\"data\":{\"setOwnAvailability\":null}}", res.SecondsMinimumAvailability, tomorrow+104*3600, tomorrow+104*3600+res.SecondsMinimumAvailability/2), response.Body.String())

	})

	t.Run("should not set availability if one starts in the past or ends after now + SecondsLimitAvailability", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := fmt.Sprintf(
			`mutation {
				setOwnAvailability(input: [
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d }
				])
			}`,
			tomorrow-40*3600,
			tomorrow-32*3600,
			tomorrow+32*3600,
			tomorrow+40*3600,
			tomorrow+56*3600,
			tomorrow+64*3600,
			tomorrow+80*3600,
			tomorrow+88*3600,
			tomorrow+104*3600,
			tomorrow+112*3600,
		)

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, fmt.Sprintf("{\"errors\":[{\"message\":\"availabilities must not start in the past: one availability starting at %d, current time is %d\",\"path\":[\"setOwnAvailability\"]}],\"data\":{\"setOwnAvailability\":null}}", tomorrow-40*3600, time.Now().Unix()), response.Body.String())

		query = fmt.Sprintf(
			`mutation {
				setOwnAvailability(input: [
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d },
					{ start: %d, end: %d }
				])
			}`,
			tomorrow+8*3600,
			tomorrow+16*3600,
			tomorrow+32*3600,
			tomorrow+40*3600,
			tomorrow+56*3600,
			tomorrow+64*3600,
			tomorrow+80*3600,
			tomorrow+88*3600,
			tomorrow+104*3600,
			time.Now().Unix()+res.SecondsLimitAvailability+1,
		)

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, fmt.Sprintf("{\"errors\":[{\"message\":\"availabilities must not finish later than %d seconds from now: one availability ending at %d, limit time is %d\",\"path\":[\"setOwnAvailability\"]}],\"data\":{\"setOwnAvailability\":null}}", res.SecondsLimitAvailability, time.Now().Unix()+res.SecondsLimitAvailability+1, time.Now().Unix()+res.SecondsLimitAvailability), response.Body.String())

	})

	t.Run("should get own availability only if coordinator or psychologist", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := `{
			getOwnAvailability {
				start
				end
			}
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getOwnAvailability\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"getOwnAvailability\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, fmt.Sprintf(
			"{\"data\":{\"getOwnAvailability\":[{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d}]}}",
			tomorrow+8*3600,
			tomorrow+16*3600,
			tomorrow+32*3600,
			tomorrow+40*3600,
			tomorrow+56*3600,
			tomorrow+64*3600,
			tomorrow+80*3600,
			tomorrow+88*3600,
			tomorrow+104*3600,
			tomorrow+112*3600,
		), response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, fmt.Sprintf(
			"{\"data\":{\"getOwnAvailability\":[{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d}]}}",
			tomorrow+8*3600,
			tomorrow+16*3600,
			tomorrow+32*3600,
			tomorrow+40*3600,
			tomorrow+56*3600,
			tomorrow+64*3600,
			tomorrow+80*3600,
			tomorrow+88*3600,
			tomorrow+104*3600,
			tomorrow+112*3600,
		), response.Body.String())

	})

	t.Run("should overwrite old availabilities when new ones are set", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := fmt.Sprintf(
			`mutation {
				setOwnAvailability(input: [
					{ start: %d, end: %d },
					{ start: %d, end: %d }
				])
			}`,
			tomorrow+33*3600,
			tomorrow+41*3600,
			tomorrow+9*3600,
			tomorrow+17*3600,
		)

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setOwnAvailability\":null}}", response.Body.String())

		query = `{
			getOwnAvailability {
				start
				end
			}
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, fmt.Sprintf(
			"{\"data\":{\"getOwnAvailability\":[{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d}]}}",
			tomorrow+9*3600,
			tomorrow+17*3600,
			tomorrow+33*3600,
			tomorrow+41*3600,
		), response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, fmt.Sprintf(
			"{\"data\":{\"getOwnAvailability\":[{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d}]}}",
			tomorrow+8*3600,
			tomorrow+16*3600,
			tomorrow+32*3600,
			tomorrow+40*3600,
			tomorrow+56*3600,
			tomorrow+64*3600,
			tomorrow+80*3600,
			tomorrow+88*3600,
			tomorrow+104*3600,
			tomorrow+112*3600,
		), response.Body.String())

	})

	t.Run("should create treatment only if coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			createOwnTreatment(input: {
				duration: 3600,
				price: 30,
				interval: 604800
			})
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createOwnTreatment\"]}],\"data\":{\"createOwnTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createOwnTreatment\"]}],\"data\":{\"createOwnTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createOwnTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createOwnTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createOwnTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createOwnTreatment\":null}}", response.Body.String())
	})

	t.Run("should get created treatments", func(t *testing.T) {

		query := `{
			getOwnPsychologistProfile {
				treatments {
					id
					duration
				}
			}
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, 3600, fastjson.GetInt(response.Body.Bytes(), "data", "getOwnPsychologistProfile", "treatments", "0", "duration"))

		treatmentID := fastjson.GetString(response.Body.Bytes(), "data", "getOwnPsychologistProfile", "treatments", "0", "id")
		storedVariables["psychologist_treatment_id"] = treatmentID
		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "getOwnPsychologistProfile", "treatments", "1", "id")
		storedVariables["psychologist_treatment_2_id"] = treatmentID

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, 3600, fastjson.GetInt(response.Body.Bytes(), "data", "getOwnPsychologistProfile", "treatments", "0", "duration"))

		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "getOwnPsychologistProfile", "treatments", "0", "id")
		storedVariables["coordinator_treatment_id"] = treatmentID

		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "getOwnPsychologistProfile", "treatments", "1", "id")
		storedVariables["coordinator_treatment_2_id"] = treatmentID
	})

	t.Run("should update treatment only if coordinator or psychologist and also owner", func(t *testing.T) {

		query := `mutation {
			updateOwnTreatment(
				id: %q,
				input: {
					duration: %d,
					price: 30,
					interval: 604800
				}
			)
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_id"], 2700), "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"updateOwnTreatment\"]}],\"data\":{\"updateOwnTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_id"], 2700), storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"updateOwnTreatment\"]}],\"data\":{\"updateOwnTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_id"], 2700), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"updateOwnTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["coordinator_treatment_id"], 2700), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"updateOwnTreatment\"]}],\"data\":{\"updateOwnTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["coordinator_treatment_id"], 3000), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"updateOwnTreatment\":null}}", response.Body.String())
	})

	t.Run("should get updated treatments", func(t *testing.T) {

		query := `{
			getOwnPsychologistProfile {
				treatments {
					duration
					price
					interval
				}
			}
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"treatments\":[{\"duration\":2700,\"price\":30,\"interval\":604800},{\"duration\":3600,\"price\":30,\"interval\":604800}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"getOwnPsychologistProfile\":{\"treatments\":[{\"duration\":3000,\"price\":30,\"interval\":604800},{\"duration\":3600,\"price\":30,\"interval\":604800}]}}}", response.Body.String())

	})

	t.Run("should assign a treatment to a patient profile", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q)
		}`, storedVariables["psychologist_treatment_id"])

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"assignTreatment\"]}],\"data\":{\"assignTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"assignTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"treatments can only be assigned if their current status is PENDING. current status is ACTIVE\",\"path\":[\"assignTreatment\"]}],\"data\":{\"assignTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"treatments can only be assigned if their current status is PENDING. current status is ACTIVE\",\"path\":[\"assignTreatment\"]}],\"data\":{\"assignTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["no_profile_user_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"assignTreatment\"]}],\"data\":{\"assignTreatment\":null}}", response.Body.String())

	})

	t.Run("should not assign a treatment to a patient profile if they already have an active treatment taken", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q)
		}`, storedVariables["coordinator_treatment_id"])

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"patient is already in an active treatment\",\"path\":[\"assignTreatment\"]}],\"data\":{\"assignTreatment\":null}}", response.Body.String())

	})

	t.Run("should finalize a treatment if coordinator or psychologist", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			finalizeOwnTreatment(id: %q)
		}`, storedVariables["psychologist_treatment_id"])

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"finalizeOwnTreatment\"]}],\"data\":{\"finalizeOwnTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"finalizeOwnTreatment\"]}],\"data\":{\"finalizeOwnTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"finalizeOwnTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"finalizeOwnTreatment\"]}],\"data\":{\"finalizeOwnTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["no_profile_user_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"finalizeOwnTreatment\"]}],\"data\":{\"finalizeOwnTreatment\":null}}", response.Body.String())

	})

	t.Run("should not assign a treatment that was finalized", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q)
		}`, storedVariables["psychologist_treatment_id"])

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"treatments can only be assigned if their current status is PENDING. current status is FINALIZED\",\"path\":[\"assignTreatment\"]}],\"data\":{\"assignTreatment\":null}}", response.Body.String())

	})

	t.Run("should assign a treatment to the same patient profile now that the other was finalized", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q)
		}`, storedVariables["coordinator_treatment_id"])

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"assignTreatment\":null}}", response.Body.String())

	})

	t.Run("should interrupt by patient if user owns the patient profile of the treatment", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			interruptTreatmentByPatient(id: %q, reason: "synergy with psychologist was not good")
		}`, storedVariables["coordinator_treatment_id"])

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"interruptTreatmentByPatient\"]}],\"data\":{\"interruptTreatmentByPatient\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"interruptTreatmentByPatient\"]}],\"data\":{\"interruptTreatmentByPatient\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"interruptTreatmentByPatient\"]}],\"data\":{\"interruptTreatmentByPatient\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["no_profile_user_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"interruptTreatmentByPatient\"]}],\"data\":{\"interruptTreatmentByPatient\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"interruptTreatmentByPatient\":null}}", response.Body.String())

	})

	t.Run("should interrupt by psychologist if user owns the psychologist profile of the treatment", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q)
		}`, storedVariables["coordinator_treatment_2_id"])

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"assignTreatment\":null}}", response.Body.String())

		query = fmt.Sprintf(`mutation {
			interruptTreatmentByPsychologist(id: %q, reason: "patient is not responding")
		}`, storedVariables["coordinator_treatment_2_id"])

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"interruptTreatmentByPsychologist\"]}],\"data\":{\"interruptTreatmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"interruptTreatmentByPsychologist\"]}],\"data\":{\"interruptTreatmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"interruptTreatmentByPsychologist\"]}],\"data\":{\"interruptTreatmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["no_profile_user_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"interruptTreatmentByPsychologist\"]}],\"data\":{\"interruptTreatmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"interruptTreatmentByPsychologist\":null}}", response.Body.String())

	})

	t.Run("should propose appointment", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q)
		}`, storedVariables["psychologist_treatment_2_id"])

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"assignTreatment\":null}}", response.Body.String())

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query = `mutation {
			proposeAppointment(input: {
				treatmentId: %q,
				start: %d
			})
		}`

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_id"], tomorrow+9*3600), storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"proposeAppointment\"]}],\"data\":{\"proposeAppointment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_2_id"], tomorrow+8.9*3600), storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"the psychologist is not available during the requested time slot\",\"path\":[\"proposeAppointment\"]}],\"data\":{\"proposeAppointment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_2_id"], tomorrow+16.1*3600), storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"the psychologist is not available during the requested time slot\",\"path\":[\"proposeAppointment\"]}],\"data\":{\"proposeAppointment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_2_id"], tomorrow+9*3600), storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"proposeAppointment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_2_id"], tomorrow+16*3600), storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"patient already has an appointment with status PROPOSED\",\"path\":[\"proposeAppointment\"]}],\"data\":{\"proposeAppointment\":null}}", response.Body.String())
	})

	t.Run("should get proposed appointment from both patient and psychologist", func(t *testing.T) {

		query := `{
			getOwnPatientProfile {
				appointments {
					id
					status
				}
			}
		}`

		response := gql(router, query, storedVariables["patient_token"])

		appointmentID := fastjson.GetString(response.Body.Bytes(), "data", "getOwnPatientProfile", "appointments", "0", "id")
		assert.NotEqual(t, "", appointmentID)
		storedVariables["appointment_id"] = appointmentID

		assert.Equal(t, "PROPOSED", fastjson.GetString(response.Body.Bytes(), "data", "getOwnPatientProfile", "appointments", "0", "status"))

		query = `{
			getOwnPsychologistProfile {
				appointments {
					id
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		appointmentID = fastjson.GetString(response.Body.Bytes(), "data", "getOwnPsychologistProfile", "appointments", "0", "id")
		assert.Equal(t, storedVariables["appointment_id"], appointmentID)

		assert.Equal(t, "PROPOSED", fastjson.GetString(response.Body.Bytes(), "data", "getOwnPsychologistProfile", "appointments", "0", "status"))

	})

}
