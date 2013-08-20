package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Parser struct {
	systemEnvironmentVariables []EnvironmentPair
	userEnvironmentVariables   []EnvironmentPair
	profileDScript             string
}

type InputJSON struct {
	NatsData                     InputNatsDataJSON `json:"nats_data"`
	InstanceContainerPort        int               `json:"instance_container_port"`
	InstanceConsoleContainerPort int               `json:"instance_console_container_port"`
	InstanceDebugContainerPort   int               `json:"instance_debug_container_port"`
	InstanceGuid                 string            `json:"instance_guid"`
	StartedAtTimestamp           int64             `json:"started_at_timestamp"`
}

type InputNatsDataJSON struct {
	Limits             InputNatsLimitsJSON `json:"limits"`
	Debug              string              `json:"debug"`
	Index              int                 `json:"index"`
	ApplicationVersion string              `json:"version"`
	Name               string              `json:"name"`
	Uris               []string            `json:"uris"`
	Services           []InputServiceJSON  `json:"services"`
	Env                []string            `json:"env"`
}

type InputNatsLimitsJSON struct {
	Mem  int `json:"mem"`
	Disk int `json:"disk"`
	Fds  int `json:"fds"`
}

type EnvironmentPair struct {
	Name  string
	Value string
}

func NewParser() *Parser {
	parser := &Parser{}
	return parser
}

func (parser *Parser) GenerateEnvironmentScriptFromJSON(rawJSON string) (string, error) {
	var input InputJSON
	err := json.Unmarshal([]byte(rawJSON), &input)
	if err != nil {
		return "", err
	}

	parser.addSystemEnvironmentVariable("MEMORY_LIMIT", fmt.Sprintf("%dm", input.NatsData.Limits.Mem))
	parser.addSystemEnvironmentVariable("HOME", "$PWD/app")
	parser.addSystemEnvironmentVariable("TMPDIR", "$PWD/tmp")
	parser.addSystemEnvironmentVariable("VCAP_APP_HOST", "0.0.0.0")
	parser.addSystemEnvironmentVariable("VCAP_APP_PORT", strconv.Itoa(input.InstanceContainerPort))
	parser.addSystemEnvironmentVariable("VCAP_CONSOLE_IP", "0.0.0.0")
	parser.addSystemEnvironmentVariable("VCAP_CONSOLE_PORT", strconv.Itoa(input.InstanceConsoleContainerPort))
	parser.addSystemEnvironmentVariable("PORT", strconv.Itoa(input.InstanceContainerPort))

	if input.NatsData.Debug != "" {
		parser.addSystemEnvironmentVariable("VCAP_DEBUG_IP", "0.0.0.0")
		parser.addSystemEnvironmentVariable("VCAP_DEBUG_PORT", strconv.Itoa(input.InstanceDebugContainerPort))
		parser.addSystemEnvironmentVariable("VCAP_DEBUG_MODE", input.NatsData.Debug)
	}

	applicationJSON, err := parser.generateApplicationJSON(input)
	if err != nil {
		return "", err
	}
	parser.addSystemEnvironmentVariable("VCAP_APPLICATION", string(applicationJSON))

	servicesJSON, err := parser.generateServicesJSON(input.NatsData.Services)
	if err != nil {
		return "", err
	}
	parser.addSystemEnvironmentVariable("VCAP_SERVICES", string(servicesJSON))

	dbServicesRepresentations := parser.generateDBServiceRepresentationArray(input)
	if len(dbServicesRepresentations) > 0 {
		databaseUriGenerator := NewDatabaseURIGenerator(dbServicesRepresentations)
		databaseUrl, err := databaseUriGenerator.Generate()
		if err != nil {
			return "", err
		}
		parser.addSystemEnvironmentVariable("DATABASE_URL", databaseUrl)
	}

	parser.profileDScript = generateProfileDReader()

	userEnvironmentVariables := generateUserEnvironment(input.NatsData.Env)
	for _, envPair := range userEnvironmentVariables {
		parser.userEnvironmentVariables = append(parser.userEnvironmentVariables, envPair)
	}

	return parser.generateOutput(), nil
}

func (parser *Parser) generateDBServiceRepresentationArray(input InputJSON) []DBServiceRepresentation {
	services := input.NatsData.Services
	servicesData := make([]DBServiceRepresentation, len(services))

	for _, service := range services {
		uri, ok := service.Credentials["uri"].(string)
		if ok {
			servicesData = append(servicesData, DBServiceRepresentation{Name: service.Name, URI: uri})
		}
	}
	return servicesData
}

func (parser *Parser) addSystemEnvironmentVariable(name string, value string) {
	parser.systemEnvironmentVariables = append(parser.systemEnvironmentVariables, EnvironmentPair{Name: name, Value: value})
}

func (parser *Parser) generateOutput() string {
	output := ""
	for _, pair := range parser.systemEnvironmentVariables {
		output = fmt.Sprintf("%sexport %s=%s\n", output, pair.Name, strconv.Quote(pair.Value))
	}
	output += parser.profileDScript
	for _, pair := range parser.userEnvironmentVariables {
		output = fmt.Sprintf("%sexport %s=%s\n", output, pair.Name, strconv.Quote(pair.Value))
	}
	return output
}
