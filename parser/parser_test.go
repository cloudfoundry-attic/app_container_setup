package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os/exec"
	"strconv"
	"strings"
)

type InputData struct {
	Debug                        string
	Index                        int
	InstanceGuid                 string
	InstanceContainerPort        int
	InstanceConsoleContainerPort int
	InstanceDebugContainerPort   int
	Services                     string
}

const Mem = 1024
const Index = 2
const InstanceGuid = "BEEF"
const InstanceContainerPort = 1234
const InstanceConsoleContainerPort = 4321
const StartedAtTimestamp = 1376503895
const ApplicationVersion = "2467er728"
const ApplicationName = "simple-app"

func GenerateJSON(input *InputData) string {
	return fmt.Sprintf(`
	{
		"nats_data":{
			"limits":{
				"mem":%d,
				"disk":512,
				"fds":16384
			},
			"debug":%s,
			"index":%d,
			"version":"%s",
			"name":"%s",
			"uris":[
				"simple-app.cfapp.com",
				"other-simple-app.cfapp.com"
			],
			"services":%s
		},
		"instance_guid":"%s",
		"instance_container_port":%d,
		"instance_console_container_port":%d,
		"instance_debug_container_port":%d,
		"started_at_timestamp":%d
	}
	`,
		Mem,
		StringOrNull(input.Debug),
		Index,
		ApplicationVersion,
		ApplicationName,
		input.Services,
		InstanceGuid,
		InstanceContainerPort,
		InstanceConsoleContainerPort,
		input.InstanceDebugContainerPort,
		StartedAtTimestamp,
	)
}

func GenerateServiceJson(data *ServiceData) string {
	return fmt.Sprintf(`
	{
		"credentials":{
			"password":"PASSWORD",
			"uri":"%s"
		},
		"options":{},
		"plan_option":{"what":"is this?"},
		"label":"%s",
		"provider":"%s",
		"version":"%s",
		"vendor":"%s",
		"plan":"%s",
		"name":"%s",
		"tags":["some-tag","some-other-tag"]
	}
	`,
		data.URI,
		data.Label,
		data.Provider,
		data.Version,
		data.Vendor,
		data.Plan,
		data.Name,
	)
}

func GenerateMostlyEmptyServiceJson() string {
	return `
	{
		"label":"label-1"
	}
	`
}

type ServiceData struct {
	Label    string
	Provider string
	Version  string
	Vendor   string
	Plan     string
	Name     string
	URI      string
}

type ParserSuite struct {
	inputData *InputData
}

func init() {
	Suite(&ParserSuite{})
}

func (suite *ParserSuite) SetUpTest(c *C) {
	suite.inputData = new(InputData)
	suite.inputData.Debug = ""
	suite.inputData.InstanceDebugContainerPort = 0
	suite.inputData.Services = "[]"
}

func (suite *ParserSuite) GetEnvironmentVariablesForJSON(json string, c *C) map[string]string {
	parser := NewParser()
	script, err := parser.GenerateEnvironmentScriptFromJSON(json)
	c.Assert(err, IsNil)

	ioutil.WriteFile("/tmp/env_test", []byte(script+"\n env"), 0777)
	output, err := exec.Command("/bin/bash", "/tmp/env_test").Output()
	c.Assert(err, IsNil)

	lines := strings.Split(string(output), "\n")
	result := make(map[string]string)
	for _, line := range lines {
		env_pair := strings.SplitN(line, "=", 2)
		if len(env_pair) == 2 {
			result[env_pair[0]] = env_pair[1]
		}
	}

	return result
}

func (suite *ParserSuite) TestHandlingInvalidJSON(c *C) {
	parser := NewParser()
	output, err := parser.GenerateEnvironmentScriptFromJSON(`kaboom`)
	c.Assert(err, NotNil)
	c.Assert(output, Equals, "")
}

