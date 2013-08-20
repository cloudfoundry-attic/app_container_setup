package parser

import (
	"encoding/json"
	. "launchpad.net/gocheck"
)

type VcapApplicationGeneratorSuite struct {
	*Parser
}

func init() {
	Suite(&VcapApplicationGeneratorSuite{
		Parser: NewParser(),
	})
}

func (suite *VcapApplicationGeneratorSuite) TestSerialization(c *C) {
	limitsJson := InputNatsLimitsJSON{
		Mem:  1024,
		Disk: 512,
		Fds:  16384,
	}
	natsJson := InputNatsDataJSON{
		Index:              1,
		Limits:             limitsJson,
		ApplicationVersion: "efsk234ee",
		Name:               "myapp",
		Uris:               []string{"simple-app.cfapp.com", "other-simple-app.cfapp.com"},
	}
	inputJson := InputJSON{
		NatsData:              natsJson,
		InstanceContainerPort: 5555,
		InstanceGuid:          "2467er728",
		StartedAtTimestamp:    1376503895,
	}

	result, err := suite.generateApplicationJSON(inputJson)
	c.Assert(err, IsNil)

	var output map[string]interface{}
	json.Unmarshal(result, &output)

	c.Assert(output["instance_id"], Equals, inputJson.InstanceGuid)
	c.Assert(int(output["instance_index"].(float64)), Equals, natsJson.Index)
	c.Assert(output["host"], Equals, "0.0.0.0")
	c.Assert(int(output["port"].(float64)), Equals, inputJson.InstanceContainerPort)
	c.Assert(int(output["started_at_timestamp"].(float64)), Equals, int(inputJson.StartedAtTimestamp))
	c.Assert(int(output["state_timestamp"].(float64)), Equals, int(inputJson.StartedAtTimestamp))
	c.Assert(output["start"], Equals, "2013-08-14 18:11:35 +0000")

	limits := (output["limits"]).(map[string]interface{})
	c.Assert(int(limits["mem"].(float64)), Equals, limitsJson.Mem)
	c.Assert(int(limits["disk"].(float64)), Equals, limitsJson.Disk)
	c.Assert(int(limits["fds"].(float64)), Equals, limitsJson.Fds)

	c.Assert(output["version"], Equals, natsJson.ApplicationVersion)
	c.Assert(output["application_version"], Equals, natsJson.ApplicationVersion)

	c.Assert(output["application_name"], Equals, natsJson.Name)
	c.Assert(output["name"], Equals, natsJson.Name)

	c.Assert(output["uris"], DeepEquals, []interface{}{"simple-app.cfapp.com", "other-simple-app.cfapp.com"})

	c.Assert(output["users"], IsNil)
}
