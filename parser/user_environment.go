package parser

import (
	"strings"
)

func generateUserEnvironment(input []string) []EnvironmentPair {
	result := make([]EnvironmentPair, len(input))

	for i, pair := range input {
		split := strings.SplitN(pair, "=", 2)
		result[i] = EnvironmentPair{split[0], split[1]}
	}

	return result
}
