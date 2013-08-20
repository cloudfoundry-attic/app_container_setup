package parser

import (
	. "launchpad.net/gocheck"
)

type DatabaseUriGeneratorSuite struct{}

func init() {
	Suite(&DatabaseUriGeneratorSuite{})
}

func NewDBServiceRepresentation(name string, uri string) DBServiceRepresentation {
	return DBServiceRepresentation{Name: name, URI: uri}
}

func (suite *DatabaseUriGeneratorSuite) TestEmptyServicesCase(c *C) {
	uri, err := NewDatabaseURIGenerator([]DBServiceRepresentation{}).Generate()
	c.Assert(err, IsNil)
	c.Assert(uri, Equals, "")
}

func (suite *DatabaseUriGeneratorSuite) TestNoBoundRelationalValidDatabases(c *C) {
	generator := NewDatabaseURIGenerator([]DBServiceRepresentation{NewDBServiceRepresentation("foo", "sendgrid://foo.com")})
	uri, err := generator.Generate()
	c.Assert(err, IsNil)
	c.Assert(uri, Equals, "")
}

func (suite *DatabaseUriGeneratorSuite) TestEmptyUri(c *C) {
	generator := NewDatabaseURIGenerator([]DBServiceRepresentation{NewDBServiceRepresentation("foo", "")})
	uri, err := generator.Generate()
	c.Assert(err, IsNil)
	c.Assert(uri, Equals, "")
}

func (suite *DatabaseUriGeneratorSuite) TestInvalidUri(c *C) {
	generator := NewDatabaseURIGenerator([]DBServiceRepresentation{NewDBServiceRepresentation("foo", "://user:s3cr3t@oops")})
	uri, err := generator.Generate()

	c.Assert(err.Error(), Equals, "Invalid database URI \"://USER_NAME_PASS@oops\"")
	c.Assert(uri, Equals, "")
}

func (suite *DatabaseUriGeneratorSuite) TestMysql(c *C) {
	generator := NewDatabaseURIGenerator([]DBServiceRepresentation{
		NewDBServiceRepresentation("foo", "mysql://user@pass:foo.com/path/to/db?q=hi"),
		NewDBServiceRepresentation("bar", "bar://bar.com"),
	})
	uri, err := generator.Generate()
	c.Assert(err, IsNil)
	c.Assert(uri, Equals, "mysql2://user@pass:foo.com/path/to/db?q=hi")
}

func (suite *DatabaseUriGeneratorSuite) TestMysql2(c *C) {
	generator := NewDatabaseURIGenerator([]DBServiceRepresentation{
		NewDBServiceRepresentation("foo", "mysql2://user@pass:foo.com/path/to/db?q=hi"),
		NewDBServiceRepresentation("bar", "bar://bar.com"),
	})
	uri, err := generator.Generate()
	c.Assert(err, IsNil)
	c.Assert(uri, Equals, "mysql2://user@pass:foo.com/path/to/db?q=hi")
}

func (suite *DatabaseUriGeneratorSuite) TestPostgres(c *C) {
	generator := NewDatabaseURIGenerator([]DBServiceRepresentation{
		NewDBServiceRepresentation("foo", "postgres://user@pass:foo.com/path/to/db?q=hi"),
		NewDBServiceRepresentation("bar", "bar://bar.com"),
	})
	uri, err := generator.Generate()
	c.Assert(err, IsNil)
	c.Assert(uri, Equals, "postgres://user@pass:foo.com/path/to/db?q=hi")
}

func (suite *DatabaseUriGeneratorSuite) TestPostgreSQL(c *C) {
	generator := NewDatabaseURIGenerator([]DBServiceRepresentation{
		NewDBServiceRepresentation("foo", "postgresql://user@pass:foo.com/path/to/db?q=hi"),
		NewDBServiceRepresentation("bar", "bar://bar.com"),
	})
	uri, err := generator.Generate()
	c.Assert(err, IsNil)
	c.Assert(uri, Equals, "postgres://user@pass:foo.com/path/to/db?q=hi")
}

func (suite *DatabaseUriGeneratorSuite) TestProductionDatabaseName(c *C) {
	generator := NewDatabaseURIGenerator([]DBServiceRepresentation{
		NewDBServiceRepresentation("foo", "mysql://alpha.com"),
		NewDBServiceRepresentation("fooproduction", "mysql://user@pass:beta.com/path/to/db?q=hi"),
		NewDBServiceRepresentation("bar", "bar://bar.com"),
	})
	uri, err := generator.Generate()
	c.Assert(err, IsNil)
	c.Assert(uri, Equals, "mysql2://user@pass:beta.com/path/to/db?q=hi")
}

func (suite *DatabaseUriGeneratorSuite) TestProdDatabaseName(c *C) {
	generator := NewDatabaseURIGenerator([]DBServiceRepresentation{
		NewDBServiceRepresentation("foo", "mysql://alpha.com"),
		NewDBServiceRepresentation("fooprod", "mysql://user@pass:beta.com/path/to/db?q=hi"),
		NewDBServiceRepresentation("bar", "bar://bar.com"),
	})
	uri, err := generator.Generate()
	c.Assert(err, IsNil)
	c.Assert(uri, Equals, "mysql2://user@pass:beta.com/path/to/db?q=hi")
}

func (suite *DatabaseUriGeneratorSuite) TestMultipleProductionDatabaseNames(c *C) {
	generator := NewDatabaseURIGenerator([]DBServiceRepresentation{
		NewDBServiceRepresentation("foo", "mysql://alpha.com"),
		NewDBServiceRepresentation("fooprod", "mysql://user@pass:beta.com/path/to/db?q=hi"),
		NewDBServiceRepresentation("fooproduction", "mysql://user@pass:gamma.com/path/to/db?q=hi"),
		NewDBServiceRepresentation("bar", "bar://bar.com"),
	})
	uri, err := generator.Generate()
	c.Assert(err, IsNil)
	c.Assert(uri, Equals, "mysql2://user@pass:beta.com/path/to/db?q=hi")
}

func (suite *DatabaseUriGeneratorSuite) TestUnidentifiableProdDatabase(c *C) {
	generator := NewDatabaseURIGenerator([]DBServiceRepresentation{
		NewDBServiceRepresentation("foo", "mysql://alpha.com"),
		NewDBServiceRepresentation("foopro", "mysql://user@pass:beta.com/path/to/db?q=hi"),
		NewDBServiceRepresentation("fooproduction", "sendgrid://user@pass:gamma.com/path/to/db?q=hi"),
		NewDBServiceRepresentation("bar", "bar://bar.com"),
	})
	uri, err := generator.Generate()
	c.Assert(err.Error(), Equals, "Unable to determine primary database from multiple. Please bind only one database service to Rails applications.")
	c.Assert(uri, Equals, "")
}