func (suite *ParserSuite) TestExportingBasicSystemEnvironmentVariables(c *C) {
	environment := suite.GetEnvironmentVariablesForJSON(GenerateJSON(suite.inputData), c)

	c.Assert(environment["MEMORY_LIMIT"], Equals, "1024m")
	c.Assert(environment["HOME"], Equals, environment["PWD"]+"/app")
	c.Assert(environment["TMPDIR"], Equals, environment["PWD"]+"/tmp")
	c.Assert(environment["VCAP_APP_HOST"], Equals, "0.0.0.0")
	c.Assert(environment["VCAP_APP_PORT"], Equals, strconv.Itoa(InstanceContainerPort))
	c.Assert(environment["VCAP_CONSOLE_IP"], Equals, "0.0.0.0")
	c.Assert(environment["VCAP_CONSOLE_PORT"], Equals, strconv.Itoa(InstanceConsoleContainerPort))
	c.Assert(environment["PORT"], Equals, strconv.Itoa(InstanceContainerPort))
}

func (suite *ParserSuite) TestDebugEnvironmentVariablesIfSet(c *C) {
	suite.inputData.Debug = "run"
	suite.inputData.InstanceDebugContainerPort = 1235
	environment := suite.GetEnvironmentVariablesForJSON(GenerateJSON(suite.inputData), c)

	c.Assert(environment["VCAP_DEBUG_IP"], Equals, "0.0.0.0")
	c.Assert(environment["VCAP_DEBUG_PORT"], Equals, "1235")
	c.Assert(environment["VCAP_DEBUG_MODE"], Equals, "run")
}

func (suite *ParserSuite) TestDebugEnvironmentVariablesIfNotSet(c *C) {
	environment := suite.GetEnvironmentVariablesForJSON(GenerateJSON(suite.inputData), c)

	c.Assert(environment, HasKey, "MEMORY_LIMIT")
	c.Assert(environment, BetterNot(HasKey), "VCAP_DEBUG_IP")
	c.Assert(environment, BetterNot(HasKey), "VCAP_DEBUG_PORT")
	c.Assert(environment, BetterNot(HasKey), "VCAP_DEBUG_MODE")
}

func (suite *ParserSuite) TestApplicationJsonEnvironmentVariables(c *C) {
	environment := suite.GetEnvironmentVariablesForJSON(GenerateJSON(suite.inputData), c)
	c.Assert(environment, HasKey, "VCAP_APPLICATION")

	var application_json map[string]interface{}
	err := json.Unmarshal([]byte(environment["VCAP_APPLICATION"]), &application_json)
	c.Assert(err, IsNil)

	c.Assert(application_json["instance_id"], Equals, "BEEF")
	c.Assert(int(application_json["instance_index"].(float64)), Equals, 2)
	c.Assert(application_json["host"], Equals, "0.0.0.0")
	c.Assert(int(application_json["port"].(float64)), Equals, InstanceContainerPort)
	c.Assert(int(application_json["port"].(float64)), Equals, InstanceContainerPort)
	c.Assert(int(application_json["started_at_timestamp"].(float64)), Equals, StartedAtTimestamp)
	c.Assert(application_json["started_at"], Equals, "2013-08-14 18:11:35 +0000")
	c.Assert(int(application_json["state_timestamp"].(float64)), Equals, StartedAtTimestamp)
	c.Assert(application_json["start"], Equals, "2013-08-14 18:11:35 +0000")

	limits := (application_json["limits"]).(map[string]interface{})
	c.Assert(int(limits["mem"].(float64)), Equals, 1024)
	c.Assert(int(limits["disk"].(float64)), Equals, 512)
	c.Assert(int(limits["fds"].(float64)), Equals, 16384)

	c.Assert(application_json["version"], Equals, ApplicationVersion)
	c.Assert(application_json["application_version"], Equals, ApplicationVersion)

	c.Assert(application_json["application_name"], Equals, ApplicationName)
	c.Assert(application_json["name"], Equals, ApplicationName)

	c.Assert(application_json["uris"], DeepEquals, []interface{}{"simple-app.cfapp.com", "other-simple-app.cfapp.com"})

	c.Assert(application_json["users"], IsNil)
}

