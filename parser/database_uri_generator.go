package parser

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

type DBServiceRepresentation struct {
	Name string
	URI  string
}

type DatabaseURIGenerator struct {
	services              []DBServiceRepresentation
	databaseSchemeMapping map[string]string
}

func NewDatabaseURIGenerator(services []DBServiceRepresentation) *DatabaseURIGenerator {
	return &DatabaseURIGenerator{
		services:              services,
		databaseSchemeMapping: map[string]string{"mysql": "mysql2", "mysql2": "mysql2", "postgres": "postgres", "postgresql": "postgres"},
	}
}

func (generator *DatabaseURIGenerator) Generate() (string, error) {
	var err error
	generator.services, err = generator.filterRelationalDatabasesAndFixScheme()
	if err != nil {
		return "", err
	}

	productionService, err := generator.findProductionDatabaseService()
	if err != nil {
		return "", err
	}

	return productionService.URI, nil
}

func (generator *DatabaseURIGenerator) filterRelationalDatabasesAndFixScheme() ([]DBServiceRepresentation, error) {
	filteredServices := make([]DBServiceRepresentation, 0)
	for _, service := range generator.services {
		parsedURI, err := url.Parse(service.URI)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid database URI \"%s\"", sanitizeURI(service.URI)))
		}
		mappedScheme, ok := generator.databaseSchemeMapping[parsedURI.Scheme]
		if ok {
			parsedURI.Scheme = mappedScheme
			filteredServices = append(filteredServices, DBServiceRepresentation{
				Name: service.Name,
				URI:  parsedURI.String(),
			})
		}
	}
	return filteredServices, nil
}

func (generator *DatabaseURIGenerator) findProductionDatabaseService() (DBServiceRepresentation, error) {
	switch len(generator.services) {
	case 0:
		return DBServiceRepresentation{}, nil
	case 1:
		return generator.services[0], nil
	default:
		re, _ := regexp.Compile(`^.*production$|^.*prod$`)
		for _, service := range generator.services {
			if re.Match([]byte(service.Name)) {
				return service, nil
			}
		}
		err := errors.New("Unable to determine primary database from multiple. Please bind only one database service to Rails applications.")
		return DBServiceRepresentation{}, err
	}
}

func sanitizeURI(uri string) string {
	re, _ := regexp.Compile(`\/\/.+@`)
	return string(re.ReplaceAll([]byte(uri), []byte("//USER_NAME_PASS@")))
}
