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
		MaxAffinityNumber:          int64(5),
		SecondsLimitAvailability:   int64(2419200),
		SecondsMinimumAvailability: int64(1800),
		SecondsToCooldownReset:     int64(86400),
		SecondsToExpire:            int64(1800),
		SecondsToExpireReset:       int64(86400),
	}

	// If you need to debug the contents of the database at a specific point, insert this code:
	// db, _ := res.DatabaseUtil.GetMockedDatabases()
	// ioutil.WriteFile("./db.json", db, 0644);

	os.Setenv("PSI_BOOTSTRAP_USER", "coordinator@psi.com.br|Abc123!@#")

	router := graph.CreateServer(res)

	t.Run("should not log in with incorrect email", func(t *testing.T) {

		query := `{
			authenticateUser(input: { 
				email: "coordinator@psi.com", 
				password: "Abc123!@#"
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
				password: "123!@#Abc"
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
				password: "Abc123!@#"
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

	t.Run("should create jobrunner user only if coordinator", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
			  input: {
				email: "jobrunner@psi.com.br"
				password: "Xyz*()890"
				role: JOBRUNNER
			  }
			)
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createUserWithPassword\"]}],\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: { 
				email: "jobrunner@psi.com.br", 
				password: "Xyz*()890"
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

		storedVariables["jobrunner_token"] = token
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

	t.Run("should not create psychologist user with same email but should not warn hackers that this email is already registered", func(t *testing.T) {

		query := `mutation {
			createPsychologistUser(input: {
				email: "tom.brady@psi.com.br"
			})
		}`

		response := gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createPsychologistUser\":null}}", response.Body.String())

		// TODO: navigate GetMockedDatabases to check if there is only one tom.brady@psi.com.br
	})

	t.Run("should reset password with token sent via email", func(t *testing.T) {

		query := `mutation {
			processPendingMail
		}`

		response := gql(router, query, storedVariables["jobrunner_token"])

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

		query = `mutation {
			resetPassword(input: {
				token: "randomtoken",
				password: "Def456$$$"
			})
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"invalid token\",\"path\":[\"resetPassword\"]}],\"data\":{\"resetPassword\":null}}", response.Body.String())

		query = fmt.Sprintf(`mutation {
			resetPassword(input: {
				token: %q,
				password: "Def456$$$"
			})
		}`, match[1])

		response = gql(router, query, "")

		assert.Equal(t, "{\"data\":{\"resetPassword\":null}}", response.Body.String())

		query = fmt.Sprintf(`mutation {
			resetPassword(input: {
				token: %q,
				password: "Def456$$*"
			})
		}`, match[1])

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"invalid token\",\"path\":[\"resetPassword\"]}],\"data\":{\"resetPassword\":null}}", response.Body.String())

	})

	t.Run("should log in as psychologist user", func(t *testing.T) {

		query := `{
			authenticateUser(input: { 
				email: "tom.brady@psi.com.br", 
				password: "Def456$$$"
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

		response = gql(router, query, storedVariables["jobrunner_token"])

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
				password: "Def456$%^"
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

		assert.Equal(t, "{\"data\":{\"createPatientUser\":null}}", response.Body.String())
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

		response := gql(router, query, storedVariables["jobrunner_token"])

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
				password: "Ghi789&*("
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

		response := gql(router, query, storedVariables["jobrunner_token"])

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
				password: "Jkl012)!@"
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
			myUser {
				id
				email
			}
		}`

		response := gql(router, query, storedVariables["coordinator_token"])

		email := fastjson.GetString(response.Body.Bytes(), "data", "myUser", "email")
		assert.Equal(t, "coordinator@psi.com.br", email)
		id := fastjson.GetString(response.Body.Bytes(), "data", "myUser", "id")
		storedVariables["coordinator_id"] = id

		response = gql(router, query, storedVariables["psychologist_token"])

		email = fastjson.GetString(response.Body.Bytes(), "data", "myUser", "email")
		assert.Equal(t, "tom.brady@psi.com.br", email)
		id = fastjson.GetString(response.Body.Bytes(), "data", "myUser", "id")
		storedVariables["psychologist_id"] = id

		response = gql(router, query, storedVariables["patient_token"])

		email = fastjson.GetString(response.Body.Bytes(), "data", "myUser", "email")
		assert.Equal(t, "patrick.mahomes@psi.com.br", email)
		id = fastjson.GetString(response.Body.Bytes(), "data", "myUser", "id")
		storedVariables["patient_id"] = id

		response = gql(router, query, storedVariables["no_profile_user_token"])

		email = fastjson.GetString(response.Body.Bytes(), "data", "myUser", "email")
		assert.Equal(t, "aaron.rodgers@psi.com.br", email)
		id = fastjson.GetString(response.Body.Bytes(), "data", "myUser", "id")
		storedVariables["no_profile_user_id"] = id
	})

	t.Run("should set psychologist characteristics only if user is coordinator", func(t *testing.T) {

		query := `mutation {
			setPsychologistCharacteristics(input: [
				{
					name: "black",
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

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setPsychologistCharacteristics\"]}],\"data\":{\"setPsychologistCharacteristics\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setPsychologistCharacteristics\"]}],\"data\":{\"setPsychologistCharacteristics\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setPsychologistCharacteristics\"]}],\"data\":{\"setPsychologistCharacteristics\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setPsychologistCharacteristics\":null}}", response.Body.String())

	})

	t.Run("should get all psychologist available characteristics only if logged in", func(t *testing.T) {

		query := `{
			psychologistCharacteristics {
				name
				type
				possibleValues
			}
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"psychologistCharacteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]}]}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"psychologistCharacteristics\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"psychologistCharacteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]}]}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"psychologistCharacteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]}]}}", response.Body.String())

	})

	t.Run("should create own psychologist profile only if user is coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			upsertMyPsychologistProfile(input: {
				fullName: "Thomas Edward Patrick Brady, Jr."
				likeName: "Tom Brady",
				birthDate: 239414400,
				city: "Boston - MA"
			})
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"upsertMyPsychologistProfile\"]}],\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"upsertMyPsychologistProfile\"]}],\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		query = `mutation {
			upsertMyPsychologistProfile(input: {
				fullName: "Peyton Williams Manning",
				likeName: "Peyton Manning",
				birthDate: 196484400,
				city: "Indianapolis - IN"
			})
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		query = `{
			myPsychologistProfile {
				birthDate
				city
			}
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPsychologistProfile\"]}],\"data\":{\"myPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPsychologistProfile\"]}],\"data\":{\"myPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":239414400,\"city\":\"Boston - MA\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":196484400,\"city\":\"Indianapolis - IN\"}}}", response.Body.String())

		query = `{
			myPsychologistProfile {
				id
			}
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		storedVariables["psychologist_1_id"] = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "id")
		assert.NotEqual(t, "", storedVariables["psychologist_1_id"])

		response = gql(router, query, storedVariables["coordinator_token"])

		storedVariables["psychologist_2_id"] = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "id")
		assert.NotEqual(t, "", storedVariables["psychologist_2_id"])

	})

	t.Run("should update own psychologist profile only if user is coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			upsertMyPsychologistProfile(input: {
				fullName: "Thomas Edward Patrick Brady, Jr."
				likeName: "Tom Brady",
				birthDate: 239414400,
				city: "Tampa - FL"
			})
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"upsertMyPsychologistProfile\"]}],\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"upsertMyPsychologistProfile\"]}],\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		query = `mutation {
			upsertMyPsychologistProfile(input: {
				fullName: "Peyton Williams Manning",
				likeName: "Peyton Manning",
				birthDate: 196484400,
				city: "Denver - CO"
			})
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		query = `{
			myPsychologistProfile {
				fullName
				likeName
				birthDate
				city
			}
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPsychologistProfile\"]}],\"data\":{\"myPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPsychologistProfile\"]}],\"data\":{\"myPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":239414400,\"city\":\"Tampa - FL\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":196484400,\"city\":\"Denver - CO\"}}}", response.Body.String())

	})

	t.Run("should choose own psychologist characteristic only if user is coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			setMyPsychologistCharacteristicChoices(input: [
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
					characteristicName: "disabilities",
					selectedValues: [
						"vision",
						"locomotion"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setMyPsychologistCharacteristicChoices\"]}],\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setMyPsychologistCharacteristicChoices\"]}],\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should not select multiple psychologist characteristics if characteristic options are true or false", func(t *testing.T) {

		query := `mutation {
			setMyPsychologistCharacteristicChoices(input: [
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
					characteristicName: "disabilities",
					selectedValues: [
						"vision",
						"hearing"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"characteristic 'black' must be either true or false\",\"path\":[\"setMyPsychologistCharacteristicChoices\"]}],\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should not select multiple psychologist characteristics if many option is false", func(t *testing.T) {

		query := `mutation {
			setMyPsychologistCharacteristicChoices(input: [
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
					characteristicName: "disabilities",
					selectedValues: [
						"vision",
						"hearing"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"characteristic 'gender' needs exactly one value\",\"path\":[\"setMyPsychologistCharacteristicChoices\"]}],\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should select one or more psychologist characteristics if many option is true", func(t *testing.T) {

		query := `mutation {
			setMyPsychologistCharacteristicChoices(input: [
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
					characteristicName: "disabilities",
					selectedValues: [
						"hearing"
					]
				},
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPsychologistCharacteristicChoices(input: [
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
					characteristicName: "disabilities",
					selectedValues: [
						"locomotion"
					]
				},
			])
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should get own psychologist profile", func(t *testing.T) {

		query := `{
			myPsychologistProfile {
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

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":239414400,\"city\":\"Tampa - FL\",\"characteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"selectedValues\":[\"hearing\"]}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":196484400,\"city\":\"Denver - CO\",\"characteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"true\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"female\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"selectedValues\":[\"locomotion\"]}]}}}", response.Body.String())

	})

	t.Run("should create own patient profile if user is logged in", func(t *testing.T) {

		query := `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Patrick Lavon Mahomes II",
				likeName: "Patrick Mahomes",
				birthDate: 811296000,
				city: "Tyler - TX"
			})
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"upsertMyPatientProfile\"]}],\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Thomas Edward Patrick Brady, Jr."
				likeName: "Tom Brady",
				birthDate: 239414400,
				city: "Boston - MA"
			})
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Peyton Williams Manning",
				likeName: "Peyton Manning",
				birthDate: 196484400,
				city: "Indianapolis - IN"
			})
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				fullName
				likeName
				birthDate
				city
			}
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPatientProfile\"]}],\"data\":{\"myPatientProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Patrick Lavon Mahomes II\",\"likeName\":\"Patrick Mahomes\",\"birthDate\":811296000,\"city\":\"Tyler - TX\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":239414400,\"city\":\"Boston - MA\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":196484400,\"city\":\"Indianapolis - IN\"}}}", response.Body.String())

	})

	t.Run("should update own patient profile if user is logged in", func(t *testing.T) {

		query := `mutation {
		upsertMyPatientProfile(input: {
			fullName: "Patrick Lavon Mahomes II",
			likeName: "Patrick Mahomes",
			birthDate: 811296000,
			city: "Kansas City - MS"
		})
	}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"upsertMyPatientProfile\"]}],\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
		upsertMyPatientProfile(input: {
			fullName: "Thomas Edward Patrick Brady, Jr."
			likeName: "Tom Brady",
			birthDate: 239414400,
			city: "Tampa - FL"
		})
	}`

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
		upsertMyPatientProfile(input: {
			fullName: "Peyton Williams Manning",
			likeName: "Peyton Manning",
			birthDate: 196484400,
			city: "Denver - CO"
		})
	}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `{
		myPatientProfile {
			fullName
			likeName
			birthDate
			city
		}
	}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPatientProfile\"]}],\"data\":{\"myPatientProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Patrick Lavon Mahomes II\",\"likeName\":\"Patrick Mahomes\",\"birthDate\":811296000,\"city\":\"Kansas City - MS\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":239414400,\"city\":\"Tampa - FL\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":196484400,\"city\":\"Denver - CO\"}}}", response.Body.String())

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
			patientCharacteristics {
				name
				type
				possibleValues
			}
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"patientCharacteristics\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"patientCharacteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]}]}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"patientCharacteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]}]}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"patientCharacteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]}]}}", response.Body.String())

	})

	t.Run("should choose own patient characteristic only if user is logged in", func(t *testing.T) {

		query := `mutation {
			setMyPatientCharacteristicChoices(input: [
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

		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setMyPatientCharacteristicChoices\"]}],\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should not select multiple patient characteristics if characteristic options are true or false", func(t *testing.T) {

		query := `mutation {
			setMyPatientCharacteristicChoices(input: [
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

		assert.Equal(t, "{\"errors\":[{\"message\":\"characteristic 'has-consulted-before' must be either true or false\",\"path\":[\"setMyPatientCharacteristicChoices\"]}],\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should not select multiple patient characteristics if many option is false", func(t *testing.T) {

		query := `mutation {
			setMyPatientCharacteristicChoices(input: [
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

		assert.Equal(t, "{\"errors\":[{\"message\":\"characteristic 'gender' needs exactly one value\",\"path\":[\"setMyPatientCharacteristicChoices\"]}],\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should select one or more patient characteristics if many option is true", func(t *testing.T) {

		query := `mutation {
			setMyPatientCharacteristicChoices(input: [
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

		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientCharacteristicChoices(input: [
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

		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

	})

	t.Run("should get own patient profile", func(t *testing.T) {

		query := `{
			myPatientProfile {
				birthDate
				city
				characteristics {
					name
					type
					selectedValues
					possibleValues
				}
			}
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"birthDate\":239414400,\"city\":\"Tampa - FL\",\"characteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"false\"],\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"non-binary\"],\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"selectedValues\":[],\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"birthDate\":196484400,\"city\":\"Denver - CO\",\"characteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"true\"],\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"female\"],\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"selectedValues\":[\"hearing\",\"vision\"],\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]}]}}}", response.Body.String())

	})

	t.Run("should set own psychologist preferences if coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			setMyPsychologistPreferences(input: [
				{
					characteristicName: "disabilities",
					selectedValue: "locomotion",
					weight: 5
				},
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 6
				}
			])
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setMyPsychologistPreferences\"]}],\"data\":{\"setMyPsychologistPreferences\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setMyPsychologistPreferences\"]}],\"data\":{\"setMyPsychologistPreferences\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setMyPsychologistPreferences\":null}}", response.Body.String())

		query = `mutation {
			setMyPsychologistPreferences(input: [
				{
					characteristicName: "disabilities",
					selectedValue: "vision",
					weight: 7
				},
				{
					characteristicName: "gender",
					selectedValue: "male",
					weight: 8
				}
			])
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setMyPsychologistPreferences\":null}}", response.Body.String())

	})

	t.Run("should get own psychologist preferences if psychologist or coordinator", func(t *testing.T) {

		query := `{
			myPsychologistProfile {
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

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":239414400,\"city\":\"Tampa - FL\",\"preferences\":[{\"characteristicName\":\"gender\",\"selectedValue\":\"female\",\"weight\":6},{\"characteristicName\":\"disabilities\",\"selectedValue\":\"locomotion\",\"weight\":5}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":196484400,\"city\":\"Denver - CO\",\"preferences\":[{\"characteristicName\":\"gender\",\"selectedValue\":\"male\",\"weight\":8},{\"characteristicName\":\"disabilities\",\"selectedValue\":\"vision\",\"weight\":7}]}}}", response.Body.String())

	})

	t.Run("should set own patient preferences if logged in", func(t *testing.T) {

		query := `mutation {
			setMyPatientPreferences(input: [
				{
					characteristicName: "black",
					selectedValue: "true",
					weight: 1
				},
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 2
				}
			])
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setMyPatientPreferences\"]}],\"data\":{\"setMyPatientPreferences\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientPreferences\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientPreferences\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientPreferences(input: [
				{
					characteristicName: "black",
					selectedValue: "true",
					weight: 3
				},
				{
					characteristicName: "gender",
					selectedValue: "non-binary",
					weight: 4
				},
				{
					characteristicName: "disabilities",
					selectedValue: "hearing",
					weight: 5
				}
			])
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientPreferences\":null}}", response.Body.String())

	})

	t.Run("should get own psychologist preferences if logged in", func(t *testing.T) {

		query := `{
			myPatientProfile {
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

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPatientProfile\"]}],\"data\":{\"myPatientProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Patrick Lavon Mahomes II\",\"likeName\":\"Patrick Mahomes\",\"birthDate\":811296000,\"city\":\"Kansas City - MS\",\"preferences\":[{\"characteristicName\":\"black\",\"selectedValue\":\"true\",\"weight\":1},{\"characteristicName\":\"gender\",\"selectedValue\":\"female\",\"weight\":2}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":239414400,\"city\":\"Tampa - FL\",\"preferences\":[{\"characteristicName\":\"black\",\"selectedValue\":\"true\",\"weight\":1},{\"characteristicName\":\"gender\",\"selectedValue\":\"female\",\"weight\":2}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":196484400,\"city\":\"Denver - CO\",\"preferences\":[{\"characteristicName\":\"black\",\"selectedValue\":\"true\",\"weight\":3},{\"characteristicName\":\"gender\",\"selectedValue\":\"non-binary\",\"weight\":4},{\"characteristicName\":\"disabilities\",\"selectedValue\":\"hearing\",\"weight\":5}]}}}", response.Body.String())

	})

	t.Run("should set own availability only if coordinator or psychologist", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := fmt.Sprintf(
			`mutation {
				setMyAvailability(input: [
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

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setMyAvailability\"]}],\"data\":{\"setMyAvailability\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setMyAvailability\"]}],\"data\":{\"setMyAvailability\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setMyAvailability\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setMyAvailability\":null}}", response.Body.String())

	})

	t.Run("should not set availability if span is less than SecondsMinimumAvailability", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := fmt.Sprintf(
			`mutation {
				setMyAvailability(input: [
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

		assert.Equal(t, fmt.Sprintf("{\"errors\":[{\"message\":\"availabilities must last at least %d seconds: one availability starting at %d, ending at %d\",\"path\":[\"setMyAvailability\"]}],\"data\":{\"setMyAvailability\":null}}", res.SecondsMinimumAvailability, tomorrow+104*3600, tomorrow+104*3600+res.SecondsMinimumAvailability/2), response.Body.String())

	})

	t.Run("should not set availability if one starts in the past or ends after now + SecondsLimitAvailability", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := fmt.Sprintf(
			`mutation {
				setMyAvailability(input: [
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

		assert.Equal(t, fmt.Sprintf("{\"errors\":[{\"message\":\"availabilities must not start in the past: one availability starting at %d, current time is %d\",\"path\":[\"setMyAvailability\"]}],\"data\":{\"setMyAvailability\":null}}", tomorrow-40*3600, time.Now().Unix()), response.Body.String())

		query = fmt.Sprintf(
			`mutation {
				setMyAvailability(input: [
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

		assert.Equal(t, fmt.Sprintf("{\"errors\":[{\"message\":\"availabilities must not finish later than %d seconds from now: one availability ending at %d, limit time is %d\",\"path\":[\"setMyAvailability\"]}],\"data\":{\"setMyAvailability\":null}}", res.SecondsLimitAvailability, time.Now().Unix()+res.SecondsLimitAvailability+1, time.Now().Unix()+res.SecondsLimitAvailability), response.Body.String())

	})

	t.Run("should get own availability only if coordinator or psychologist", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := `{
			myAvailability {
				start
				end
			}
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myAvailability\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myAvailability\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, fmt.Sprintf(
			"{\"data\":{\"myAvailability\":[{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d}]}}",
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
			"{\"data\":{\"myAvailability\":[{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d}]}}",
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
				setMyAvailability(input: [
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

		assert.Equal(t, "{\"data\":{\"setMyAvailability\":null}}", response.Body.String())

		query = `{
			myAvailability {
				start
				end
			}
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, fmt.Sprintf(
			"{\"data\":{\"myAvailability\":[{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d}]}}",
			tomorrow+9*3600,
			tomorrow+17*3600,
			tomorrow+33*3600,
			tomorrow+41*3600,
		), response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, fmt.Sprintf(
			"{\"data\":{\"myAvailability\":[{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d},{\"start\":%d,\"end\":%d}]}}",
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
			createTreatment(input: {
				duration: 3600,
				price: 30,
				interval: 604800
			})
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createTreatment\"]}],\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createTreatment\"]}],\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())
	})

	t.Run("should get created treatments", func(t *testing.T) {

		query := `{
			myPsychologistProfile {
				treatments {
					id
					duration
				}
			}
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, 3600, fastjson.GetInt(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "0", "duration"))

		treatmentID := fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "0", "id")
		storedVariables["psychologist_treatment_id"] = treatmentID
		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "1", "id")
		storedVariables["psychologist_treatment_2_id"] = treatmentID
		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "2", "id")
		storedVariables["psychologist_treatment_3_id"] = treatmentID
		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "3", "id")
		storedVariables["psychologist_treatment_4_id"] = treatmentID

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, 3600, fastjson.GetInt(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "0", "duration"))

		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "0", "id")
		storedVariables["coordinator_treatment_id"] = treatmentID

		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "1", "id")
		storedVariables["coordinator_treatment_2_id"] = treatmentID
	})

	t.Run("should update treatment only if coordinator or psychologist and also owner", func(t *testing.T) {

		query := `mutation {
			updateTreatment(
				id: %q,
				input: {
					duration: %d,
					price: 30,
					interval: 604800
				}
			)
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_id"], 2700), "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"updateTreatment\"]}],\"data\":{\"updateTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_id"], 2700), storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"updateTreatment\"]}],\"data\":{\"updateTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_id"], 2700), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"updateTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["coordinator_treatment_id"], 2700), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"updateTreatment\"]}],\"data\":{\"updateTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["coordinator_treatment_id"], 3000), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"updateTreatment\":null}}", response.Body.String())
	})

	t.Run("should get updated treatments", func(t *testing.T) {

		query := `{
			myPsychologistProfile {
				treatments {
					duration
					price
					interval
				}
			}
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"treatments\":[{\"duration\":2700,\"price\":30,\"interval\":604800},{\"duration\":3600,\"price\":30,\"interval\":604800},{\"duration\":3600,\"price\":30,\"interval\":604800},{\"duration\":3600,\"price\":30,\"interval\":604800},{\"duration\":3600,\"price\":30,\"interval\":604800}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"treatments\":[{\"duration\":3000,\"price\":30,\"interval\":604800},{\"duration\":3600,\"price\":30,\"interval\":604800},{\"duration\":3600,\"price\":30,\"interval\":604800}]}}}", response.Body.String())

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
			finalizeTreatment(id: %q)
		}`, storedVariables["psychologist_treatment_id"])

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"finalizeTreatment\"]}],\"data\":{\"finalizeTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"finalizeTreatment\"]}],\"data\":{\"finalizeTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"finalizeTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"finalizeTreatment\"]}],\"data\":{\"finalizeTreatment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["no_profile_user_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"finalizeTreatment\"]}],\"data\":{\"finalizeTreatment\":null}}", response.Body.String())

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
			myPatientProfile {
				appointments {
					id
					status
				}
			}
		}`

		response := gql(router, query, storedVariables["patient_token"])

		appointmentID := fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "0", "id")
		assert.NotEqual(t, "", appointmentID)
		storedVariables["appointment_id"] = appointmentID

		assert.Equal(t, "PROPOSED", fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "0", "status"))

		query = `{
			myPsychologistProfile {
				appointments {
					id
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		appointmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "0", "id")
		assert.Equal(t, storedVariables["appointment_id"], appointmentID)

		assert.Equal(t, "PROPOSED", fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "0", "status"))

	})

	t.Run("should deny appointment", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			denyAppointment(id: %q, reason: "please reschedule to another day")
		}`, storedVariables["appointment_id"])

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"denyAppointment\"]}],\"data\":{\"denyAppointment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"denyAppointment\"]}],\"data\":{\"denyAppointment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"denyAppointment\"]}],\"data\":{\"denyAppointment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"denyAppointment\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"appointments can only be denied if their current status is PROPOSED. current status is DENIED\",\"path\":[\"denyAppointment\"]}],\"data\":{\"denyAppointment\":null}}", response.Body.String())

	})

	t.Run("should propose new appointment, confirm, and cancel by psychologist", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := `mutation {
			proposeAppointment(input: {
				treatmentId: %q,
				start: %d
			})
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_2_id"], tomorrow+15*3600), storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"proposeAppointment\":null}}", response.Body.String())

		query = `{
			myPsychologistProfile {
				appointments {
					id
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		appointmentID := fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "1", "id")
		assert.NotEqual(t, "", appointmentID)
		storedVariables["appointment_2_id"] = appointmentID

		query = fmt.Sprintf(`mutation {
			confirmAppointment(id: %q)
		}`, storedVariables["appointment_2_id"])

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"confirmAppointment\":null}}", response.Body.String())

		query = fmt.Sprintf(`mutation {
			cancelAppointmentByPsychologist(id: %q, reason: "")
		}`, storedVariables["appointment_2_id"])

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"cancelAppointmentByPsychologist\"]}],\"data\":{\"cancelAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"cancelAppointmentByPsychologist\"]}],\"data\":{\"cancelAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"cancelAppointmentByPsychologist\"]}],\"data\":{\"cancelAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"cancelAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"appointments can only be canceled if their current status is CONFIRMED. current status is CANCELED_BY_PSYCHOLOGIST\",\"path\":[\"cancelAppointmentByPsychologist\"]}],\"data\":{\"cancelAppointmentByPsychologist\":null}}", response.Body.String())
	})

	t.Run("should propose new appointment, confirm, and cancel by patient", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := `mutation {
			proposeAppointment(input: {
				treatmentId: %q,
				start: %d
			})
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_2_id"], tomorrow+13*3600), storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"proposeAppointment\":null}}", response.Body.String())

		query = `{
			myPsychologistProfile {
				appointments {
					id
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		appointmentID := fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "2", "id")
		assert.NotEqual(t, "", appointmentID)
		storedVariables["appointment_3_id"] = appointmentID

		query = fmt.Sprintf(`mutation {
			confirmAppointment(id: %q)
		}`, storedVariables["appointment_3_id"])

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"confirmAppointment\":null}}", response.Body.String())

		query = fmt.Sprintf(`mutation {
			cancelAppointmentByPatient(id: %q, reason: "")
		}`, storedVariables["appointment_3_id"])

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"cancelAppointmentByPatient\"]}],\"data\":{\"cancelAppointmentByPatient\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"cancelAppointmentByPatient\"]}],\"data\":{\"cancelAppointmentByPatient\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"cancelAppointmentByPatient\"]}],\"data\":{\"cancelAppointmentByPatient\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"cancelAppointmentByPatient\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"appointments can only be canceled if their current status is CONFIRMED. current status is CANCELED_BY_PATIENT\",\"path\":[\"cancelAppointmentByPatient\"]}],\"data\":{\"cancelAppointmentByPatient\":null}}", response.Body.String())
	})

	t.Run("should cancel future appointments when patient interrupts treatment", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := fmt.Sprintf(`mutation {
			proposeAppointment(input: {
				treatmentId: %q,
				start: %d
			})
		}`, storedVariables["psychologist_treatment_2_id"], tomorrow+9*3600)

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"proposeAppointment\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				appointments {
					id
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "PROPOSED", fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "3", "status"))

		query = fmt.Sprintf(`mutation {
			interruptTreatmentByPatient(id: %q, reason: "synergy with psychologist was not good")
		}`, storedVariables["psychologist_treatment_2_id"])

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"interruptTreatmentByPatient\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				appointments {
					id
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "CANCELED_BY_PATIENT", fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "3", "status"))

	})

	t.Run("should cancel future appointments when psychologist interrupts treatment", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q)
		}`, storedVariables["psychologist_treatment_3_id"])

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"assignTreatment\":null}}", response.Body.String())

		query = fmt.Sprintf(`mutation {
			proposeAppointment(input: {
				treatmentId: %q,
				start: %d
			})
		}`, storedVariables["psychologist_treatment_3_id"], tomorrow+9*3600)

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"proposeAppointment\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				appointments {
					id
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "PROPOSED", fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "4", "status"))

		query = fmt.Sprintf(`mutation {
			interruptTreatmentByPsychologist(id: %q, reason: "")
		}`, storedVariables["psychologist_treatment_3_id"])

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"interruptTreatmentByPsychologist\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				appointments {
					id
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "CANCELED_BY_PSYCHOLOGIST", fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "4", "status"))

	})

	t.Run("should cancel future appointments when psychologist finalizes treatment", func(t *testing.T) {

		tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Unix()

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q)
		}`, storedVariables["psychologist_treatment_4_id"])

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"assignTreatment\":null}}", response.Body.String())

		query = fmt.Sprintf(`mutation {
			proposeAppointment(input: {
				treatmentId: %q,
				start: %d
			})
		}`, storedVariables["psychologist_treatment_4_id"], tomorrow+9*3600)

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"proposeAppointment\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				appointments {
					id
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "PROPOSED", fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "5", "status"))

		query = fmt.Sprintf(`mutation {
			finalizeTreatment(id: %q)
		}`, storedVariables["psychologist_treatment_4_id"])

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"finalizeTreatment\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				appointments {
					id
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "CANCELED_BY_PSYCHOLOGIST", fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "5", "status"))

	})

	t.Run("should set and check affinities", func(t *testing.T) {

		query := `mutation {
			setMyPatientTopAffinities
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setMyPatientTopAffinities\"]}],\"data\":{\"setMyPatientTopAffinities\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientTopAffinities\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientTopAffinities\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientTopAffinities\":null}}", response.Body.String())

		query = `{
			myPatientTopAffinities {
				psychologist {
					id
				}
			}
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPatientTopAffinities\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, storedVariables["psychologist_2_id"], fastjson.GetString(response.Body.Bytes(), "data", "myPatientTopAffinities", "0", "psychologist", "id"))

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, storedVariables["psychologist_2_id"], fastjson.GetString(response.Body.Bytes(), "data", "myPatientTopAffinities", "0", "psychologist", "id"))

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, storedVariables["psychologist_1_id"], fastjson.GetString(response.Body.Bytes(), "data", "myPatientTopAffinities", "0", "psychologist", "id"))

	})

	t.Run("should set and get translations", func(t *testing.T) {

		query := `mutation {
				setTranslations(
					lang: "pt-BR"
					input: [
						{
							key: "pat-char:has-consulted-before"
							value: "Você já se consultou com um psicólogo alguma vez?"
						}
						{ key: "pat-char:has-consulted-before:true", value: "Sim" }
						{ key: "pat-char:has-consulted-before:false", value: "Não" }
						{
							key: "pat-char:gender"
							value: "Com qual desses gêneros você mais se identifica?"
						}
						{ key: "pat-char:gender:male", value: "Masculino" }
						{ key: "pat-char:gender:female", value: "Feminino" }
						{ key: "pat-char:gender:non-binary", value: "Não binário" }
						{
							key: "pat-char:disabilities"
							value: "Você possui alguma dessas deficiências?"
						}
						{ key: "pat-char:disabilities:vision", value: "Visual" }
						{ key: "pat-char:disabilities:hearing", value: "Auditiva" }
						{ key: "pat-char:disabilities:locomotion", value: "Locomotiva" }
						{
							key: "psy-pref:black:true"
							value: "Quão confortável você se sente sendo atendido por um psicólogo negro?"
						}
						{
							key: "psy-pref:black:false"
							value: "Quão confortável você se sente sendo atendido por um psicólogo que não seja negro?"
						}
						{
							key: "psy-pref:gender:male"
							value: "Quão confortável você se sente sendo atendido por um psicólogo do gênero masculino?"
						}
						{
							key: "psy-pref:gender:female"
							value: "Quão confortável você se sente sendo atendido por um psicólogo do gênero feminino?"
						}
						{
							key: "psy-pref:gender:non-binary"
							value: "Quão confortável você se sente sendo atendido por um psicólogo de gênero não binário?"
						}
						{
							key: "psy-pref:disabilities:vision"
							value: "Quão confortável você se sente sendo atendido por um psicólogo com deficiência visual?"
						}
						{
							key: "psy-pref:disabilities:hearing"
							value: "Quão confortável você se sente sendo atendido por um psicólogo com deficiência auditiva?"
						}
						{
							key: "psy-pref:disabilities:locomotion"
							value: "Quão confortável você se sente sendo atendido por um psicólogo com deficiência locomotiva?"
						}
						{
							key: "psy-char:black"
							value: "Você é negro(a)?"
						}
						{ key: "psy-char:black:true", value: "Sim" }
						{ key: "psy-char:black:false", value: "Não" }
						{
							key: "psy-char:gender"
							value: "Com qual desses gêneros você mais se identifica?"
						}
						{ key: "psy-char:gender:male", value: "Masculino" }
						{ key: "psy-char:gender:female", value: "Feminino" }
						{ key: "psy-char:gender:non-binary", value: "Não binário" }
						{
							key: "psy-char:disabilities"
							value: "Você possui alguma dessas deficiências?"
						}
						{ key: "psy-char:disabilities:vision", value: "Visual" }
						{ key: "psy-char:disabilities:hearing", value: "Auditiva" }
						{ key: "psy-char:disabilities:locomotion", value: "Locomotiva" }
						{
							key: "pat-pref:has-consulted-before:true"
							value: "Quão interessado você está em atender pacientes que já fizeram tratamento psicológico anteriormente?"
						}
						{
							key: "pat-pref:has-consulted-before:false"
							value: "Quão interessado você está em atender pacientes que nunca fizeram tratamento psicológico?"
						}
						{
							key: "pat-pref:gender:male"
							value: "Quão interessado você está em atender pacientes do gênero masculino?"
						}
						{
							key: "pat-pref:gender:female"
							value: "Quão interessado você está em atender pacientes do gênero feminino?"
						}
						{
							key: "pat-pref:gender:non-binary"
							value: "Quão interessado você está em atender pacientes de gênero não binário?"
						}
						{
							key: "pat-pref:disabilities:vision"
							value: "Quão interessado você está em atender pacientes com deficiência visual?"
						}
						{
							key: "pat-pref:disabilities:hearing"
							value: "Quão interessado você está em atender pacientes com deficiência auditiva?"
						}
						{
							key: "pat-pref:disabilities:locomotion"
							value: "Quão interessado você está em atender pacientes com deficiência locomotiva?"
						}
					]
				)
			}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setTranslations\"]}],\"data\":{\"setTranslations\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setTranslations\"]}],\"data\":{\"setTranslations\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setTranslations\"]}],\"data\":{\"setTranslations\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setTranslations\":null}}", response.Body.String())

		query = `{
			translations(lang: "pt-BR", keys: ["psy-pref:gender:female", "psy-pref:gender:male"]) {
				lang
				key
				value
			}
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"data\":{\"translations\":[{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:male\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero masculino?\"},{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:female\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero feminino?\"}]}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"translations\":[{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:male\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero masculino?\"},{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:female\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero feminino?\"}]}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"translations\":[{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:male\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero masculino?\"},{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:female\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero feminino?\"}]}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"translations\":[{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:male\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero masculino?\"},{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:female\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero feminino?\"}]}}", response.Body.String())

	})

}
