package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Parser struct {
	systemEnvironmentVariables []EnvironmentPair
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
	parser.systemEnvironmentVariables = []EnvironmentPair{}
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

	return parser.generateOutput(), nil
}

// VCAP_APPLICATION_JSON

type ApplicationJSON struct {
	InstanceId         string              `json:"instance_id"`
	InstanceIndex      int                 `json:"instance_index"`
	Host               string              `json:"host"`
	Port               int                 `json:"port"`
	StartedAtTimestamp int                 `json:"started_at_timestamp"`
	StartedAt          string              `json:"started_at"`
	Start              string              `json:"start"`
	StateTimestamp     int                 `json:"state_timestamp"`
	Limits             InputNatsLimitsJSON `json:"limits"`
	ApplicationVersion string              `json:"application_version"`
	Version            string              `json:"version"`
	ApplicationName    string              `json:"application_name"`
	Name               string              `json:"name"`
	Uris               []string            `json:"uris"`
	ApplicationUris    []string            `json:"application_uris"`
	Users              interface{}         `json:"users"`
}

func (parser *Parser) generateApplicationJSON(input InputJSON) ([]byte, error) {
	applicationData := new(ApplicationJSON)
	applicationData.InstanceId = input.InstanceGuid
	applicationData.InstanceIndex = input.NatsData.Index
	applicationData.Host = "0.0.0.0"
	applicationData.Port = input.InstanceContainerPort

	applicationData.StartedAtTimestamp = int(input.StartedAtTimestamp)
	applicationData.StateTimestamp = int(input.StartedAtTimestamp)
	startTime := time.Unix(input.StartedAtTimestamp, 0).UTC().Format("2006-01-02 15:04:05 -0700")
	applicationData.Start = startTime
	applicationData.StartedAt = startTime

	applicationData.Limits = input.NatsData.Limits

	applicationData.ApplicationVersion = input.NatsData.ApplicationVersion
	applicationData.Version = input.NatsData.ApplicationVersion

	applicationData.ApplicationName = input.NatsData.Name
	applicationData.Name = input.NatsData.Name

	applicationData.ApplicationUris = input.NatsData.Uris
	applicationData.Uris = input.NatsData.Uris

	applicationData.Users = nil

	return json.Marshal(applicationData)
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
	return output
}
