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

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/go-chi/chi"
	"github.com/guicostaarantes/psi-server/graph"
	"github.com/guicostaarantes/psi-server/graph/resolvers"
	users_models "github.com/guicostaarantes/psi-server/modules/users/models"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/logging"
	"github.com/guicostaarantes/psi-server/utils/mail"
	"github.com/guicostaarantes/psi-server/utils/match"
	"github.com/guicostaarantes/psi-server/utils/orm"
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

	os.Setenv("TZ", "UTC")

	storedVariables := map[string]string{}

	loggingUtil := logging.PrintLoggingUtil{}

	hashUtil := hash.BcryptHashUtil{
		Cost:        4,
		LoggingUtil: loggingUtil,
	}

	identifierUtil := identifier.UuidIdentifierUtil{
		LoggingUtil: loggingUtil,
	}

	mailUtil := mail.FakeMailUtil{
		MockedMessages: &[]map[string]interface{}{},
	}

	matchUtil := match.RegexpMatchUtil{
		LoggingUtil: loggingUtil,
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Username("green").
		Password("blue").
		Database("red").
		Port(9876).
		Version(embeddedpostgres.V12).
		RuntimePath(fmt.Sprintf("./test-%s", time.Now().Format(time.RFC3339))).
		StartTimeout(10 * time.Second))
	err := postgres.Start()
	if err != nil {
		panic(err)
	}

	defer postgres.Stop()

	ormUtil := orm.PostgresOrmUtil{}
	ormUtil.Connect("host=localhost user=green password=blue dbname=red port=9876")

	// ormUtil := orm.SqliteOrmUtil{}
	// ormUtil.Connect(fmt.Sprintf("./test-%s.db", time.Now().Format(time.RFC3339)))

	serializingUtil := serializing.JsonSerializingUtil{
		LoggingUtil: loggingUtil,
	}

	tokenUtil := token.RngTokenUtil{
		Runes: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		Size:  8,
	}

	res := &resolvers.Resolver{
		HashUtil:                           hashUtil,
		IdentifierUtil:                     identifierUtil,
		MailUtil:                           mailUtil,
		MatchUtil:                          matchUtil,
		OrmUtil:                            &ormUtil,
		SerializingUtil:                    serializingUtil,
		TokenUtil:                          tokenUtil,
		MaxAffinityNumber:                  int64(5),
		ScheduleIntervalDuration:           time.Duration(604800) * time.Second,
		ExpireAuthTokenDuration:            time.Duration(1800) * time.Second,
		ExpireResetTokenDuration:           time.Duration(86400) * time.Second,
		InterruptTreatmentCooldownDuration: time.Duration(259200) * time.Second,
		TopAffinitiesCooldownDuration:      time.Duration(86400) * time.Second,
	}

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
		assert.NotEqual(t, time.Now().Add(res.ExpireAuthTokenDuration), expiresAt)

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
		assert.NotEqual(t, time.Now().Add(res.ExpireAuthTokenDuration), expiresAt)

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

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createPsychologistUser\"]}],\"data\":{\"createPsychologistUser\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["jobruuner_token"])

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

		users := []users_models.User{}

		ormUtil.Db().Where("email = ?", "tom.brady@psi.com.br").Limit(2).Find(&users)

		assert.Equal(t, 1, len(users))
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
		assert.NotEqual(t, time.Now().Add(res.ExpireAuthTokenDuration), expiresAt)

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
		assert.NotEqual(t, time.Now().Add(res.ExpireAuthTokenDuration), expiresAt)

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
		assert.NotEqual(t, time.Now().Add(res.ExpireAuthTokenDuration), expiresAt)

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
		assert.NotEqual(t, time.Now().Add(res.ExpireAuthTokenDuration), expiresAt)

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
				birthDate: "1977-08-03T00:00:00Z",
				city: "Boston - MA",
				bio: "Hey there, my name is Tom",
				crp: "01/123456",
				whatsapp: "(11) 2345-6789",
				instagram: "@tombrady"
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
				birthDate: "1976-03-24T00:00:00Z",
				city: "Indianapolis - IN",
				bio: "Hey there, my name is Peyton",
				crp: "01/123457",
				whatsapp: "(11) 2345-6780",
				instagram: "@peytonmanning"
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

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":\"1977-08-03T00:00:00Z\",\"city\":\"Boston - MA\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":\"1976-03-24T00:00:00Z\",\"city\":\"Indianapolis - IN\"}}}", response.Body.String())

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
				birthDate: "1977-08-03T00:00:00Z",
				city: "Tampa - FL",
				bio: "Hey there, my name is Tom",
				crp: "01/123456",
				whatsapp: "(11) 2345-6789",
				instagram: "@tombrady"
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
				birthDate: "1976-03-24T00:00:00Z",
				city: "Denver - CO",
				bio: "Hey there, my name is Peyton",
				crp: "01/123457",
				whatsapp: "(11) 2345-6780",
				instagram: "@peytonmanning"
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

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":\"1977-08-03T00:00:00Z\",\"city\":\"Tampa - FL\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":\"1976-03-24T00:00:00Z\",\"city\":\"Denver - CO\"}}}", response.Body.String())

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

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":\"1977-08-03T00:00:00Z\",\"city\":\"Tampa - FL\",\"characteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"selectedValues\":[\"hearing\"]}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":\"1976-03-24T00:00:00Z\",\"city\":\"Denver - CO\",\"characteristics\":[{\"name\":\"black\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"true\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"female\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"selectedValues\":[\"locomotion\"]}]}}}", response.Body.String())

	})

	t.Run("should create own patient profile if user is logged in", func(t *testing.T) {

		query := `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Patrick Lavon Mahomes II",
				likeName: "Patrick Mahomes",
				birthDate: "1995-09-27T00:00:00Z",
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
				birthDate: "1977-08-03T00:00:00Z",
				city: "Boston - MA"
			})
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Peyton Williams Manning",
				likeName: "Peyton Manning",
				birthDate: "1976-03-24T00:00:00Z",
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

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Patrick Lavon Mahomes II\",\"likeName\":\"Patrick Mahomes\",\"birthDate\":\"1995-09-27T00:00:00Z\",\"city\":\"Tyler - TX\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":\"1977-08-03T00:00:00Z\",\"city\":\"Boston - MA\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":\"1976-03-24T00:00:00Z\",\"city\":\"Indianapolis - IN\"}}}", response.Body.String())

	})

	t.Run("should update own patient profile if user is logged in", func(t *testing.T) {

		query := `mutation {
		upsertMyPatientProfile(input: {
			fullName: "Patrick Lavon Mahomes II",
			likeName: "Patrick Mahomes",
			birthDate: "1995-09-27T00:00:00Z",
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
			birthDate: "1977-08-03T00:00:00Z",
			city: "Tampa - FL"
		})
	}`

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
		upsertMyPatientProfile(input: {
			fullName: "Peyton Williams Manning",
			likeName: "Peyton Manning",
			birthDate: "1976-03-24T00:00:00Z",
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

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Patrick Lavon Mahomes II\",\"likeName\":\"Patrick Mahomes\",\"birthDate\":\"1995-09-27T00:00:00Z\",\"city\":\"Kansas City - MS\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":\"1977-08-03T00:00:00Z\",\"city\":\"Tampa - FL\"}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":\"1976-03-24T00:00:00Z\",\"city\":\"Denver - CO\"}}}", response.Body.String())

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
				},
				{
					name: "income",
					type: SINGLE,
					possibleValues: [
						"D",
						"C",
						"B",
						"A",
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

	t.Run("should set treatment price ranges only if user is coordinator", func(t *testing.T) {

		query := `mutation {
			setTreatmentPriceRanges(input: [
				{
					name: "free",
					minimumPrice: 0,
					maximumPrice: 0,
					eligibleFor: "D"
				},
				{
					name: "low",
					minimumPrice: 25,
					maximumPrice: 50,
					eligibleFor: "D,C"
				},
				{
					name: "medium",
					minimumPrice: 50,
					maximumPrice: 100,
					eligibleFor: "D,C,B"
				},
				{
					name: "high",
					minimumPrice: 100,
					maximumPrice: 200,
					eligibleFor: "D,C,B,A"
				}
			])
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setTreatmentPriceRanges\"]}],\"data\":{\"setTreatmentPriceRanges\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setTreatmentPriceRanges\"]}],\"data\":{\"setTreatmentPriceRanges\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"setTreatmentPriceRanges\"]}],\"data\":{\"setTreatmentPriceRanges\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"setTreatmentPriceRanges\":null}}", response.Body.String())

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

		assert.Equal(t, "{\"data\":{\"patientCharacteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]},{\"name\":\"income\",\"type\":\"SINGLE\",\"possibleValues\":[\"D\",\"C\",\"B\",\"A\"]}]}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"patientCharacteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]},{\"name\":\"income\",\"type\":\"SINGLE\",\"possibleValues\":[\"D\",\"C\",\"B\",\"A\"]}]}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"patientCharacteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]},{\"name\":\"income\",\"type\":\"SINGLE\",\"possibleValues\":[\"D\",\"C\",\"B\",\"A\"]}]}}", response.Body.String())

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
				{
					characteristicName: "income",
					selectedValues: [
						"C"
					]
				}
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
				,
				{
					characteristicName: "income",
					selectedValues: [
						"C"
					]
				}
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
				{
					characteristicName: "income",
					selectedValues: [
						"C"
					]
				}
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
				{
					characteristicName: "income",
					selectedValues: [
						"C"
					]
				}
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
				{
					characteristicName: "income",
					selectedValues: [
						"C"
					]
				}
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

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"birthDate\":\"1977-08-03T00:00:00Z\",\"city\":\"Tampa - FL\",\"characteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"false\"],\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"non-binary\"],\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"selectedValues\":[],\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]},{\"name\":\"income\",\"type\":\"SINGLE\",\"selectedValues\":[\"C\"],\"possibleValues\":[\"D\",\"C\",\"B\",\"A\"]}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"birthDate\":\"1976-03-24T00:00:00Z\",\"city\":\"Denver - CO\",\"characteristics\":[{\"name\":\"has-consulted-before\",\"type\":\"BOOLEAN\",\"selectedValues\":[\"true\"],\"possibleValues\":[\"true\",\"false\"]},{\"name\":\"gender\",\"type\":\"SINGLE\",\"selectedValues\":[\"female\"],\"possibleValues\":[\"male\",\"female\",\"non-binary\"]},{\"name\":\"disabilities\",\"type\":\"MULTIPLE\",\"selectedValues\":[\"hearing\",\"vision\"],\"possibleValues\":[\"vision\",\"hearing\",\"locomotion\"]},{\"name\":\"income\",\"type\":\"SINGLE\",\"selectedValues\":[\"C\"],\"possibleValues\":[\"D\",\"C\",\"B\",\"A\"]}]}}}", response.Body.String())

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

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":\"1977-08-03T00:00:00Z\",\"city\":\"Tampa - FL\",\"preferences\":[{\"characteristicName\":\"disabilities\",\"selectedValue\":\"locomotion\",\"weight\":5},{\"characteristicName\":\"gender\",\"selectedValue\":\"female\",\"weight\":6}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"birthDate\":\"1976-03-24T00:00:00Z\",\"city\":\"Denver - CO\",\"preferences\":[{\"characteristicName\":\"disabilities\",\"selectedValue\":\"vision\",\"weight\":7},{\"characteristicName\":\"gender\",\"selectedValue\":\"male\",\"weight\":8}]}}}", response.Body.String())

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

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Patrick Lavon Mahomes II\",\"likeName\":\"Patrick Mahomes\",\"birthDate\":\"1995-09-27T00:00:00Z\",\"city\":\"Kansas City - MS\",\"preferences\":[{\"characteristicName\":\"black\",\"selectedValue\":\"true\",\"weight\":1},{\"characteristicName\":\"gender\",\"selectedValue\":\"female\",\"weight\":2}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Thomas Edward Patrick Brady, Jr.\",\"likeName\":\"Tom Brady\",\"birthDate\":\"1977-08-03T00:00:00Z\",\"city\":\"Tampa - FL\",\"preferences\":[{\"characteristicName\":\"black\",\"selectedValue\":\"true\",\"weight\":1},{\"characteristicName\":\"gender\",\"selectedValue\":\"female\",\"weight\":2}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"fullName\":\"Peyton Williams Manning\",\"likeName\":\"Peyton Manning\",\"birthDate\":\"1976-03-24T00:00:00Z\",\"city\":\"Denver - CO\",\"preferences\":[{\"characteristicName\":\"black\",\"selectedValue\":\"true\",\"weight\":3},{\"characteristicName\":\"gender\",\"selectedValue\":\"non-binary\",\"weight\":4},{\"characteristicName\":\"disabilities\",\"selectedValue\":\"hearing\",\"weight\":5}]}}}", response.Body.String())

	})

	t.Run("should set terms only if user is coordinator", func(t *testing.T) {

		query := `mutation {
			upsertTerm(input: {
				name: "emergency",
				version: 1,
				profileType: PATIENT,
				active: true
			})
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"upsertTerm\"]}],\"data\":{\"upsertTerm\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"upsertTerm\"]}],\"data\":{\"upsertTerm\":null}}", response.Body.String())

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"upsertTerm\"]}],\"data\":{\"upsertTerm\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"upsertTerm\":null}}", response.Body.String())

		query = `mutation {
			upsertTerm(input: {
				name: "price",
				version: 1,
				profileType: PSYCHOLOGIST,
				active: true
			})
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"upsertTerm\":null}}", response.Body.String())

	})

	t.Run("should get terms", func(t *testing.T) {

		query := `{
			patientTerms {
				name
				version
				active
			}
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"patientTerms\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"patientTerms\":[{\"name\":\"emergency\",\"version\":1,\"active\":true}]}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"patientTerms\":[{\"name\":\"emergency\",\"version\":1,\"active\":true}]}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"patientTerms\":[{\"name\":\"emergency\",\"version\":1,\"active\":true}]}}", response.Body.String())

		query = `{
			psychologistTerms {
				name
				version
				active
			}
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"psychologistTerms\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"psychologistTerms\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"psychologistTerms\":[{\"name\":\"price\",\"version\":1,\"active\":true}]}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"psychologistTerms\":[{\"name\":\"price\",\"version\":1,\"active\":true}]}}", response.Body.String())

	})

	t.Run("should get own patient agreements if logged in", func(t *testing.T) {

		query := `mutation {
			upsertPatientAgreement(input: {
				termName: "emergency",
				termVersion: 1,
				agreed: true,
			})
		}`

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"upsertPatientAgreement\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				agreements {
					termName
					termVersion
				}
			}
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPatientProfile\"]}],\"data\":{\"myPatientProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"agreements\":[{\"termName\":\"emergency\",\"termVersion\":1}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"agreements\":[]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"agreements\":[]}}}", response.Body.String())

		query = `mutation {
			upsertPatientAgreement(input: {
				termName: "emergency",
				termVersion: 1,
				agreed: false,
			})
		}`

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"upsertPatientAgreement\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				agreements {
					termName
					termVersion
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"agreements\":[]}}}", response.Body.String())

	})

	t.Run("should get own psychologist agreements if logged in", func(t *testing.T) {

		query := `mutation {
			upsertPsychologistAgreement(input: {
				termName: "price",
				termVersion: 1,
				agreed: true,
			})
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"upsertPsychologistAgreement\":null}}", response.Body.String())

		query = `{
			myPsychologistProfile {
				agreements {
					termName
					termVersion
				}
			}
		}`

		response = gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPsychologistProfile\"]}],\"data\":{\"myPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPsychologistProfile\"]}],\"data\":{\"myPsychologistProfile\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"agreements\":[{\"termName\":\"price\",\"termVersion\":1}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"agreements\":[]}}}", response.Body.String())

		query = `mutation {
			upsertPsychologistAgreement(input: {
				termName: "price",
				termVersion: 1,
				agreed: false,
			})
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"upsertPsychologistAgreement\":null}}", response.Body.String())

		query = `{
			myPsychologistProfile {
				agreements {
					termName
					termVersion
				}
			}
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"agreements\":[]}}}", response.Body.String())

	})

	t.Run("should create patient 2", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
			  input: {
				email: "patient2@psi.com.br"
				password: "Xyz*()890"
				role: PATIENT
			  }
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: { 
				email: "patient2@psi.com.br", 
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
		assert.NotEqual(t, time.Now().Add(res.ExpireAuthTokenDuration), expiresAt)

		storedVariables["patient_2_token"] = token

		query = `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Patient Two"
				likeName: "Two",
				birthDate: "1977-08-03T00:00:00Z",
				city: "Boston - MA"
			})
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

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
					characteristicName: "income",
					selectedValues: [
						"D"
					]
				}
			])
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientPreferences(input: [
				{
					characteristicName: "black",
					selectedValue: "true",
					weight: 3
				},
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 4
				}
			])
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientPreferences\":null}}", response.Body.String())
	})

	t.Run("should create patient 3", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
			  input: {
				email: "patient3@psi.com.br"
				password: "Xyz*()890"
				role: PATIENT
			  }
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: { 
				email: "patient3@psi.com.br", 
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
		assert.NotEqual(t, time.Now().Add(res.ExpireAuthTokenDuration), expiresAt)

		storedVariables["patient_3_token"] = token

		query = `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Patient Three"
				likeName: "Three",
				birthDate: "1977-08-03T00:00:00Z",
				city: "Boston - MA"
			})
		}`

		response = gql(router, query, storedVariables["patient_3_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
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
					characteristicName: "income",
					selectedValues: [
						"B"
					]
				}
			])
		}`

		response = gql(router, query, storedVariables["patient_3_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientPreferences(input: [
				{
					characteristicName: "black",
					selectedValue: "true",
					weight: 6
				},
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 5
				}
			])
		}`

		response = gql(router, query, storedVariables["patient_3_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientPreferences\":null}}", response.Body.String())
	})

	t.Run("should create patient 4", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
			  input: {
				email: "patient4@psi.com.br"
				password: "Xyz*()890"
				role: PATIENT
			  }
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: { 
				email: "patient4@psi.com.br", 
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
		assert.NotEqual(t, time.Now().Add(res.ExpireAuthTokenDuration), expiresAt)

		storedVariables["patient_4_token"] = token

		query = `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Patient Four"
				likeName: "Four",
				birthDate: "1977-08-03T00:00:00Z",
				city: "Boston - MA"
			})
		}`

		response = gql(router, query, storedVariables["patient_4_token"])

		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
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
					characteristicName: "income",
					selectedValues: [
						"B"
					]
				}
			])
		}`

		response = gql(router, query, storedVariables["patient_4_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientPreferences(input: [
				{
					characteristicName: "black",
					selectedValue: "true",
					weight: 6
				},
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 5
				}
			])
		}`

		response = gql(router, query, storedVariables["patient_4_token"])

		assert.Equal(t, "{\"data\":{\"setMyPatientPreferences\":null}}", response.Body.String())
	})

	t.Run("should create treatment only if coordinator or psychologist", func(t *testing.T) {

		query := `mutation {
			createTreatment(input: {
				frequency: %d,
				phase: %d,
				duration: 3600,
				priceRangeName: %q
			})
		}`

		response := gql(router, fmt.Sprintf(query, 2, 226800, "low"), "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createTreatment\"]}],\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 226800, "low"), storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createTreatment\"]}],\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 1, res.ScheduleIntervalDuration/time.Second, "low"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"phase cannot be bigger than the schedule interval\",\"path\":[\"createTreatment\"]}],\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 1, 226800, "low"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 226800, "low"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"there is another treatment in the same period\",\"path\":[\"createTreatment\"]}],\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 228600, "low"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"there is another treatment in the same period\",\"path\":[\"createTreatment\"]}],\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 228600+res.ScheduleIntervalDuration/time.Second, "low"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"there is another treatment in the same period\",\"path\":[\"createTreatment\"]}],\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 230400, "free"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 230400+res.ScheduleIntervalDuration/time.Second, "low"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 234000, "medium"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 237600+res.ScheduleIntervalDuration/time.Second, "low"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 1, 239400, "low"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"there is another treatment in the same period\",\"path\":[\"createTreatment\"]}],\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 241200, "medium"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 244800, "low"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 248400, "free"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 226800, "low"), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 230400, "low"), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"createTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, 2, 234000, "low"), storedVariables["coordinator_token"])

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
		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "4", "id")
		storedVariables["psychologist_treatment_5_id"] = treatmentID
		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "5", "id")
		storedVariables["psychologist_treatment_6_id"] = treatmentID
		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "6", "id")
		storedVariables["psychologist_treatment_7_id"] = treatmentID
		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "7", "id")
		storedVariables["psychologist_treatment_8_id"] = treatmentID

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, 3600, fastjson.GetInt(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "0", "duration"))

		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "0", "id")
		storedVariables["coordinator_treatment_id"] = treatmentID

		treatmentID = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "treatments", "1", "id")
		storedVariables["coordinator_treatment_2_id"] = treatmentID
	})

	t.Run("should delete treatment if price range offering exists", func(t *testing.T) {

		query := `mutation {
			deleteTreatment(id: %q, priceRangeName: %q)
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_8_id"], "high"), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"deleteTreatment\"]}],\"data\":{\"deleteTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_8_id"], "high"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"price range offering not found\",\"path\":[\"deleteTreatment\"]}],\"data\":{\"deleteTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_8_id"], "free"), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"deleteTreatment\":null}}", response.Body.String())

	})

	t.Run("should update treatment only if coordinator or psychologist and also owner", func(t *testing.T) {

		query := `mutation {
			updateTreatment(
				id: %q,
				input: {
					frequency: 2,
					phase: 226800,
					duration: %d,
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
					frequency
					phase
					duration
					priceRange {
						name
					}
				}
			}
		}`

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"treatments\":[{\"frequency\":2,\"phase\":226800,\"duration\":2700,\"priceRange\":{\"name\":\"\"}},{\"frequency\":2,\"phase\":230400,\"duration\":3600,\"priceRange\":{\"name\":\"\"}},{\"frequency\":2,\"phase\":835200,\"duration\":3600,\"priceRange\":{\"name\":\"\"}},{\"frequency\":2,\"phase\":234000,\"duration\":3600,\"priceRange\":{\"name\":\"\"}},{\"frequency\":2,\"phase\":842400,\"duration\":3600,\"priceRange\":{\"name\":\"\"}},{\"frequency\":2,\"phase\":241200,\"duration\":3600,\"priceRange\":{\"name\":\"\"}},{\"frequency\":2,\"phase\":244800,\"duration\":3600,\"priceRange\":{\"name\":\"\"}}]}}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPsychologistProfile\":{\"treatments\":[{\"frequency\":2,\"phase\":226800,\"duration\":3000,\"priceRange\":{\"name\":\"\"}},{\"frequency\":2,\"phase\":230400,\"duration\":3600,\"priceRange\":{\"name\":\"\"}},{\"frequency\":2,\"phase\":234000,\"duration\":3600,\"priceRange\":{\"name\":\"\"}}]}}}", response.Body.String())

	})

	t.Run("should not assign a treatment to a patient profile that is not eligible for that price range", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q, priceRangeName: "free")
		}`, storedVariables["psychologist_treatment_id"])

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"patient is not eligible for this price range\",\"path\":[\"assignTreatment\"]}],\"data\":{\"assignTreatment\":null}}", response.Body.String())
	})

	t.Run("should assign a treatment to a patient profile", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q, priceRangeName: "low")
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
			assignTreatment(id: %q, priceRangeName: "medium")
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
			assignTreatment(id: %q, priceRangeName: "low")
		}`, storedVariables["psychologist_treatment_id"])

		response := gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"treatments can only be assigned if their current status is PENDING. current status is FINALIZED\",\"path\":[\"assignTreatment\"]}],\"data\":{\"assignTreatment\":null}}", response.Body.String())

	})

	t.Run("should assign a treatment to the same patient profile now that the other was finalized", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q, priceRangeName: "low")
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

	t.Run("should not be able to assign a new treatment due to interrupt cooldown", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q, priceRangeName: "low")
		}`, storedVariables["coordinator_treatment_2_id"])

		response := gql(router, query, storedVariables["patient_token"])

		storedVariables["interrupt_cooldown"] = time.Now().Add(res.InterruptTreatmentCooldownDuration).Format(time.RFC3339)

		assert.Equal(t, fmt.Sprintf("assign treatment is blocked for this user until %s", storedVariables["interrupt_cooldown"]), fastjson.GetString(response.Body.Bytes(), "errors", "0", "message"))

	})

	t.Run("should interrupt by psychologist if user owns the psychologist profile of the treatment", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			assignTreatment(id: %q, priceRangeName: "low")
		}`, storedVariables["coordinator_treatment_2_id"])

		response := gql(router, query, storedVariables["patient_2_token"])

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

	t.Run("should create pending appointments based on active treatments only if user is jobrunner", func(t *testing.T) {
		query := `mutation {
			assignTreatment(id: %q, priceRangeName: %q)
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_5_id"], "high"), storedVariables["patient_2_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"treatment price range offering not found\",\"path\":[\"assignTreatment\"]}],\"data\":{\"assignTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_5_id"], "low"), storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"assignTreatment\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_6_id"], "low"), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"assignTreatment\":null}}", response.Body.String())

		query = `mutation {
			createPendingAppointments
		}`

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createPendingAppointments\"]}],\"data\":{\"createPendingAppointments\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createPendingAppointments\"]}],\"data\":{\"createPendingAppointments\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["no_profile_user_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createPendingAppointments\"]}],\"data\":{\"createPendingAppointments\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"createPendingAppointments\"]}],\"data\":{\"createPendingAppointments\":null}}", response.Body.String())

		response = gql(router, query, storedVariables["jobrunner_token"])

		assert.Equal(t, "{\"data\":{\"createPendingAppointments\":null}}", response.Body.String())

		query = `query {
			myPsychologistProfile {
				appointments {
					id
					start
					status
					treatment {
						frequency
						phase
					}
				}
			}
		}`

		response = gql(router, query, storedVariables["psychologist_token"])

		appointmentStatus := fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "0", "status")
		assert.Equal(t, "CREATED", appointmentStatus)
		appointmentStatus = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "1", "status")
		assert.Equal(t, "CREATED", appointmentStatus)

		appointmentFrequency := int64(fastjson.GetInt(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "0", "treatment", "frequency"))
		appointmentPhase := int64(fastjson.GetInt(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "0", "treatment", "phase"))
		appointmentStartString := fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "0", "start")
		appointmentStart, _ := time.Parse(time.RFC3339, appointmentStartString)
		intervalDuration := int64(res.ScheduleIntervalDuration/time.Second) * appointmentFrequency
		assert.Equal(t, appointmentStart.Unix()%intervalDuration, appointmentPhase)

		appointmentFrequency = int64(fastjson.GetInt(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "1", "treatment", "frequency"))
		appointmentPhase = int64(fastjson.GetInt(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "1", "treatment", "phase"))
		appointmentStartString = fastjson.GetString(response.Body.Bytes(), "data", "myPsychologistProfile", "appointments", "1", "start")
		appointmentStart, _ = time.Parse(time.RFC3339, appointmentStartString)
		intervalDuration = int64(res.ScheduleIntervalDuration/time.Second) * appointmentFrequency
		assert.Equal(t, appointmentStart.Unix()%intervalDuration, appointmentPhase)

		query = `query {
			myPatientProfile {
				appointments {
					id
					start
					status
					treatment {
						frequency
						phase
					}
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		storedVariables["appointment_1_id"] = fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "0", "id")

		response = gql(router, query, storedVariables["coordinator_token"])

		storedVariables["appointment_2_id"] = fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "0", "id")

	})

	t.Run("should confirm appointment by patient", func(t *testing.T) {
		query := `mutation {
			confirmAppointmentByPatient(id: %q)
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"confirmAppointmentByPatient\"]}],\"data\":{\"confirmAppointmentByPatient\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"confirmAppointmentByPatient\"]}],\"data\":{\"confirmAppointmentByPatient\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"confirmAppointmentByPatient\":null}}", response.Body.String())

		query = `query {
			myPatientProfile {
				appointments {
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, "CONFIRMED_BY_PATIENT", fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "0", "status"))
	})

	t.Run("should confirm appointment by psychologist", func(t *testing.T) {
		query := `mutation {
			confirmAppointmentByPsychologist(id: %q)
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"confirmAppointmentByPsychologist\"]}],\"data\":{\"confirmAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"confirmAppointmentByPsychologist\"]}],\"data\":{\"confirmAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"confirmAppointmentByPsychologist\":null}}", response.Body.String())

		query = `query {
			myPatientProfile {
				appointments {
					status
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, "CONFIRMED_BY_BOTH", fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "0", "status"))
	})

	t.Run("should edit appointment by psychologist", func(t *testing.T) {
		query := `mutation {
			editAppointmentByPsychologist(id: %q, input: {
				start: %q
				end: %q
				priceRangeName: "medium"
				reason: "I will be on vacations this day."
			})
		}`

		start := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		end := time.Now().Add(25 * time.Hour).Format(time.RFC3339)

		response := gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"], start, end), storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"editAppointmentByPsychologist\"]}],\"data\":{\"editAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"], start, end), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"editAppointmentByPsychologist\"]}],\"data\":{\"editAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"], start, end), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"editAppointmentByPsychologist\":null}}", response.Body.String())

		query = `query {
			myPatientProfile {
				appointments {
					status
					start
					end
					reason
					priceRange {
						name
					}
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, fmt.Sprintf("{\"data\":{\"myPatientProfile\":{\"appointments\":[{\"status\":\"EDITED_BY_PSYCHOLOGIST\",\"start\":%q,\"end\":%q,\"reason\":\"I will be on vacations this day.\",\"priceRange\":{\"name\":\"medium\"}}]}}}", start, end), response.Body.String())
	})

	t.Run("should edit appointment by psychologist", func(t *testing.T) {
		query := `mutation {
			editAppointmentByPsychologist(id: %q, input: {
				start: %q
				end: %q
				priceRangeName: "medium"
				reason: "I will be on vacations this day."
			})
		}`

		start := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		end := time.Now().Add(25 * time.Hour).Format(time.RFC3339)

		response := gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"], start, end), storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"editAppointmentByPsychologist\"]}],\"data\":{\"editAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"], start, end), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"editAppointmentByPsychologist\"]}],\"data\":{\"editAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"], start, end), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"editAppointmentByPsychologist\":null}}", response.Body.String())

		query = `query {
			myPatientProfile {
				appointments {
					status
					start
					end
					reason
					priceRange {
						name
					}
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, fmt.Sprintf("{\"data\":{\"myPatientProfile\":{\"appointments\":[{\"status\":\"EDITED_BY_PSYCHOLOGIST\",\"start\":%q,\"end\":%q,\"reason\":\"I will be on vacations this day.\",\"priceRange\":{\"name\":\"medium\"}}]}}}", start, end), response.Body.String())
	})

	t.Run("should cancel appointment by patient", func(t *testing.T) {
		query := `mutation {
			cancelAppointmentByPatient(id: %q, reason: "Maybe it's better to skip this week.")
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"cancelAppointmentByPatient\"]}],\"data\":{\"cancelAppointmentByPatient\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"cancelAppointmentByPatient\"]}],\"data\":{\"cancelAppointmentByPatient\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"cancelAppointmentByPatient\":null}}", response.Body.String())

		query = `query {
			myPatientProfile {
				appointments {
					status
					reason
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"appointments\":[{\"status\":\"CANCELED_BY_PATIENT\",\"reason\":\"Maybe it's better to skip this week.\"}]}}}", response.Body.String())
	})

	t.Run("should not edit or confirm by psychologist if canceled by patient", func(t *testing.T) {
		query := `mutation {
			editAppointmentByPsychologist(id: %q, input: {
				start: %q
				end: %q
				priceRangeName: "medium"
				reason: "I will be on vacations this day."
			})
		}`

		start := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		end := time.Now().Add(25 * time.Hour).Format(time.RFC3339)

		response := gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"], start, end), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"appointment status cannot change from CANCELED_BY_PATIENT to EDITED_BY_PSYCHOLOGIST\",\"path\":[\"editAppointmentByPsychologist\"]}],\"data\":{\"editAppointmentByPsychologist\":null}}", response.Body.String())

		query = `mutation {
			confirmAppointmentByPsychologist(id: %q)
		}`

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"appointment status cannot change from CANCELED_BY_PATIENT to CONFIRMED_BY_PSYCHOLOGIST\",\"path\":[\"confirmAppointmentByPsychologist\"]}],\"data\":{\"confirmAppointmentByPsychologist\":null}}", response.Body.String())

		query = `query {
			myPatientProfile {
				appointments {
					status
					reason
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"appointments\":[{\"status\":\"CANCELED_BY_PATIENT\",\"reason\":\"Maybe it's better to skip this week.\"}]}}}", response.Body.String())
	})

	t.Run("should edit appointment by patient", func(t *testing.T) {
		query := `mutation {
			editAppointmentByPatient(id: %q, input: {
				start: %q
				reason: "I can only do it this time in that day."
			})
		}`

		start := time.Now().Add(21 * time.Hour).Format(time.RFC3339)
		end := time.Now().Add(22 * time.Hour).Format(time.RFC3339)

		response := gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"], start), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"editAppointmentByPatient\"]}],\"data\":{\"editAppointmentByPatient\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"], start), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"editAppointmentByPatient\"]}],\"data\":{\"editAppointmentByPatient\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"], start), storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"editAppointmentByPatient\":null}}", response.Body.String())

		query = `query {
			myPatientProfile {
				appointments {
					status
					start
					end
					reason
					priceRange {
						name
					}
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, fmt.Sprintf("{\"data\":{\"myPatientProfile\":{\"appointments\":[{\"status\":\"EDITED_BY_PATIENT\",\"start\":%q,\"end\":%q,\"reason\":\"I can only do it this time in that day.\",\"priceRange\":{\"name\":\"medium\"}}]}}}", start, end), response.Body.String())
	})

	t.Run("should cancel appointment by psychologist", func(t *testing.T) {
		query := `mutation {
			cancelAppointmentByPsychologist(id: %q, reason: "I had a problem and will not be able to do it this week.")
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["patient_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"cancelAppointmentByPsychologist\"]}],\"data\":{\"cancelAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["coordinator_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"resource not found\",\"path\":[\"cancelAppointmentByPsychologist\"]}],\"data\":{\"cancelAppointmentByPsychologist\":null}}", response.Body.String())

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"cancelAppointmentByPsychologist\":null}}", response.Body.String())

		query = `query {
			myPatientProfile {
				appointments {
					status
					reason
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"appointments\":[{\"status\":\"CANCELED_BY_PSYCHOLOGIST\",\"reason\":\"I had a problem and will not be able to do it this week.\"}]}}}", response.Body.String())
	})

	t.Run("should not edit or confirm by patient if canceled by patient", func(t *testing.T) {
		query := `mutation {
			editAppointmentByPatient(id: %q, input: {
				start: %q
				reason: "I will be on vacations this day."
			})
		}`

		start := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

		response := gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"], start), storedVariables["patient_2_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"appointment status cannot change from CANCELED_BY_PSYCHOLOGIST to EDITED_BY_PATIENT\",\"path\":[\"editAppointmentByPatient\"]}],\"data\":{\"editAppointmentByPatient\":null}}", response.Body.String())

		query = `mutation {
			confirmAppointmentByPatient(id: %q)
		}`

		response = gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["patient_2_token"])

		assert.Equal(t, "{\"errors\":[{\"message\":\"appointment status cannot change from CANCELED_BY_PSYCHOLOGIST to CONFIRMED_BY_PATIENT\",\"path\":[\"confirmAppointmentByPatient\"]}],\"data\":{\"confirmAppointmentByPatient\":null}}", response.Body.String())

		query = `query {
			myPatientProfile {
				appointments {
					status
					reason
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"appointments\":[{\"status\":\"CANCELED_BY_PSYCHOLOGIST\",\"reason\":\"I had a problem and will not be able to do it this week.\"}]}}}", response.Body.String())
	})

	t.Run("should cancel future appointments when patient interrupts treatment", func(t *testing.T) {

		query := `mutation {
			confirmAppointmentByPsychologist(id: %q)
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["appointment_1_id"]), storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"confirmAppointmentByPsychologist\":null}}", response.Body.String())

		query = fmt.Sprintf(`mutation {
			interruptTreatmentByPatient(id: %q, reason: "Synergy with psychologist was not good.")
		}`, storedVariables["psychologist_treatment_5_id"])

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"interruptTreatmentByPatient\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				appointments {
					status
					reason
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_2_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"appointments\":[{\"status\":\"TREATMENT_INTERRUPTED_BY_PATIENT\",\"reason\":\"Synergy with psychologist was not good.\"}]}}}", response.Body.String())

	})

	t.Run("should cancel future appointments when psychologist interrupts treatment", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			interruptTreatmentByPsychologist(id: %q, reason: "Patient has not shown in last three appointments.")
		}`, storedVariables["psychologist_treatment_6_id"])

		response := gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"interruptTreatmentByPsychologist\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				appointments {
					status
					reason
				}
			}
		}`

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"myPatientProfile\":{\"appointments\":[{\"status\":\"TREATMENT_INTERRUPTED_BY_PSYCHOLOGIST\",\"reason\":\"Patient has not shown in last three appointments.\"}]}}}", response.Body.String())

	})

	t.Run("should cancel future appointments when psychologist finalizes treatment", func(t *testing.T) {

		query := `mutation {
			assignTreatment(id: %q, priceRangeName: "medium")
		}`

		response := gql(router, fmt.Sprintf(query, storedVariables["psychologist_treatment_4_id"]), storedVariables["patient_3_token"])

		assert.Equal(t, "{\"data\":{\"assignTreatment\":null}}", response.Body.String())

		query = `mutation {
			createPendingAppointments
		}`

		response = gql(router, query, storedVariables["jobrunner_token"])

		assert.Equal(t, "{\"data\":{\"createPendingAppointments\":null}}", response.Body.String())

		query = fmt.Sprintf(`mutation {
			finalizeTreatment(id: %q)
		}`, storedVariables["psychologist_treatment_4_id"])

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"finalizeTreatment\":null}}", response.Body.String())

		query = `{
			myPatientProfile {
				appointments {
					status
					reason
				}
			}
		}`

		response = gql(router, query, storedVariables["patient_3_token"])

		assert.Equal(t, "TREATMENT_FINALIZED", fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "0", "status"))
		assert.Equal(t, "Tratamento finalizado", fastjson.GetString(response.Body.Bytes(), "data", "myPatientProfile", "appointments", "0", "reason"))

	})

	t.Run("should set and check affinities", func(t *testing.T) {

		query := `{
			myPatientTopAffinities {
				psychologist {
					id
				}
			}
		}`

		response := gql(router, query, "")

		assert.Equal(t, "{\"errors\":[{\"message\":\"forbidden\",\"path\":[\"myPatientTopAffinities\"]}],\"data\":null}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, fmt.Sprintf("assign treatment is blocked for this user until %s", storedVariables["interrupt_cooldown"]), fastjson.GetString(response.Body.Bytes(), "errors", "0", "message"))

		response = gql(router, query, storedVariables["patient_3_token"])

		assert.Equal(t, storedVariables["psychologist_1_id"], fastjson.GetString(response.Body.Bytes(), "data", "myPatientTopAffinities", "0", "psychologist", "id"))

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

		assert.Equal(t, "{\"data\":{\"translations\":[{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:female\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero feminino?\"},{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:male\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero masculino?\"}]}}", response.Body.String())

		response = gql(router, query, storedVariables["patient_token"])

		assert.Equal(t, "{\"data\":{\"translations\":[{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:female\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero feminino?\"},{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:male\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero masculino?\"}]}}", response.Body.String())

		response = gql(router, query, storedVariables["psychologist_token"])

		assert.Equal(t, "{\"data\":{\"translations\":[{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:female\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero feminino?\"},{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:male\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero masculino?\"}]}}", response.Body.String())

		response = gql(router, query, storedVariables["coordinator_token"])

		assert.Equal(t, "{\"data\":{\"translations\":[{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:female\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero feminino?\"},{\"lang\":\"pt-BR\",\"key\":\"psy-pref:gender:male\",\"value\":\"Quão confortável você se sente sendo atendido por um psicólogo do gênero masculino?\"}]}}", response.Body.String())

	})

}
