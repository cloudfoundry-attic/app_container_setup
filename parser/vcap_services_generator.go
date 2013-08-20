package parser

import (
	"encoding/json"
	"errors"
)

type ServiceJSON struct {
	Name        string                 `json:"name,omitempty"`
	Label       string                 `json:"label"`
	Tags        []string               `json:"tags,omitempty"`
	Credentials map[string]interface{} `json:"credentials,omitempty"`
	Plan        string                 `json:"plan,omitempty"`
	PlanOption  map[string]interface{} `json:"plan_option,omitempty"`
}

type InputServiceJSON struct {
	Credentials map[string]interface{} `json:"credentials"`
	Tags        []string               `json:"tags"`
	PlanOption  map[string]interface{} `json:"plan_option"`
	Label       string                 `json:"label"`
	Provider    string                 `json:"provider"`
	Version     string                 `json:"version"`
	Vendor      string                 `json:"vendor"`
	Plan        string                 `json:"plan"`
	Name        string                 `json:"name"`
}

var ErrMissingLabel = errors.New("Label cannot be empty")

func (parser *Parser) generateServicesJSON(services []InputServiceJSON) ([]byte, error) {
	servicesData := make(map[string]*ServiceJSON)

	for _, service := range services {
		if service.Label == "" {
			return nil, ErrMissingLabel
		}
		servicesData[service.Label] = &ServiceJSON{
			Name:        service.Name,
			Label:       service.Label,
			Tags:        service.Tags,
			Credentials: service.Credentials,
			Plan:        service.Plan,
			PlanOption:  service.PlanOption,
		}
	}

	return json.Marshal(servicesData)
}
