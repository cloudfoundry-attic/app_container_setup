package parser

import (
	. "launchpad.net/gocheck"
	"testing"
	"regexp"
	"fmt"
)

func Test(t *testing.T) { TestingT(t) }

func StringOrNull(s string) string {
	if s == "" {
		return "null"
	}
	return "\"" + s + "\""
}

//Custom Checkers:

//BetterMatches

var BetterMatches Checker = &betterMatchesChecker{
	&CheckerInfo{Name: "BetterMatches", Params: []string{"value", "regex"}},
}

type betterMatchesChecker struct {
	*CheckerInfo
}

func (checker *betterMatchesChecker) Check(params []interface{}, names []string) (result bool, error string) {
	value := params[0]
	regex := params[1]

	reStr, ok := regex.(string)
	if !ok {
		return false, "Regex must be a string"
	}
	valueStr, valueIsStr := value.(string)
	if !valueIsStr {
		if valueWithStr, valueHasStr := value.(fmt.Stringer); valueHasStr {
			valueStr, valueIsStr = valueWithStr.String(), true
		}
	}
	if valueIsStr {
		matches, err := regexp.MatchString(reStr, valueStr)
		if err != nil {
			return false, "Can't compile regex: " + err.Error()
		}
		return matches, ""
	}
	return false, "Obtained value is not a string and has no .String()"
}

//BetterChecker

type BetterChecker interface {
	Checker
	NegationMessage() string
}

type NegatedChecker struct {
	NegationString string
}

func (checker *NegatedChecker) NegationMessage() string {
	return checker.NegationString
}

//BetterNot

func BetterNot(checker BetterChecker) Checker {
	return &betterNotChecker{checker}
}

type betterNotChecker struct {
	sub BetterChecker
}

func (checker *betterNotChecker) Info() *CheckerInfo {
	info := *checker.sub.Info()
	info.Name = "Not(" + info.Name + ")"
	return &info
}

func (checker *betterNotChecker) Check(params []interface{}, names []string) (result bool, error string) {
	result, error = checker.sub.Check(params, names)
	if (result) {
		return false, checker.sub.NegationMessage()
	}
	return true, ""
}

//Has Key

var HasKey BetterChecker = &hasKeyChecker{
	&CheckerInfo{Name: "HasKey", Params: []string{"map", "key"}},
	&NegatedChecker{},
}

type hasKeyChecker struct {
	*CheckerInfo
	*NegatedChecker
}

func (checker *hasKeyChecker) Check(params []interface{}, names []string) (result bool, error string) {
	mapToCheck := params[0]
	key := params[1]

	keyAsString, ok := key.(string)
	if !ok {
		return false, "key must be a string"
	}

	typecastMapToCheck, ok := mapToCheck.(map[string]string)
	if !ok {
		return false, "map must be a map[string]string"
	}

	_, ok = typecastMapToCheck[keyAsString]
	if (!ok) {
		return false, "Map does not have key " + keyAsString
	} else {
		checker.NegationString = "Map has unexpected key " + keyAsString
	}
	return ok, ""
}
