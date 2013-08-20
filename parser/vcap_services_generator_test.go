package parser

import (
	"encoding/json"

	. "launchpad.net/gocheck"
)

type VcapServicesGeneratorSuite struct {
	*Parser
}

func init() {
	Suite(&VcapServicesGeneratorSuite{
		Parser: NewParser(),
	})
}

func (suite *VcapServicesGeneratorSuite) TestEmptyServicesCase(c *C) {
	result, err := suite.generateServicesJSON([]InputServiceJSON{})
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, "{}")
}

func (suite *VcapServicesGeneratorSuite) TestSerialization(c *C) {
	input1 := InputServiceJSON{
		Name:        "mysql",
		Label:       "rds",
		Tags:        []string{"first-tag", "second-tag"},
		Credentials: map[string]interface{}{"username": "hello", "password": "my-password"},
		Plan:        "large",
		PlanOption:  map[string]interface{}{"speed": "fast"},
	}

	input2 := InputServiceJSON{
		Name:        "mongodb",
		Label:       "document",
		Tags:        []string{"mongo", "doc"},
		Credentials: map[string]interface{}{"username": "admin", "password": "123"},
		Plan:        "small",
		PlanOption:  map[string]interface{}{"speed": "slow"},
	}

	input := []InputServiceJSON{
		input1,
		input2,
	}

	result, err := suite.generateServicesJSON(input)
	c.Assert(err, IsNil)

	output := make(map[string]ServiceJSON)
	json.Unmarshal(result, &output)

	output1 := output["rds"]

	c.Assert(output1.Name, Equals, input1.Name)
	c.Assert(output1.Label, Equals, input1.Label)
	c.Assert(output1.Tags, DeepEquals, input1.Tags)
	c.Assert(output1.Credentials, DeepEquals, input1.Credentials)
	c.Assert(output1.Plan, Equals, input1.Plan)
	c.Assert(output1.PlanOption, DeepEquals, input1.PlanOption)

	output2 := output["document"]

	c.Assert(output2.Name, Equals, input2.Name)
	c.Assert(output2.Label, Equals, input2.Label)
	c.Assert(output2.Tags, DeepEquals, input2.Tags)
	c.Assert(output2.Credentials, DeepEquals, input2.Credentials)
	c.Assert(output2.Plan, Equals, input2.Plan)
	c.Assert(output2.PlanOption, DeepEquals, input2.PlanOption)
}
