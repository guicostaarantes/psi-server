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
					name: "issues",
					type: MULTIPLE,
					possibleValues: [
						"aging",
						"anxiety",
						"autism",
						"body",
						"career",
						"career-choice",
						"compulsion",
						"death-of-keen",
						"depression",
						"desease-in-keen",
						"desease-in-self",
						"discrimination",
						"drug-abuse-keen",
						"drug-abuse-self",
						"existential-crisis",
						"family-conflicts",
						"financial",
						"food-disorder",
						"learning-challenges",
						"memory-loss",
						"mood-control",
						"panic",
						"parenting",
						"partner-abuse",
						"partner-conflicts",
						"physical-pain",
						"post-traumatic-stress",
						"pre-pregnancy",
						"pregnancy",
						"schizophrenia",
						"self-development",
						"self-harm",
						"sexual",
						"sexuality",
						"sleep",
						"social-interactions",
						"stammering",
						"stress",
						"unknown",
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
					{
						key: "pat-char:issues"
						value: "Marque as opções que te incetivam a procurar tratamento psicológico:"
					}
					{ key: "pat-char:issues:aging", value: "Tenho dificuldades de aceitar o meu envelhecimento." }
					{ key: "pat-char:issues:anxiety", value: "Sofro de crises de ansiedade." }
					{ key: "pat-char:issues:autism", value: "Sou diagnosticado com Aspergers ou autismo." }
					{ key: "pat-char:issues:body", value: "Sofro pela forma como meu corpo é, não o aceito ou não consigo mudá-lo." }
					{ key: "pat-char:issues:career", value: "Sofro com assuntos relacionados à minha carreira profissional." }
					{ key: "pat-char:issues:career-choice", value: "Tenho dificuldades em escolher uma carreira, ou gostaria de fazer uma transição de carreira." }
					{ key: "pat-char:issues:compulsion", value: "Sofro de transtornos compulsivos." }
					{ key: "pat-char:issues:death-of-keen", value: "Sofro pelo falecimento de um ente querido." }
					{ key: "pat-char:issues:depression", value: "Sou diagnosticado ou acredito que tenha depressão." }
					{ key: "pat-char:issues:desease-in-keen", value: "Sofro por uma doença que afeta um ente querido." }
					{ key: "pat-char:issues:desease-in-self", value: "Sofro por uma doença que afeta a mim." }
					{ key: "pat-char:issues:discrimination", value: "Sofro por perceber discriminação de raça, gênero, ou similar, dirigida a mim." }
					{ key: "pat-char:issues:drug-abuse-keen", value: "Sofro por um ente querido que possui um vício em álcool, cigarro ou drogas psicoativas." }
					{ key: "pat-char:issues:drug-abuse-self", value: "Sou viciado em álcool, cigarro ou drogas psicoativas e quero me libertar do vício." }
					{ key: "pat-char:issues:existential-crisis", value: "Sofro de crises existenciais." }
					{ key: "pat-char:issues:family-conflicts", value: "Sofro por conflitos com meus familiares." }
					{ key: "pat-char:issues:financial", value: "Sofro de dificuldades financeiras e isso afeta o meu estado mental." }
					{ key: "pat-char:issues:food-disorder", value: "Sofro de transtornos alimentares (ex: anorexia, bulimia, consumo compulsivo)." }
					{ key: "pat-char:issues:learning-challenges", value: "Sofro de dificuldades de aprendizado." }
					{ key: "pat-char:issues:memory-loss", value: "Sofro de esquecimentos ou perda de memória." }
					{ key: "pat-char:issues:mood-control", value: "Sofro de mudanças frequentes de humor sem razão aparente." }
					{ key: "pat-char:issues:panic", value: "Sofro de ataques de pânico." }
					{ key: "pat-char:issues:parenting", value: "Tenho dificuldades de educar meus filhos." }
					{ key: "pat-char:issues:partner-abuse", value: "Estou em um relacionamento mentalmente ou fisicamente abusivo." }
					{ key: "pat-char:issues:partner-conflicts", value: "Estou passando por um momento difícil no meu casamento ou relacionamento." }
					{ key: "pat-char:issues:physical-pain", value: "Sofro de dores físicas crônicas que podem ser atenuadas com tratamento psicológico." }
					{ key: "pat-char:issues:post-traumatic-stress", value: "Sofro prolongadamente em razão de um evento terrível que vivenciei ou testemunhei e não consigo me recuperar." }
					{ key: "pat-char:issues:pre-pregnancy", value: "Sofro por conflitos relacionados a novos filhos (ex: dificuldade de engravidar, adoção, reprodução assistida)." }
					{ key: "pat-char:issues:pregnancy", value: "Sofro por problemas relativos à gestação ou após a gestação (ex: luto perinatal, depressão pós-parto)." }
					{ key: "pat-char:issues:schizophrenia", value: "Sou diagnosticado com esquizofrenia." }
					{ key: "pat-char:issues:self-development", value: "Procuro desenvolver habilidades pessoais e profissionais por meio do tratamento psicológico." }
					{ key: "pat-char:issues:self-harm", value: "Frequentemente tenho vontade de cometer suicídio ou me auto-mutilar." }
					{ key: "pat-char:issues:sexual", value: "Sofro de problemas de desempenho sexual (ex: impotência, perda de libido)." }
					{ key: "pat-char:issues:sexuality", value: "Tenho dificuldades de entender ou aceitar a minha sexualidade." }
					{ key: "pat-char:issues:sleep", value: "Sofro de insônia ou outros problemas psicológicos relacionados ao sono." }
					{ key: "pat-char:issues:stammering", value: "Sofro de gagueira ou outros problemas psicológicos relacionados à fala." }
					{ key: "pat-char:issues:social-interactions", value: "Tenho dificuldades de me relacionar socialmente (ex: timidez, insegurança)." }
					{ key: "pat-char:issues:stress", value: "Tenho estado estressado e nervoso de forma recorrente ou constante." }
					{ key: "pat-char:issues:unknown", value: "Busco um diagnóstico e entendo que o tratamento psicológico pode me auxiliar." }
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
					{ key: "pat-pref:issues:aging", value: "Quão interessado você está em tratar pacientes com dificuldades de aceitar o seu envelhecimento." }
					{ key: "pat-pref:issues:anxiety", value: "Quão interessado você está em tratar pacientes que sofrem de crises de ansiedade." }
					{ key: "pat-pref:issues:autism", value: "Quão interessado você está em tratar pacientes que são diagnosticados com Aspergers ou autismo." }
					{ key: "pat-pref:issues:body", value: "Quão interessado você está em tratar pacientes que sofrem pela forma como seu corpo é, não o aceita ou não consegue mudá-lo." }
					{ key: "pat-pref:issues:career", value: "Quão interessado você está em tratar pacientes que sofrem com assuntos relacionados à carreira profissional." }
					{ key: "pat-pref:issues:career-choice", value: "Quão interessado você está em tratar pacientes que têm dificuldades em escolher uma carreira, ou gostariam de fazer uma transição de carreira." }
					{ key: "pat-pref:issues:compulsion", value: "Quão interessado você está em tratar pacientes que sofrem de transtornos compulsivos." }
					{ key: "pat-pref:issues:death-of-keen", value: "Quão interessado você está em tratar pacientes que sofrem pelo falecimento de um ente querido." }
					{ key: "pat-pref:issues:depression", value: "Quão interessado você está em tratar pacientes que são diagnosticados ou acreditam ter depressão." }
					{ key: "pat-pref:issues:desease-in-keen", value: "Quão interessado você está em tratar pacientes que sofrem por uma doença que afeta um ente querido." }
					{ key: "pat-pref:issues:desease-in-self", value: "Quão interessado você está em tratar pacientes que sofrem por uma doença que afeta a si próprios." }
					{ key: "pat-pref:issues:discrimination", value: "Quão interessado você está em tratar pacientes que sofrem por perceber discriminação de raça, gênero, ou similar, dirigida a si próprios." }
					{ key: "pat-pref:issues:drug-abuse-keen", value: "Quão interessado você está em tratar pacientes que sofrem por um ente querido que possui um vício em álcool, cigarro ou drogas psicoativas." }
					{ key: "pat-pref:issues:drug-abuse-self", value: "Quão interessado você está em tratar pacientes que são viciados em álcool, cigarro ou drogas psicoativas e querem se libertar do vício." }
					{ key: "pat-pref:issues:existential-crisis", value: "Quão interessado você está em tratar pacientes que sofrem de crises existenciais." }
					{ key: "pat-pref:issues:family-conflicts", value: "Quão interessado você está em tratar pacientes que sofrem por conflitos com seus familiares." }
					{ key: "pat-pref:issues:financial", value: "Quão interessado você está em tratar pacientes que sofrem de dificuldades financeiras e isso afeta o seu estado mental." }
					{ key: "pat-pref:issues:food-disorder", value: "Quão interessado você está em tratar pacientes que sofrem de transtornos alimentares (ex: anorexia, bulimia, consumo compulsivo)." }
					{ key: "pat-pref:issues:learning-challenges", value: "Quão interessado você está em tratar pacientes que sofrem de dificuldades de aprendizado." }
					{ key: "pat-pref:issues:memory-loss", value: "Quão interessado você está em tratar pacientes que sofrem de esquecimentos ou perda de memória." }
					{ key: "pat-pref:issues:mood-control", value: "Quão interessado você está em tratar pacientes que sofrem de mudanças frequentes de humor sem razão aparente." }
					{ key: "pat-pref:issues:panic", value: "Quão interessado você está em tratar pacientes que sofrem de ataques de pânico." }
					{ key: "pat-pref:issues:parenting", value: "Quão interessado você está em tratar pacientes que têm dificuldades de educar seus filhos." }
					{ key: "pat-pref:issues:partner-abuse", value: "Quão interessado você está em tratar pacientes que estão em um relacionamento mentalmente ou fisicamente abusivo." }
					{ key: "pat-pref:issues:partner-conflicts", value: "Quão interessado você está em tratar pacientes que estão passando por um momento difícil no seu casamento ou relacionamento." }
					{ key: "pat-pref:issues:physical-pain", value: "Quão interessado você está em tratar pacientes que sofrem de dores físicas crônicas que podem ser atenuadas com tratamento psicológico." }
					{ key: "pat-pref:issues:post-traumatic-stress", value: "Quão interessado você está em tratar pacientes que sofrem prolongadamente em razão de um evento terrível que vivenciaram ou testemunharam e não conseguem se recuperar." }
					{ key: "pat-pref:issues:pre-pregnancy", value: "Quão interessado você está em tratar pacientes que sofrem por conflitos relacionados a novos filhos (ex: dificuldade de engravidar, adoção, reprodução assistida)." }
					{ key: "pat-pref:issues:pregnancy", value: "Quão interessado você está em tratar pacientes que sofrem por problemas relativos à gestação ou após a gestação (ex: luto perinatal, depressão pós-parto)." }
					{ key: "pat-pref:issues:schizophrenia", value: "Quão interessado você está em tratar pacientes que são diagnosticados com esquizofrenia." }
					{ key: "pat-pref:issues:self-development", value: "Quão interessado você está em tratar pacientes que procuram desenvolver habilidades pessoais e profissionais por meio do tratamento psicológico." }
					{ key: "pat-pref:issues:self-harm", value: "Quão interessado você está em tratar pacientes que frequentemente têm vontade de cometer suicídio ou se auto-mutilar." }
					{ key: "pat-pref:issues:sexual", value: "Quão interessado você está em tratar pacientes que sofrem de problemas de desempenho sexual (ex: impotência, perda de libido)." }
					{ key: "pat-pref:issues:sexuality", value: "Quão interessado você está em tratar pacientes que têm dificuldades de entender ou aceitar a sua sexualidade." }
					{ key: "pat-pref:issues:sleep", value: "Quão interessado você está em tratar pacientes que sofrem de insônia ou outros problemas psicológicos relacionados ao sono." }
					{ key: "pat-pref:issues:stammering", value: "Quão interessado você está em tratar pacientes que sofrem de gagueira ou outros problemas psicológicos relacionados à fala." }
					{ key: "pat-pref:issues:social-interactions", value: "Quão interessado você está em tratar pacientes que têm dificuldades de se relacionar socialmente (ex: timidez, insegurança)." }
					{ key: "pat-pref:issues:stress", value: "Quão interessado você está em tratar pacientes que têm estado estressados e nervosos de forma recorrente ou constante." }
					{ key: "pat-pref:issues:unknown", value: "Quão interessado você está em tratar pacientes que buscam um diagnóstico e entendem que o tratamento psicológico podem os auxiliar." }
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
				email := fmt.Sprintf("psi%04d@exemplo.com", i)
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
						crp: "01/23%04d",
						whatsapp: "(11) 98765-%04d",
						instagram: "@psi.%s.%s",
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
						{
							characteristicName: "issues",
							selectedValue: "aging",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "anxiety",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "autism",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "body",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "career",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "career-choice",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "compulsion",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "death-of-keen",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "depression",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "desease-in-keen",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "desease-in-self",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "discrimination",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "drug-abuse-keen",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "drug-abuse-self",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "existential-crisis",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "family-conflicts",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "financial",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "food-disorder",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "learning-challenges",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "memory-loss",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "mood-control",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "panic",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "parenting",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "partner-abuse",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "partner-conflicts",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "physical-pain",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "post-traumatic-stress",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "pre-pregnancy",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "pregnancy",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "schizophrenia",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "self-development",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "self-harm",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "sexual",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "sexuality",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "sleep",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "social-interactions",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "stammering",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "stress",
							weight: %d
						}
						{
							characteristicName: "issues",
							selectedValue: "unknown",
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
				email := fmt.Sprintf("pac%04d@exemplo.com", i)
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
							characteristicName: "issues",
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
					[]string{"aging", "anxiety", "autism", "body", "career", "career-choice", "compulsion", "death-of-keen", "depression", "desease-in-keen", "desease-in-self", "discrimination", "drug-abuse-keen", "drug-abuse-self", "existential-crisis", "family-conflicts", "financial", "food-disorder", "learning-challenges", "memory-loss", "mood-control", "panic", "parenting", "partner-abuse", "partner-conflicts", "physical-pain", "post-traumatic-stress", "pre-pregnancy", "pregnancy", "schizophrenia", "self-development", "self-harm", "sexual", "sexuality", "sleep", "social-interactions", "stammering", "stress", "unknown"}[rand.Intn(39)],
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
