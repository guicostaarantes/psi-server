package populate

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/guicostaarantes/psi-server/graph"
	"github.com/guicostaarantes/psi-server/graph/resolvers"
	"github.com/guicostaarantes/psi-server/utils/database"
	"github.com/guicostaarantes/psi-server/utils/hash"
	"github.com/guicostaarantes/psi-server/utils/identifier"
	"github.com/guicostaarantes/psi-server/utils/logging"
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

	loggingUtil := logging.PrintLoggingUtil{}

	databaseUtil := database.MongoDatabaseUtil{
		Context:      context.Background(),
		DatabaseName: "psi_db",
		LoggingUtil:  loggingUtil,
	}

	databaseUtil.Connect("mongodb://root:pass@localhost:27017")

	hashUtil := hash.BcryptHashUtil{
		Cost:        8,
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

	serializingUtil := serializing.JsonSerializingUtil{
		LoggingUtil: loggingUtil,
	}

	tokenUtil := token.RngTokenUtil{
		Runes: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		Size:  64,
	}

	res := &resolvers.Resolver{
		DatabaseUtil:                 &databaseUtil,
		HashUtil:                     hashUtil,
		IdentifierUtil:               identifierUtil,
		MailUtil:                     mailUtil,
		MatchUtil:                    matchUtil,
		SerializingUtil:              serializingUtil,
		TokenUtil:                    tokenUtil,
		MaxAffinityNumber:            int64(5),
		SecondsToCooldownReset:       int64(86400),
		SecondsToExpire:              int64(1800),
		SecondsToExpireReset:         int64(86400),
		TopAffinitiesCooldownSeconds: int64(86400),
	}

	// If you need to debug the contents of the database at a specific point, insert this code:
	// db, _ := res.DatabaseUtil.GetMockedDatabases()
	// ioutil.WriteFile("./db.json", db, 0644);

	os.Setenv("PSI_BOOTSTRAP_USER", "coordinator@psi.com.br|Abc123!@#")

	router := graph.CreateServer(res)

	t.Run("should log in as bootstrap coordinator", func(t *testing.T) {

		query := `{
			authenticateUser(input: {
				email: "coordinator@psi.com.br",
				password: "Abc123!@#"
			}) {
				token
			}
		}`

		response := gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["coordinator_token"] = token

	})

	t.Run("should create characteristics", func(t *testing.T) {

		query := `mutation {
			setPatientCharacteristics(input: [
				{
					name: "has-consulted-before",
					type: BOOLEAN,
					possibleValues: [
						"true",
						"false",
					]
				},
				{
					name: "gender",
					type: SINGLE,
					possibleValues: [
						"female",
						"male",
						"non-binary",
					]
				},
				{
					name: "age",
					type: SINGLE,
					possibleValues: [
						"child",
						"teen",
						"young-adult",
						"adult",
						"elderly",
					]
				},
				{
					name: "lgbtqiaplus",
					type: BOOLEAN,
					possibleValues: [
						"true",
						"false",
					]
				},
				{
					name: "skin-tone",
					type: SINGLE,
					possibleValues: [
						"black",
						"red",
						"yellow",
						"white",
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

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"setPatientCharacteristics\":null}}", response.Body.String())

		query = `mutation {
			setPsychologistCharacteristics(input: [
				{
					name: "gender",
					type: SINGLE,
					possibleValues: [
						"female",
						"male",
						"non-binary",
					]
				},
				{
					name: "lgbtqiaplus",
					type: BOOLEAN,
					possibleValues: [
						"true",
						"false",
					]
				},
				{
					name: "skin-tone",
					type: SINGLE,
					possibleValues: [
						"black",
						"red",
						"yellow",
						"white",
					]
				},
				{
					name: "sign-language",
					type: BOOLEAN,
					possibleValues: [
						"true",
						"false",
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

		response = gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"setPsychologistCharacteristics\":null}}", response.Body.String())

		query = `mutation {
			setTranslations(
				lang: "pt-BR"
				input: [
					{
						key: "pat-char:has-consulted-before"
						value: "Você já se consultou com um psicólogo alguma vez?"
					}
					{ key: "pat-char:has-consulted-before:true", value: "Já" }
					{ key: "pat-char:has-consulted-before:false", value: "Nunca" }
					{
						key: "pat-char:gender"
						value: "Com qual desses gêneros você mais se identifica?"
					}
					{ key: "pat-char:gender:female", value: "Feminino" }
					{ key: "pat-char:gender:male", value: "Masculino" }
					{ key: "pat-char:gender:non-binary", value: "Não binário" }
					{
						key: "pat-char:age"
						value: "Em qual dessas faixas etárias você se encaixa?"
					}
					{ key: "pat-char:age:child", value: "Entre 0 e 13 anos" }
					{ key: "pat-char:age:teen", value: "Entre 14 e 18 anos" }
					{ key: "pat-char:age:young-adult", value: "Entre 19 e 30 anos" }
					{ key: "pat-char:age:adult", value: "Entre 31 e 50 anos" }
					{ key: "pat-char:age:elderly", value: "Mais de 50 anos" }
					{
						key: "pat-char:lgbtqiaplus"
						value: "Você é uma pessoa LBGTQIA+?"
					}
					{ key: "pat-char:lgbtqiaplus:true", value: "Sim" }
					{ key: "pat-char:lgbtqiaplus:false", value: "Não" }
					{
						key: "pat-char:skin-tone"
						value: "Qual desses tons de pele mais se aproxima ao seu?"
					}
					{ key: "pat-char:skin-tone:black", value: "Pele negra" }
					{ key: "pat-char:skin-tone:red", value: "Pele vermelha" }
					{ key: "pat-char:skin-tone:yellow", value: "Pele amarela" }
					{ key: "pat-char:skin-tone:white", value: "Pele branca" }
					{
						key: "pat-char:disabilities"
						value: "Você possui alguma dessas deficiências?"
					}
					{ key: "pat-char:disabilities:vision", value: "Visual" }
					{ key: "pat-char:disabilities:hearing", value: "Auditiva" }
					{ key: "pat-char:disabilities:locomotion", value: "Locomotiva" }
					{
						key: "pat-pref:has-consulted-before:true"
						value: "Quão interessado você está em atender pacientes que já fizeram tratamento psicológico anteriormente?"
					}
					{
						key: "pat-pref:has-consulted-before:false"
						value: "Quão interessado você está em atender pacientes que nunca fizeram tratamento psicológico?"
					}
					{
						key: "pat-pref:gender:female"
						value: "Quão interessado você está em atender pacientes de gênero feminino?"
					}
					{
						key: "pat-pref:gender:male"
						value: "Quão interessado você está em atender pacientes de gênero masculino?"
					}
					{
						key: "pat-pref:gender:non-binary"
						value: "Quão interessado você está em atender pacientes de gênero não binário?"
					}
					{
						key: "pat-pref:age:child"
						value: "Quão interessado você está em atender pacientes entre 0 e 13 anos?"
					}
					{
						key: "pat-pref:age:teen"
						value: "Quão interessado você está em atender pacientes entre 14 e 18 anos?"
					}
					{
						key: "pat-pref:age:young-adult"
						value: "Quão interessado você está em atender pacientes entre 19 e 30 anos?"
					}
					{
						key: "pat-pref:age:adult"
						value: "Quão interessado você está em atender pacientes entre 31 e 50 anos?"
					}
					{
						key: "pat-pref:age:elderly"
						value: "Quão interessado você está em atender pacientes com mais de 50 anos?"
					}
					{ 
						key: "pat-pref:lgbtqiaplus:true"
						value: "Quão interessado você está em atender pacientes LGBTQIA+?"
					}
					{ 
						key: "pat-pref:lgbtqiaplus:false"
						value: "Quão interessado você está em atender pacientes não LGBTQIA+?"
					}
					{
						key: "pat-pref:skin-tone:black"
						value: "Quão interessado você está em atender pacientes de pele negra?"
					}
					{
						key: "pat-pref:skin-tone:red"
						value: "Quão interessado você está em atender pacientes de pele vermelha?"
					}
					{
						key: "pat-pref:skin-tone:yellow"
						value: "Quão interessado você está em atender pacientes de pele amarela?"
					}
					{
						key: "pat-pref:skin-tone:white"
						value: "Quão interessado você está em atender pacientes de pele branca?"
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
					{
						key: "psy-char:gender"
						value: "Com qual desses gêneros você mais se identifica?"
					}
					{ key: "psy-char:gender:male", value: "Masculino" }
					{ key: "psy-char:gender:female", value: "Feminino" }
					{ key: "psy-char:gender:non-binary", value: "Não binário" }
					{
						key: "psy-char:lgbtqiaplus"
						value: "Você é uma pessoa LBGTQIA+?"
					}
					{ key: "psy-char:lgbtqiaplus:true", value: "Sim" }
					{ key: "psy-char:lgbtqiaplus:false", value: "Não" }
					{
						key: "psy-char:skin-tone"
						value: "Qual desses tons de pele mais se aproxima ao seu?"
					}
					{ key: "psy-char:skin-tone:black", value: "Pele negra" }
					{ key: "psy-char:skin-tone:red", value: "Pele vermelha" }
					{ key: "psy-char:skin-tone:yellow", value: "Pele amarela" }
					{ key: "psy-char:skin-tone:white", value: "Pele branca" }
					{
						key: "psy-char:sign-language"
						value: "Você é capaz de atender a paciente utilizando LIBRAS (Linguagem Brasileira de Sinais)?"
					}
					{ key: "psy-char:sign-language:true", value: "Sim" }
					{ key: "psy-char:sign-language:false", value: "Não" }
					{
						key: "psy-char:disabilities"
						value: "Você possui alguma dessas deficiências?"
					}
					{ key: "psy-char:disabilities:vision", value: "Visual" }
					{ key: "psy-char:disabilities:hearing", value: "Auditiva" }
					{ key: "psy-char:disabilities:locomotion", value: "Locomotiva" }
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
						key: "psy-pref:lgbtqiaplus:true"
						value: "Quão confortável você se sente sendo atendido por um psicólogo LGBTQIA+?"
					}
					{ 
						key: "psy-pref:lgbtqiaplus:false"
						value: "Quão confortável você se sente sendo atendido por um psicólogo que não é LGBTQIA+?"
					}
					{
						key: "psy-pref:skin-tone:black"
						value: "Quão confortável você se sente sendo atendido por um psicólogo de pele negra?"
					}
					{
						key: "psy-pref:skin-tone:red"
						value: "Quão confortável você se sente sendo atendido por um psicólogo de pele vermelha?"
					}
					{
						key: "psy-pref:skin-tone:yellow"
						value: "Quão confortável você se sente sendo atendido por um psicólogo de pele amarela?"
					}
					{
						key: "psy-pref:skin-tone:white"
						value: "Quão confortável você se sente sendo atendido por um psicólogo de pele branca?"
					}
					{
						key: "psy-pref:sign-language:true"
						value: "Quão confortável você se sente sendo atendido por um psicólogo fluente em linguagem de sinais?"
					}
					{
						key: "psy-pref:sign-language:false"
						value: "Quão confortável você se sente sendo atendido por um psicólogo que não é fluente em linguagem de sinais?"
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
				]
			)
		}`

		response = gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"setTranslations\":null}}", response.Body.String())

	})

	t.Run("should create jobrunner", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
			  input: {
				email: "jobrunner@psi.com.br"
				password: "Xyz*()890"
				role: JOBRUNNER
			  }
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: { 
				email: "jobrunner@psi.com.br", 
				password: "Xyz*()890"
			}) { 
				token
			} 
		}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["jobrunner_token"] = token

	})

	t.Run("should create psychologist 1", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
			  input: {
				email: "ana.duarte@psi.com.br"
				password: "Abc123!@#"
				role: PSYCHOLOGIST
			  }
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: {
				email: "ana.duarte@psi.com.br"
				password: "Abc123!@#"
			}) {
				token
			}
		}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["psychologist_1_token"] = token

		query = `mutation {
			upsertMyPsychologistProfile(input: {
				fullName: "Ana Duarte"
				likeName: "Dra. Ana",
				birthDate: 239414400,
				city: "Belo Horizonte - MG"
			})
		}`

		response = gql(router, query, storedVariables["psychologist_1_token"])
		assert.Equal(t, "{\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		query = `mutation {
			setMyPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "gender",
					selectedValues: ["female"]
				},
				{
					characteristicName: "lgbtqiaplus",
					selectedValues: ["true"]
				},
				{
					characteristicName: "skin-tone",
					selectedValues: ["white"]
				},
				{
					characteristicName: "sign-language",
					selectedValues: ["false"]
				},
				{
					characteristicName: "disabilities",
					selectedValues: []
				},
			])
		}`

		response = gql(router, query, storedVariables["psychologist_1_token"])
		assert.Equal(t, "{\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPsychologistPreferences(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValue: "true",
					weight: -1
				}
				{
					characteristicName: "has-consulted-before",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 3
				}
				{
					characteristicName: "gender",
					selectedValue: "male",
					weight: 0
				}
				{
					characteristicName: "gender",
					selectedValue: "non-binary",
					weight: 1
				}
				{
					characteristicName: "age",
					selectedValue: "child",
					weight: 1
				}
				{
					characteristicName: "age",
					selectedValue: "teen",
					weight: 3
				}
				{
					characteristicName: "age",
					selectedValue: "young-adult",
					weight: 0
				}
				{
					characteristicName: "age",
					selectedValue: "adult",
					weight: 0
				}
				{
					characteristicName: "age",
					selectedValue: "elderly",
					weight: -3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "true",
					weight: 3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "false",
					weight: 0
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "black",
					weight: 1
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "red",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "yellow",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "white",
					weight: 1
				}
				{
					characteristicName: "disabilities",
					selectedValue: "vision",
					weight: 1
				}
				{
					characteristicName: "disabilities",
					selectedValue: "hearing",
					weight: -3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "locomotion",
					weight: 0
				}
			])
		}`

		response = gql(router, query, storedVariables["psychologist_1_token"])
		assert.Equal(t, "{\"data\":{\"setMyPsychologistPreferences\":null}}", response.Body.String())

	})

	t.Run("should create psychologist 2", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
				input: {
				email: "lucio.fonseca@psi.com.br"
				password: "Abc123!@#"
				role: PSYCHOLOGIST
				}
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: {
				email: "lucio.fonseca@psi.com.br"
				password: "Abc123!@#"
			}) {
				token
			}
		}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["psychologist_2_token"] = token

		query = `mutation {
			upsertMyPsychologistProfile(input: {
				fullName: "Lúcio Fonseca"
				likeName: "Dr. Lúcio",
				birthDate: -199918800,
				city: "Curitiba - PR"
			})
		}`

		response = gql(router, query, storedVariables["psychologist_2_token"])
		assert.Equal(t, "{\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		query = `mutation {
			setMyPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "gender",
					selectedValues: ["male"]
				},
				{
					characteristicName: "lgbtqiaplus",
					selectedValues: ["false"]
				},
				{
					characteristicName: "skin-tone",
					selectedValues: ["yellow"]
				},
				{
					characteristicName: "sign-language",
					selectedValues: ["true"]
				},
				{
					characteristicName: "disabilities",
					selectedValues: ["hearing"]
				},
			])
		}`

		response = gql(router, query, storedVariables["psychologist_2_token"])
		assert.Equal(t, "{\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPsychologistPreferences(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValue: "true",
					weight: 0
				}
				{
					characteristicName: "has-consulted-before",
					selectedValue: "false",
					weight: -3
				}
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: -3
				}
				{
					characteristicName: "gender",
					selectedValue: "male",
					weight: 1
				}
				{
					characteristicName: "gender",
					selectedValue: "non-binary",
					weight: 3
				}
				{
					characteristicName: "age",
					selectedValue: "child",
					weight: -3
				}
				{
					characteristicName: "age",
					selectedValue: "teen",
					weight: -3
				}
				{
					characteristicName: "age",
					selectedValue: "young-adult",
					weight: 1
				}
				{
					characteristicName: "age",
					selectedValue: "adult",
					weight: 0
				}
				{
					characteristicName: "age",
					selectedValue: "elderly",
					weight: 3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "true",
					weight: -3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "black",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "red",
					weight: -3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "yellow",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "white",
					weight: -1
				}
				{
					characteristicName: "disabilities",
					selectedValue: "vision",
					weight: -1
				}
				{
					characteristicName: "disabilities",
					selectedValue: "hearing",
					weight: 3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "locomotion",
					weight: -1
				}
			])
		}`

		response = gql(router, query, storedVariables["psychologist_2_token"])
		assert.Equal(t, "{\"data\":{\"setMyPsychologistPreferences\":null}}", response.Body.String())

	})

	t.Run("should create psychologist 3", func(t *testing.T) {

		query := `mutation {
				createUserWithPassword(
					input: {
					email: "laura.carvalho@psi.com.br"
					password: "Abc123!@#"
					role: PSYCHOLOGIST
					}
				)
			}`

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
				authenticateUser(input: {
					email: "laura.carvalho@psi.com.br"
					password: "Abc123!@#"
				}) {
					token
				}
			}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["psychologist_3_token"] = token

		query = `mutation {
				upsertMyPsychologistProfile(input: {
					fullName: "Laura Carvalho"
					likeName: "Dra. Laura",
					birthDate: 519966000,
					city: "Campo Grande - MS"
				})
			}`

		response = gql(router, query, storedVariables["psychologist_3_token"])
		assert.Equal(t, "{\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		query = `mutation {
				setMyPsychologistCharacteristicChoices(input: [
					{
						characteristicName: "gender",
						selectedValues: ["non-binary"]
					},
					{
						characteristicName: "lgbtqiaplus",
						selectedValues: ["true"]
					},
					{
						characteristicName: "skin-tone",
						selectedValues: ["black"]
					},
					{
						characteristicName: "sign-language",
						selectedValues: ["false"]
					},
					{
						characteristicName: "disabilities",
						selectedValues: ["locomotion"]
					},
				])
			}`

		response = gql(router, query, storedVariables["psychologist_3_token"])
		assert.Equal(t, "{\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
				setMyPsychologistPreferences(input: [
					{
						characteristicName: "has-consulted-before",
						selectedValue: "true",
						weight: 1
					}
					{
						characteristicName: "has-consulted-before",
						selectedValue: "false",
						weight: -1
					}
					{
						characteristicName: "gender",
						selectedValue: "female",
						weight: 3
					}
					{
						characteristicName: "gender",
						selectedValue: "male",
						weight: 1
					}
					{
						characteristicName: "gender",
						selectedValue: "non-binary",
						weight: -3
					}
					{
						characteristicName: "age",
						selectedValue: "child",
						weight: 3
					}
					{
						characteristicName: "age",
						selectedValue: "teen",
						weight: 3
					}
					{
						characteristicName: "age",
						selectedValue: "young-adult",
						weight: 0
					}
					{
						characteristicName: "age",
						selectedValue: "adult",
						weight: 0
					}
					{
						characteristicName: "age",
						selectedValue: "elderly",
						weight: 3
					}
					{
						characteristicName: "lgbtqiaplus",
						selectedValue: "true",
						weight: 3
					}
					{
						characteristicName: "lgbtqiaplus",
						selectedValue: "false",
						weight: 0
					}
					{
						characteristicName: "skin-tone",
						selectedValue: "black",
						weight: 3
					}
					{
						characteristicName: "skin-tone",
						selectedValue: "red",
						weight: 3
					}
					{
						characteristicName: "skin-tone",
						selectedValue: "yellow",
						weight: 3
					}
					{
						characteristicName: "skin-tone",
						selectedValue: "white",
						weight: 1
					}
					{
						characteristicName: "disabilities",
						selectedValue: "vision",
						weight: -3
					}
					{
						characteristicName: "disabilities",
						selectedValue: "hearing",
						weight: -3
					}
					{
						characteristicName: "disabilities",
						selectedValue: "locomotion",
						weight: 3
					}
				])
			}`

		response = gql(router, query, storedVariables["psychologist_3_token"])
		assert.Equal(t, "{\"data\":{\"setMyPsychologistPreferences\":null}}", response.Body.String())

	})

	t.Run("should create psychologist 4", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
				input: {
				email: "marcos.greco@psi.com.br"
				password: "Abc123!@#"
				role: PSYCHOLOGIST
				}
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: {
				email: "marcos.greco@psi.com.br"
				password: "Abc123!@#"
			}) {
				token
			}
		}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["psychologist_4_token"] = token

		query = `mutation {
			upsertMyPsychologistProfile(input: {
				fullName: "Marcos Greco"
				likeName: "Dr. Marcos",
				birthDate: 822708000,
				city: "Manaus - AM"
			})
		}`

		response = gql(router, query, storedVariables["psychologist_4_token"])
		assert.Equal(t, "{\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		query = `mutation {
			setMyPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "gender",
					selectedValues: ["male"]
				},
				{
					characteristicName: "lgbtqiaplus",
					selectedValues: ["true"]
				},
				{
					characteristicName: "skin-tone",
					selectedValues: ["white"]
				},
				{
					characteristicName: "sign-language",
					selectedValues: ["false"]
				},
				{
					characteristicName: "disabilities",
					selectedValues: []
				},
			])
		}`

		response = gql(router, query, storedVariables["psychologist_4_token"])
		assert.Equal(t, "{\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPsychologistPreferences(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValue: "true",
					weight: 3
				}
				{
					characteristicName: "has-consulted-before",
					selectedValue: "false",
					weight: 0
				}
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 0
				}
				{
					characteristicName: "gender",
					selectedValue: "male",
					weight: 0
				}
				{
					characteristicName: "gender",
					selectedValue: "non-binary",
					weight: 0
				}
				{
					characteristicName: "age",
					selectedValue: "child",
					weight: 3
				}
				{
					characteristicName: "age",
					selectedValue: "teen",
					weight: 1
				}
				{
					characteristicName: "age",
					selectedValue: "young-adult",
					weight: 1
				}
				{
					characteristicName: "age",
					selectedValue: "adult",
					weight: -1
				}
				{
					characteristicName: "age",
					selectedValue: "elderly",
					weight: -3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "true",
					weight: 3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "false",
					weight: 0
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "black",
					weight: 0
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "red",
					weight: 0
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "yellow",
					weight: 0
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "white",
					weight: 0
				}
				{
					characteristicName: "disabilities",
					selectedValue: "vision",
					weight: 0
				}
				{
					characteristicName: "disabilities",
					selectedValue: "hearing",
					weight: -3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "locomotion",
					weight: 0
				}
			])
		}`

		response = gql(router, query, storedVariables["psychologist_4_token"])
		assert.Equal(t, "{\"data\":{\"setMyPsychologistPreferences\":null}}", response.Body.String())

	})

	t.Run("should create psychologist 5", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
				input: {
				email: "tassia.lopes@psi.com.br"
				password: "Abc123!@#"
				role: PSYCHOLOGIST
				}
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: {
				email: "tassia.lopes@psi.com.br"
				password: "Abc123!@#"
			}) {
				token
			}
		}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["psychologist_5_token"] = token

		query = `mutation {
			upsertMyPsychologistProfile(input: {
				fullName: "Tássia Lopes"
				likeName: "Dra. Tássia",
				birthDate: 772513200,
				city: "Vitória - ES"
			})
		}`

		response = gql(router, query, storedVariables["psychologist_5_token"])
		assert.Equal(t, "{\"data\":{\"upsertMyPsychologistProfile\":null}}", response.Body.String())

		query = `mutation {
			setMyPsychologistCharacteristicChoices(input: [
				{
					characteristicName: "gender",
					selectedValues: ["female"]
				},
				{
					characteristicName: "lgbtqiaplus",
					selectedValues: ["false"]
				},
				{
					characteristicName: "skin-tone",
					selectedValues: ["black"]
				},
				{
					characteristicName: "sign-language",
					selectedValues: ["true"]
				},
				{
					characteristicName: "disabilities",
					selectedValues: []
				},
			])
		}`

		response = gql(router, query, storedVariables["psychologist_5_token"])
		assert.Equal(t, "{\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPsychologistPreferences(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValue: "true",
					weight: -3
				}
				{
					characteristicName: "has-consulted-before",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 3
				}
				{
					characteristicName: "gender",
					selectedValue: "male",
					weight: 0
				}
				{
					characteristicName: "gender",
					selectedValue: "non-binary",
					weight: 3
				}
				{
					characteristicName: "age",
					selectedValue: "child",
					weight: -3
				}
				{
					characteristicName: "age",
					selectedValue: "teen",
					weight: -3
				}
				{
					characteristicName: "age",
					selectedValue: "young-adult",
					weight: 3
				}
				{
					characteristicName: "age",
					selectedValue: "adult",
					weight: 1
				}
				{
					characteristicName: "age",
					selectedValue: "elderly",
					weight: 0
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "true",
					weight: 3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "black",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "red",
					weight: 1
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "yellow",
					weight: 1
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "white",
					weight: 1
				}
				{
					characteristicName: "disabilities",
					selectedValue: "vision",
					weight: 1
				}
				{
					characteristicName: "disabilities",
					selectedValue: "hearing",
					weight: 3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "locomotion",
					weight: 1
				}
			])
		}`

		response = gql(router, query, storedVariables["psychologist_5_token"])
		assert.Equal(t, "{\"data\":{\"setMyPsychologistPreferences\":null}}", response.Body.String())

	})

	t.Run("should create patient 1", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
				input: {
				email: "melissa.duarte@gmail.com"
				password: "Abc123!@#"
				role: PATIENT
				}
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: {
				email: "melissa.duarte@gmail.com"
				password: "Abc123!@#"
			}) {
				token
			}
		}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["patient_1_token"] = token

		query = `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Melissa Duarte"
				likeName: "Melissa",
				birthDate: 904618800,
				city: "Cuiabá - MT"
			})
		}`

		response = gql(router, query, storedVariables["patient_1_token"])
		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientCharacteristicChoices(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValues: ["true"]
				},
				{
					characteristicName: "gender",
					selectedValues: ["female"]
				},
				{
					characteristicName: "age",
					selectedValues: ["young-adult"]
				},
				{
					characteristicName: "lgbtqiaplus",
					selectedValues: ["false"]
				},
				{
					characteristicName: "skin-tone",
					selectedValues: ["yellow"]
				},
				{
					characteristicName: "disabilities",
					selectedValues: ["hearing"]
				},
			])
		}`

		response = gql(router, query, storedVariables["patient_1_token"])
		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientPreferences(input: [
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 3
				}
				{
					characteristicName: "gender",
					selectedValue: "male",
					weight: 3
				}
				{
					characteristicName: "gender",
					selectedValue: "non-binary",
					weight: 3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "true",
					weight: 3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "black",
					weight: 0
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "red",
					weight: 0
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "yellow",
					weight: 0
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "white",
					weight: 0
				}
				{
					characteristicName: "sign-language",
					selectedValue: "true",
					weight: 3
				}
				{
					characteristicName: "sign-language",
					selectedValue: "false",
					weight: -3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "vision",
					weight: 0
				}
				{
					characteristicName: "disabilities",
					selectedValue: "hearing",
					weight: 3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "locomotion",
					weight: 0
				}
			])
		}`

		response = gql(router, query, storedVariables["patient_1_token"])
		assert.Equal(t, "{\"data\":{\"setMyPatientPreferences\":null}}", response.Body.String())

	})

	t.Run("should create patient 2", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
				input: {
				email: "sandra.horta@gmail.com"
				password: "Abc123!@#"
				role: PATIENT
				}
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: {
				email: "sandra.horta@gmail.com"
				password: "Abc123!@#"
			}) {
				token
			}
		}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["patient_2_token"] = token

		query = `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Sandra Horta"
				likeName: "Sandra",
				birthDate: -237589200,
				city: "Rio Branco - AC"
			})
		}`

		response = gql(router, query, storedVariables["patient_2_token"])
		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientCharacteristicChoices(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValues: ["false"]
				},
				{
					characteristicName: "gender",
					selectedValues: ["female"]
				},
				{
					characteristicName: "age",
					selectedValues: ["elderly"]
				},
				{
					characteristicName: "lgbtqiaplus",
					selectedValues: ["false"]
				},
				{
					characteristicName: "skin-tone",
					selectedValues: ["red"]
				},
				{
					characteristicName: "disabilities",
					selectedValues: []
				},
			])
		}`

		response = gql(router, query, storedVariables["patient_2_token"])
		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientPreferences(input: [
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 0
				}
				{
					characteristicName: "gender",
					selectedValue: "male",
					weight: 1
				}
				{
					characteristicName: "gender",
					selectedValue: "non-binary",
					weight: -3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "true",
					weight: 0
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "black",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "red",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "yellow",
					weight: -3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "white",
					weight: -1
				}
				{
					characteristicName: "sign-language",
					selectedValue: "true",
					weight: -3
				}
				{
					characteristicName: "sign-language",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "vision",
					weight: 0
				}
				{
					characteristicName: "disabilities",
					selectedValue: "hearing",
					weight: -3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "locomotion",
					weight: 0
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
				email: "joao.melo@gmail.com"
				password: "Abc123!@#"
				role: PATIENT
				}
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: {
				email: "joao.melo@gmail.com"
				password: "Abc123!@#"
			}) {
				token
			}
		}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["patient_3_token"] = token

		query = `mutation {
			upsertMyPatientProfile(input: {
				fullName: "João Melo"
				likeName: "João",
				birthDate: 775191600,
				city: "Belém - PA"
			})
		}`

		response = gql(router, query, storedVariables["patient_3_token"])
		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientCharacteristicChoices(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValues: ["false"]
				},
				{
					characteristicName: "gender",
					selectedValues: ["non-binary"]
				},
				{
					characteristicName: "age",
					selectedValues: ["teen"]
				},
				{
					characteristicName: "lgbtqiaplus",
					selectedValues: ["true"]
				},
				{
					characteristicName: "skin-tone",
					selectedValues: ["white"]
				},
				{
					characteristicName: "disabilities",
					selectedValues: []
				},
			])
		}`

		response = gql(router, query, storedVariables["patient_3_token"])
		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientPreferences(input: [
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 3
				}
				{
					characteristicName: "gender",
					selectedValue: "male",
					weight: 0
				}
				{
					characteristicName: "gender",
					selectedValue: "non-binary",
					weight: 3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "true",
					weight: 3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "false",
					weight: 0
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "black",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "red",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "yellow",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "white",
					weight: 1
				}
				{
					characteristicName: "sign-language",
					selectedValue: "true",
					weight: -3
				}
				{
					characteristicName: "sign-language",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "vision",
					weight: 0
				}
				{
					characteristicName: "disabilities",
					selectedValue: "hearing",
					weight: -3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "locomotion",
					weight: -3
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
				email: "paula.santos@gmail.com"
				password: "Abc123!@#"
				role: PATIENT
				}
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: {
				email: "paula.santos@gmail.com"
				password: "Abc123!@#"
			}) {
				token
			}
		}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["patient_4_token"] = token

		query = `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Paula Santos"
				likeName: "Paula",
				birthDate: 800074800,
				city: "Blumenau - SC"
			})
		}`

		response = gql(router, query, storedVariables["patient_4_token"])
		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientCharacteristicChoices(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValues: ["false"]
				},
				{
					characteristicName: "gender",
					selectedValues: ["female"]
				},
				{
					characteristicName: "age",
					selectedValues: ["elderly"]
				},
				{
					characteristicName: "lgbtqiaplus",
					selectedValues: ["false"]
				},
				{
					characteristicName: "skin-tone",
					selectedValues: ["white"]
				},
				{
					characteristicName: "disabilities",
					selectedValues: ["vision"]
				},
			])
		}`

		response = gql(router, query, storedVariables["patient_4_token"])
		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientPreferences(input: [
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: -1
				}
				{
					characteristicName: "gender",
					selectedValue: "male",
					weight: 3
				}
				{
					characteristicName: "gender",
					selectedValue: "non-binary",
					weight: -3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "true",
					weight: -3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "black",
					weight: -3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "red",
					weight: -3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "yellow",
					weight: 0
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "white",
					weight: 3
				}
				{
					characteristicName: "sign-language",
					selectedValue: "true",
					weight: -3
				}
				{
					characteristicName: "sign-language",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "vision",
					weight: 3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "hearing",
					weight: -3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "locomotion",
					weight: -1
				}
			])
		}`

		response = gql(router, query, storedVariables["patient_4_token"])
		assert.Equal(t, "{\"data\":{\"setMyPatientPreferences\":null}}", response.Body.String())

	})

	t.Run("should create patient 5", func(t *testing.T) {

		query := `mutation {
			createUserWithPassword(
				input: {
				email: "jorge.martins@gmail.com"
				password: "Abc123!@#"
				role: PATIENT
				}
			)
		}`

		response := gql(router, query, storedVariables["coordinator_token"])
		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", response.Body.String())

		query = `{
			authenticateUser(input: {
				email: "jorge.martins@gmail.com"
				password: "Abc123!@#"
			}) {
				token
			}
		}`

		response = gql(router, query, "")

		token := fastjson.GetString(response.Body.Bytes(), "data", "authenticateUser", "token")
		assert.NotEqual(t, "", token)

		storedVariables["patient_5_token"] = token

		query = `mutation {
			upsertMyPatientProfile(input: {
				fullName: "Jorge Martins"
				likeName: "Jorge",
				birthDate: 826858800,
				city: "Fortaleza - CE"
			})
		}`

		response = gql(router, query, storedVariables["patient_5_token"])
		assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientCharacteristicChoices(input: [
				{
					characteristicName: "has-consulted-before",
					selectedValues: ["true"]
				},
				{
					characteristicName: "gender",
					selectedValues: ["male"]
				},
				{
					characteristicName: "age",
					selectedValues: ["young-adult"]
				},
				{
					characteristicName: "lgbtqiaplus",
					selectedValues: ["true"]
				},
				{
					characteristicName: "skin-tone",
					selectedValues: ["black"]
				},
				{
					characteristicName: "disabilities",
					selectedValues: []
				},
			])
		}`

		response = gql(router, query, storedVariables["patient_5_token"])
		assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", response.Body.String())

		query = `mutation {
			setMyPatientPreferences(input: [
				{
					characteristicName: "gender",
					selectedValue: "female",
					weight: 3
				}
				{
					characteristicName: "gender",
					selectedValue: "male",
					weight: 0
				}
				{
					characteristicName: "gender",
					selectedValue: "non-binary",
					weight: 3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "true",
					weight: 3
				}
				{
					characteristicName: "lgbtqiaplus",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "black",
					weight: 3
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "red",
					weight: 1
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "yellow",
					weight: 1
				}
				{
					characteristicName: "skin-tone",
					selectedValue: "white",
					weight: 1
				}
				{
					characteristicName: "sign-language",
					selectedValue: "true",
					weight: -3
				}
				{
					characteristicName: "sign-language",
					selectedValue: "false",
					weight: 3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "vision",
					weight: 1
				}
				{
					characteristicName: "disabilities",
					selectedValue: "hearing",
					weight: 3
				}
				{
					characteristicName: "disabilities",
					selectedValue: "locomotion",
					weight: 1
				}
			])
		}`

		response = gql(router, query, storedVariables["patient_5_token"])
		assert.Equal(t, "{\"data\":{\"setMyPatientPreferences\":null}}", response.Body.String())

	})

}
