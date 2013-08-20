package parser

import (
	"encoding/json"
	"time"
)

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