func (suite *ParserSuite) TestServicesJsonEnvironmentVariablesWithNoServices(c *C) {
	environment := suite.GetEnvironmentVariablesForJSON(GenerateJSON(suite.inputData), c)
	c.Assert(environment, HasKey, "VCAP_SERVICES")
	c.Assert(environment["VCAP_SERVICES"], Equals, "{}")

	var services_json map[string]interface{}
	err := json.Unmarshal([]byte(environment["VCAP_SERVICES"]), &services_json)
	c.Assert(err, IsNil)
	c.Assert(len(services_json), Equals, 0)
}

func (suite *ParserSuite) TestServicesJsonEnvironmentVariablesWithMultipleService(c *C) {
	service1 := &ServiceData{Label: "label-1", Provider: "provider-1", Version: "version-1", Vendor: "vendor-1", Plan: "plan-1", Name: "name-1", URI: "http://foo.com"}
	service2 := &ServiceData{Label: "label-2", Provider: "provider-2", Version: "version-2", Vendor: "vendor-2", Plan: "plan-2", Name: "name-2", URI: "http://foo.com"}

	suite.inputData.Services = fmt.Sprintf("[%s,%s]", GenerateServiceJson(service1), GenerateServiceJson(service2))

	environment := suite.GetEnvironmentVariablesForJSON(GenerateJSON(suite.inputData), c)
	c.Assert(environment, HasKey, "VCAP_SERVICES")

	var services_json map[string]interface{}
	err := json.Unmarshal([]byte(environment["VCAP_SERVICES"]), &services_json)
	c.Assert(err, IsNil)
	c.Assert(len(services_json), Equals, 2)

	for index := 1; index <= 2; index++ {
		i := strconv.Itoa(index)
		c.Assert(services_json["label-"+i], NotNil)
		service_json := services_json["label-"+i].(map[string]interface{})
		c.Assert(service_json["name"], Equals, "name-"+i)
		c.Assert(service_json["label"], Equals, "label-"+i)
		c.Assert(service_json["tags"], DeepEquals, []interface{}{"some-tag", "some-other-tag"})
		c.Assert(service_json["credentials"], DeepEquals, map[string]interface{}{"password": "PASSWORD", "uri": "http://foo.com"})
		c.Assert(service_json["plan_option"], DeepEquals, map[string]interface{}{"what": "is this?"})
		c.Assert(service_json["plan"], Equals, "plan-"+i)
	}
}

func (suite *ParserSuite) TestDatabaseURLEnvironmentVariablesWithNoServices(c *C) {
	environment := suite.GetEnvironmentVariablesForJSON(GenerateJSON(suite.inputData), c)
	c.Assert(environment, BetterNot(HasKey), "DATABASE_URL")
}

func (suite *ParserSuite) TestDatabaseURLEnvironmentVariableWithServices(c *C) {
	service1 := &ServiceData{Label: "foo", Provider: "provider-1", Version: "version-1", Vendor: "vendor-1", Plan: "plan-1", Name: "foo_production", URI: "postgresql://a:b@foo.com?q=2"}
	service2 := &ServiceData{Label: "bar", Provider: "provider-1", Version: "version-1", Vendor: "vendor-1", Plan: "plan-1", Name: "bar", URI: "mysql://a:b@bar.com?q=2"}
	suite.inputData.Services = fmt.Sprintf("[%s,%s]", GenerateServiceJson(service1), GenerateServiceJson(service2))
	environment := suite.GetEnvironmentVariablesForJSON(GenerateJSON(suite.inputData), c)
	c.Assert(environment["DATABASE_URL"], Equals, "postgres://a:b@foo.com?q=2")
}

//todo: test profile.d stuff
//todo: test user environments
