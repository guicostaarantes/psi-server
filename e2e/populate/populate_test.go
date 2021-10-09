package populate

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fastjson"
)

// WARNING: Running this test will make changes to your database state
// Proceed with caution and avoid running this in production environments

// Edit the parameters below. Since this a public repository, ALWAYS change the passwords
// before running this in a production environment
var API_URL = "http://localhost:7070/gql"
var COORDINATOR_USERNAME = "coordinator@psi.com.br"
var COORDINATOR_PASSWORD = "Abc123!@#"
var NEW_JOBRUNNER_USERNAME = "jobrunner@psi.com.br"
var NEW_JOBRUNNER_PASSWORD = "Abc123!@#"
var PSYCOLOGIST_QUANTITY = 50
var PATIENT_QUANTITY = 50
var NEW_PSYCHOLOGISTS_PASSWORD = "Abc123!@#"
var NEW_PATIENTS_PASSWORD = "Abc123!@#"

func gql(client *http.Client, query string, token string) *http.Response {

	body := fmt.Sprintf(`{"query": %q}`, query)

	request, err := http.NewRequest(http.MethodPost, API_URL, strings.NewReader(body))
	if err != nil {
		log.Fatalln("Error creating request")
	}

	request.Header["Authorization"] = []string{token}
	request.Header["Content-Type"] = []string{"application/json"}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalln("Error in response")
	}

	return response

}

