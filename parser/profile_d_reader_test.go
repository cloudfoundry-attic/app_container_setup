package parser

import (
	. "launchpad.net/gocheck"
)

type ProfileDReaderSuite struct {
}

func init() {
	Suite(&ProfileDReaderSuite{})
}

func (suite *ProfileDReaderSuite) TestScript(c *C) {
	output := generateProfileDReader()
	c.Assert(output, Equals, `
unset GEM_PATH
if [ -d app/.profile.d ]; then
for i in app/.profile.d/*.sh; do
  if [ -r $i ]; then
	. $i
  fi
done
unset i
fi
`)
}
