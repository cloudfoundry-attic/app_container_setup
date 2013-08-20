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

	var output map[string]interface{}
	json.Unmarshal(result, &output)

	output1 := output["rds"].(map[string]interface{})

	c.Assert(output1["name"], Equals, input1.Name)
	c.Assert(output1["label"], Equals, input1.Label)
	c.Assert(output1["tags"], DeepEquals, []interface{}{"first-tag", "second-tag"})
	c.Assert(output1["credentials"], DeepEquals, input1.Credentials)
	c.Assert(output1["plan"], Equals, input1.Plan)
	c.Assert(output1["plan_option"], DeepEquals, input1.PlanOption)

	output2 := output["document"].(map[string]interface{})

	c.Assert(output2["name"], Equals, input2.Name)
	c.Assert(output2["label"], Equals, input2.Label)
	c.Assert(output2["tags"], DeepEquals, []interface{}{"mongo", "doc"})
	c.Assert(output2["credentials"], DeepEquals, input2.Credentials)
	c.Assert(output2["plan"], Equals, input2.Plan)
	c.Assert(output2["plan_option"], DeepEquals, input2.PlanOption)
}

func (suite *VcapServicesGeneratorSuite) TestFailIfMissingLabel(c *C) {
	input := InputServiceJSON{Name: "some-name"}
	result, err := suite.generateServicesJSON([]InputServiceJSON{input})
	c.Assert(result, IsNil)
	c.Assert(err, Equals, ErrMissingLabel)
}

func (suite *VcapServicesGeneratorSuite) TestMissingFields(c *C) {
	input := InputServiceJSON{Label: "mysql"}
	result, err := suite.generateServicesJSON([]InputServiceJSON{input})
	c.Assert(err, IsNil)
	c.Assert(string(result), Equals, `{"mysql":{"label":"mysql"}}`)
}