func TestEnd2End(t *testing.T) {

	rand.Seed(time.Now().UnixMicro())

	storedVariables := map[string]string{}

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	t.Run("should log in as bootstrap coordinator", func(t *testing.T) {

		query := fmt.Sprintf(`{
			authenticateUser(input: {
				email: %q,
				password: %q
			}) {
				token
			}
		}`, COORDINATOR_USERNAME, COORDINATOR_PASSWORD)

		response := gql(client, query, "")
		body, _ := ioutil.ReadAll(response.Body)

		token := fastjson.GetString(body, "data", "authenticateUser", "token")
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

		response := gql(client, query, storedVariables["coordinator_token"])
		body, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, "{\"data\":{\"setPatientCharacteristics\":null}}", string(body))

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
				},
				{
					name: "methods",
					type: MULTIPLE,
					possibleValues: [
						"psychoanalysis",
						"behaviorism",
						"cognitive",
						"humanistic",
					]
				}
			])
		}`

		response = gql(client, query, storedVariables["coordinator_token"])
		body, _ = ioutil.ReadAll(response.Body)

		assert.Equal(t, "{\"data\":{\"setPsychologistCharacteristics\":null}}", string(body))

		query = `mutation {
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
					maximumPrice: 150,
					eligibleFor: "D,C,B,A"
				}
			])
		}`

		response = gql(client, query, storedVariables["coordinator_token"])
		body, _ = ioutil.ReadAll(response.Body)

		assert.Equal(t, "{\"data\":{\"setTreatmentPriceRanges\":null}}", string(body))

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
					{ key: "pat-char:income", value: "Qual é a renda mensal da sua família por cada membro?" }
					{ key: "pat-char:income:D", value: "De 0 a 1000 reais" }
					{ key: "pat-char:income:C", value: "De 1000 a 2000 reais" }
					{ key: "pat-char:income:B", value: "De 2000 a 5000 reais" }
					{ key: "pat-char:income:A", value: "Acima de 5000 reais" }
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
						value: "Você é capaz de atender pacientes utilizando LIBRAS (Linguagem Brasileira de Sinais)?"
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
						key: "psy-char:methods"
						value: "Você utiliza quais dessas abordagens da psicologia?"
					}
					{ key: "psy-char:methods:psychoanalysis", value: "Psicoanálise" }
					{ key: "psy-char:methods:behaviorism", value: "Behaviorismo" }
					{ key: "psy-char:methods:cognitive", value: "Terapia cognitivo-comportamental" }
					{ key: "psy-char:methods:humanistic", value: "Psicologia humanista" }
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
					{ key: "psy-pref:methods:psychoanalysis", value: "Quão confortável você se sente sendo atendido por um psicólogo pela abordagem de psicoanálise?" }
					{ key: "psy-pref:methods:behaviorism", value: "Quão confortável você se sente sendo atendido por um psicólogo pela abordagem de behaviorismo?" }
					{ key: "psy-pref:methods:cognitive", value: "Quão confortável você se sente sendo atendido por um psicólogo pela abordagem de terapia cognitivo-comportamental?" }
					{ key: "psy-pref:methods:humanistic", value: "Quão confortável você se sente sendo atendido por um psicólogo pela abordagem de psicologia humanista?" }
				]
			)
		}`

		response = gql(client, query, storedVariables["coordinator_token"])
		body, _ = ioutil.ReadAll(response.Body)

		assert.Equal(t, "{\"data\":{\"setTranslations\":null}}", string(body))

	})

	t.Run("should create jobrunner", func(t *testing.T) {

		query := fmt.Sprintf(`mutation {
			createUserWithPassword(
			  input: {
				email: %q
				password: %q
				role: JOBRUNNER
			  }
			)
		}`, NEW_JOBRUNNER_USERNAME, NEW_JOBRUNNER_PASSWORD)

		response := gql(client, query, storedVariables["coordinator_token"])
		body, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", string(body))

	})

	t.Run("should create psychologists", func(t *testing.T) {

		var wg sync.WaitGroup

		for i := 1; i <= PSYCOLOGIST_QUANTITY; i++ {
			wg.Add(1)
			go func(i int) {
				email := fmt.Sprintf("psi%03d@exemplo.com", i)
				gender := []string{"female", "male", "non-binary"}[rand.Intn(3)]
				lastName := LastNames[rand.Intn(40)]
				firstName := ""
				switch gender {
				case "female":
					firstName = FemaleNames[rand.Intn(40)]
				case "male":
					firstName = MaleNames[rand.Intn(40)]
				case "non-binary":
					firstName = []string{FemaleNames[rand.Intn(40)], MaleNames[rand.Intn(40)]}[rand.Intn(2)]
				}

				query := fmt.Sprintf(`mutation {
					createUserWithPassword(
						input: {
						email: %q
						password: %q
						role: PSYCHOLOGIST
						}
					)
				}`, email, NEW_PSYCHOLOGISTS_PASSWORD)

				response := gql(client, query, storedVariables["coordinator_token"])
				body, _ := ioutil.ReadAll(response.Body)

				assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", string(body))

				query = fmt.Sprintf(`{
					authenticateUser(input: {
						email: %q
						password: %q
					}) {
						token
					}
				}`, email, NEW_PSYCHOLOGISTS_PASSWORD)

				response = gql(client, query, "")
				body, _ = ioutil.ReadAll(response.Body)

				token := fastjson.GetString(body, "data", "authenticateUser", "token")
				assert.NotEqual(t, "", token)

				query = fmt.Sprintf(`mutation {
					upsertMyPsychologistProfile(input: {
						fullName: "%s %s"
						likeName: %q,
						birthDate: %d,
						city: "Belo Horizonte - MG",
						crp: "06/123%03d",
						whatsapp: "(31) 98765-4%03d",
						instagram: "@psi%s%s",
						bio: "Oi, meu nome é %s %s."
					})
				}`, firstName, lastName, firstName, 86400*rand.Intn(10000), i, i, strings.ToLower(firstName), strings.ToLower(lastName), firstName, lastName)

				response = gql(client, query, token)
				body, _ = ioutil.ReadAll(response.Body)

				assert.Equal(t, "{\"data\":{\"upsertMyPsychologistProfile\":null}}", string(body))

				query = fmt.Sprintf(`mutation {
					setMyPsychologistCharacteristicChoices(input: [
						{
							characteristicName: "gender",
							selectedValues: [%q]
						},
						{
							characteristicName: "lgbtqiaplus",
							selectedValues: [%q]
						},
						{
							characteristicName: "skin-tone",
							selectedValues: [%q]
						},
						{
							characteristicName: "sign-language",
							selectedValues: [%q]
						},
						{
							characteristicName: "disabilities",
							selectedValues: [%q]
						},
						{
							characteristicName: "methods",
							selectedValues: [%q]
						},
					])
				}`,
					gender,
					[]string{"true", "false"}[rand.Intn(2)],
					[]string{"black", "red", "yellow", "white"}[rand.Intn(4)],
					[]string{"true", "false"}[rand.Intn(2)],
					[]string{"vision", "hearing", "locomotion"}[rand.Intn(3)],
					[]string{"psychoanalysis", "behaviorism", "cognitive", "humanistic"}[rand.Intn(4)],
				)

				response = gql(client, query, token)
				body, _ = ioutil.ReadAll(response.Body)

				assert.Equal(t, "{\"data\":{\"setMyPsychologistCharacteristicChoices\":null}}", string(body))

				query = fmt.Sprintf(`mutation {
					setMyPsychologistPreferences(input: [
						{
							characteristicName: "has-consulted-before",
							selectedValue: "true",
							weight: %d
						}
						{
							characteristicName: "has-consulted-before",
							selectedValue: "false",
							weight: %d
						}
						{
							characteristicName: "gender",
							selectedValue: "female",
							weight: %d
						}
						{
							characteristicName: "gender",
							selectedValue: "male",
							weight: %d
						}
						{
							characteristicName: "gender",
							selectedValue: "non-binary",
							weight: %d
						}
						{
							characteristicName: "age",
							selectedValue: "child",
							weight: %d
						}
						{
							characteristicName: "age",
							selectedValue: "teen",
							weight: %d
						}
						{
							characteristicName: "age",
							selectedValue: "young-adult",
							weight: %d
						}
						{
							characteristicName: "age",
							selectedValue: "adult",
							weight: %d
						}
						{
							characteristicName: "age",
							selectedValue: "elderly",
							weight: %d
						}
						{
							characteristicName: "lgbtqiaplus",
							selectedValue: "true",
							weight: %d
						}
						{
							characteristicName: "lgbtqiaplus",
							selectedValue: "false",
							weight: %d
						}
						{
							characteristicName: "skin-tone",
							selectedValue: "black",
							weight: %d
						}
						{
							characteristicName: "skin-tone",
							selectedValue: "red",
							weight: %d
						}
						{
							characteristicName: "skin-tone",
							selectedValue: "yellow",
							weight: %d
						}
						{
							characteristicName: "skin-tone",
							selectedValue: "white",
							weight: %d
						}
						{
							characteristicName: "disabilities",
							selectedValue: "vision",
							weight: %d
						}
						{
							characteristicName: "disabilities",
							selectedValue: "hearing",
							weight: %d
						}
						{
							characteristicName: "disabilities",
							selectedValue: "locomotion",
							weight: %d
						}
					])
				}`,
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
				)

				response = gql(client, query, token)
				body, _ = ioutil.ReadAll(response.Body)

				assert.Equal(t, "{\"data\":{\"setMyPsychologistPreferences\":null}}", string(body))

				wg.Done()
			}(i)
		}
		wg.Wait()

	})

	t.Run("should create patients", func(t *testing.T) {

		var wg sync.WaitGroup

		for i := 1; i <= PATIENT_QUANTITY; i++ {
			wg.Add(1)
			go func(i int) {
				email := fmt.Sprintf("pac%03d@exemplo.com", i)
				gender := []string{"female", "male", "non-binary"}[rand.Intn(3)]
				birthDate := 86400 * rand.Intn(10000)
				lastName := LastNames[rand.Intn(40)]
				firstName := ""
				switch gender {
				case "female":
					firstName = FemaleNames[rand.Intn(40)]
				case "male":
					firstName = MaleNames[rand.Intn(40)]
				case "non-binary":
					firstName = []string{FemaleNames[rand.Intn(40)], MaleNames[rand.Intn(40)]}[rand.Intn(2)]
				}

				query := fmt.Sprintf(`mutation {
					createUserWithPassword(
						input: {
						email: %q
						password: %q
						role: PATIENT
						}
					)
				}`, email, NEW_PATIENTS_PASSWORD)

				response := gql(client, query, storedVariables["coordinator_token"])
				body, _ := ioutil.ReadAll(response.Body)

				assert.Equal(t, "{\"data\":{\"createUserWithPassword\":null}}", string(body))

				query = fmt.Sprintf(`{
					authenticateUser(input: {
						email: %q
						password: %q
					}) {
						token
					}
				}`, email, NEW_PATIENTS_PASSWORD)

				response = gql(client, query, "")
				body, _ = ioutil.ReadAll(response.Body)

				token := fastjson.GetString(body, "data", "authenticateUser", "token")
				assert.NotEqual(t, "", token)

				query = fmt.Sprintf(`mutation {
					upsertMyPatientProfile(input: {
						fullName: "%s %s"
						likeName: %q,
						birthDate: %d,
						city: "Belo Horizonte - MG"
					})
				}`, firstName, lastName, firstName, birthDate)

				response = gql(client, query, token)
				body, _ = ioutil.ReadAll(response.Body)

				assert.Equal(t, "{\"data\":{\"upsertMyPatientProfile\":null}}", string(body))

				query = fmt.Sprintf(`mutation {
					setMyPatientCharacteristicChoices(input: [
						{
							characteristicName: "has-consulted-before",
							selectedValues: [%q]
						},
						{
							characteristicName: "gender",
							selectedValues: [%q]
						},
						{
							characteristicName: "age",
							selectedValues: [%q]
						},
						{
							characteristicName: "lgbtqiaplus",
							selectedValues: [%q]
						},
						{
							characteristicName: "skin-tone",
							selectedValues: [%q]
						},
						{
							characteristicName: "disabilities",
							selectedValues: [%q]
						},
						{
							characteristicName: "income",
							selectedValues: [%q]
						}
					])
				}`,
					[]string{"true", "false"}[rand.Intn(2)],
					gender,
					[]string{"child", "teen", "young-adult", "adult", "elderly"}[rand.Intn(5)],
					[]string{"true", "false"}[rand.Intn(2)],
					[]string{"black", "red", "yellow", "white"}[rand.Intn(4)],
					[]string{"vision", "hearing", "locomotion"}[rand.Intn(3)],
					[]string{"D", "C", "B", "A"}[rand.Intn(4)],
				)

				response = gql(client, query, token)
				body, _ = ioutil.ReadAll(response.Body)

				assert.Equal(t, "{\"data\":{\"setMyPatientCharacteristicChoices\":null}}", string(body))

				query = fmt.Sprintf(`mutation {
					setMyPatientPreferences(input: [
						{
							characteristicName: "gender",
							selectedValue: "female",
							weight: %d
						}
						{
							characteristicName: "gender",
							selectedValue: "male",
							weight: %d
						}
						{
							characteristicName: "gender",
							selectedValue: "non-binary",
							weight: %d
						}
						{
							characteristicName: "lgbtqiaplus",
							selectedValue: "true",
							weight: %d
						}
						{
							characteristicName: "lgbtqiaplus",
							selectedValue: "false",
							weight: %d
						}
						{
							characteristicName: "skin-tone",
							selectedValue: "black",
							weight: %d
						}
						{
							characteristicName: "skin-tone",
							selectedValue: "red",
							weight: %d
						}
						{
							characteristicName: "skin-tone",
							selectedValue: "yellow",
							weight: %d
						}
						{
							characteristicName: "skin-tone",
							selectedValue: "white",
							weight: %d
						}
						{
							characteristicName: "sign-language",
							selectedValue: "true",
							weight: %d
						}
						{
							characteristicName: "sign-language",
							selectedValue: "false",
							weight: %d
						}
						{
							characteristicName: "disabilities",
							selectedValue: "vision",
							weight: %d
						}
						{
							characteristicName: "disabilities",
							selectedValue: "hearing",
							weight: %d
						}
						{
							characteristicName: "disabilities",
							selectedValue: "locomotion",
							weight: %d
						}
						{
							characteristicName: "methods",
							selectedValue: "psychoanalysis",
							weight: %d
						}
						{
							characteristicName: "methods",
							selectedValue: "behaviorism",
							weight: %d
						}
						{
							characteristicName: "methods",
							selectedValue: "cognitive",
							weight: %d
						}
						{
							characteristicName: "methods",
							selectedValue: "humanistic",
							weight: %d
						}
					])
				}`,
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
					[]int{-1, 0, 1, 3}[rand.Intn(4)],
				)

				response = gql(client, query, token)
				body, _ = ioutil.ReadAll(response.Body)

				assert.Equal(t, "{\"data\":{\"setMyPatientPreferences\":null}}", string(body))

				wg.Done()
			}(i)
		}
		wg.Wait()

	})

}
