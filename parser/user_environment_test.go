package parser

import (
	. "launchpad.net/gocheck"
)

type UserEnvironmentSuite struct {
}

func init() {
	Suite(&UserEnvironmentSuite{})
}

func (suite *UserEnvironmentSuite) TestEmptyInput(c *C) {
	var input []string

	result := generateUserEnvironment(input)
	c.Assert(len(result), Equals, 0)
}

func (suite *UserEnvironmentSuite) TestGenerateUserEnvironment(c *C) {
	input := []string{"one=1", "sentence=lots of words", "equation= 1 + 1 = 2"}

	result := generateUserEnvironment(input)
	c.Assert(result[0], Equals, EnvironmentPair{"one", "1"})
	c.Assert(result[1], Equals, EnvironmentPair{"sentence", "lots of words"})
	c.Assert(result[2], Equals, EnvironmentPair{"equation", " 1 + 1 = 2"})
}
