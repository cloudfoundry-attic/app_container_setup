package container

import (
	//	"errors"
	//	warden "github.com/cloudfoundry/gordon"
	. "launchpad.net/gocheck"
)

type MainSuite struct {
}

func init() {
	Suite(&MainSuite{})
}

func (s *MainSuite) TestMainReturnsErrorForInvalidJson(c *C) {
	_, err := Main("")
	c.Assert(err, NotNil)
}

func (s *MainSuite) TestParseForValidJson(c *C) {
	config, err := parseInput(`{
	"disk_limit_in_bytes": 100,
	"memory_limit_in_bytes": 200,
	"warden_socket_path": "/tmp/warden.sock",
	"bind_mounts": [
		{
			"src_path": "/path/src",
			"dst_path": "/path/dst",
			"mode": "ro"
		}
	]}`)
	c.Assert(err, IsNil)
	c.Assert(config.DiskLimitInBytes, Equals, uint64(100))
	c.Assert(config.MemoryLimitInBytes, Equals, uint64(200))
	c.Assert(config.BindMounts[0].SrcPath, Equals, "/path/src")
	c.Assert(config.BindMounts[0].DstPath, Equals, "/path/dst")
	c.Assert(config.BindMounts[0].Mode, Equals, "ro")
	c.Assert(config.WardenSocketPath, Equals, "/tmp/warden.sock")
}

type FakeContainer struct {
	CreateCalls         [][]*BindMount
	SetDiskLimitCalls   []uint64
	SetMemoryLimitCalls []uint64
}

func (c *FakeContainer) Create(bindMounts []*BindMount) error {
	c.CreateCalls = append(c.CreateCalls, bindMounts)
	return nil
}

func (c *FakeContainer) SetDiskLimit(limitInBytes uint64) error {
	c.SetDiskLimitCalls = append(c.SetDiskLimitCalls, limitInBytes)
	return nil
}

func (c *FakeContainer) SetMemoryLimit(limitInBytes uint64) error {
	c.SetMemoryLimitCalls = append(c.SetMemoryLimitCalls, limitInBytes)
	return nil
}

func (s *MainSuite) TestStatePerformingContainerCreation(c *C) {
	fakeContainer := &FakeContainer{}
	state := NewState(fakeContainer,
		&CommandLineJson{DiskLimitInBytes: 123, MemoryLimitInBytes: 456})
	state.Perform()

	c.Assert(len(fakeContainer.CreateCalls) > 0, Equals, true)
	c.Assert(fakeContainer.SetDiskLimitCalls, DeepEquals, []uint64{123})
	c.Assert(fakeContainer.SetMemoryLimitCalls, DeepEquals, []uint64{456})
}

//
//func (s *MainSuite) TestMainComplainingForMissingValues(c *C) {
//	config, err := parseInput(`{
//	"disk_limit_in_bytes": 100
//	}`)
//	c.Assert(err, IsNil)
//	c.Assert(config.DiskLimitInBytes, Equals, int64(100))
//}
