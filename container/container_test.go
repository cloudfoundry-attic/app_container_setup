package container

import (
	"errors"
	warden "github.com/cloudfoundry/gordon"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ContainerSuite struct {
}

func init() {
	Suite(&ContainerSuite{})
}

func (suite *ContainerSuite) TestCreate(c *C) {
	var request *warden.CreateRequest
	fakeClient := MakeFakeWardenClient()
	fakeClient.CreateByRequestFunc = func(r *warden.CreateRequest) (*warden.CreateResponse, error) {
		request = r
		handle := "wardenhandle"
		return &warden.CreateResponse{Handle: &handle}, nil
	}
	container := NewContainer(fakeClient)

	err := container.Create([]string{"/tmp/foo"})

	c.Assert(len(request.GetBindMounts()), Equals, 1)
	bindMount := request.GetBindMounts()[0]
	c.Assert(*bindMount.SrcPath, Equals, "/tmp/foo")
	c.Assert(*bindMount.DstPath, Equals, "/tmp/foo")
	c.Assert(*bindMount.Mode, Equals, warden.CreateRequest_BindMount_RO)

	c.Assert(err, IsNil)
	c.Assert(container.handle, Equals, "wardenhandle")
}

func (suite *ContainerSuite) TestCreateError(c *C) {
	fakeClient := MakeFakeWardenClient()
	fakeClient.CreateByRequestFunc = func(*warden.CreateRequest) (*warden.CreateResponse, error) {
		return nil, errors.New("no client for you")
	}
	container := NewContainer(fakeClient)

	err := container.Create([]string{})

	c.Assert(err.Error(), Equals, "no client for you")
	c.Assert(container.handle, Equals, "")
}

func (suite *ContainerSuite) TestSetDiskLimit(c *C) {
	var handle string
	var limit uint64
	fakeClient := MakeFakeWardenClient()
	fakeClient.LimitDiskFunc = func(h string, l uint64) (*warden.LimitDiskResponse, error) {
		handle = h
		limit = l
		return nil, nil
	}

	container := NewContainer(fakeClient)
	container.handle = "the_warden_handle"

	container.SetDiskLimit(123)

	c.Assert(handle, Equals, "the_warden_handle")
	c.Assert(limit, Equals, uint64(123))
}

func (suite *ContainerSuite) TestSetDiskLimitError(c *C) {
	fakeClient := MakeFakeWardenClient()
	fakeClient.LimitDiskFunc = func(h string, l uint64) (*warden.LimitDiskResponse, error) {
		return nil, errors.New("failed to limit disk")
	}

	container := NewContainer(fakeClient)
	container.handle = "the_warden_handle"

	err := container.SetDiskLimit(123)
	c.Assert(err.Error(), Equals, "failed to limit disk")
}

type fakeWardenClient struct {
	CreateByRequestFunc func(*warden.CreateRequest) (*warden.CreateResponse, error)
	LimitDiskFunc func(string, uint64) (*warden.LimitDiskResponse, error)
}

func MakeFakeWardenClient() *fakeWardenClient {
	return &fakeWardenClient{
		CreateByRequestFunc: func(*warden.CreateRequest) (*warden.CreateResponse, error) { return nil, nil },
		LimitDiskFunc: func(string, uint64) (*warden.LimitDiskResponse, error) { return nil, nil },
	}
}

func (c *fakeWardenClient) CreateByRequest(r *warden.CreateRequest) (*warden.CreateResponse, error) {
	return c.CreateByRequestFunc(r)
}

func (c *fakeWardenClient) LimitDisk(handle string, limit uint64) (*warden.LimitDiskResponse, error) {
	return c.LimitDiskFunc(handle, limit)
}
