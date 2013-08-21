package container

import (
	warden "github.com/cloudfoundry/gordon"
)

type Container struct {
	client       warden.ConnectedWardenClient
	handle string
}

func NewContainer(client warden.ConnectedWardenClient) *Container {
	return &Container{client: client}
}

func (c *Container) Create(pathsToBind []string) error {
	var bindMountRequests []*warden.CreateRequest_BindMount
	ro := warden.CreateRequest_BindMount_RO
	for _, path := range pathsToBind {
		bindMountRequests = append(bindMountRequests,
			&warden.CreateRequest_BindMount{SrcPath: &path, DstPath: &path, Mode: &ro})
	}
	request := &warden.CreateRequest{BindMounts: bindMountRequests}
	response, err := c.client.CreateByRequest(request)
	if err != nil {
		return err
	}
	c.handle = response.GetHandle()
	return nil
}

func (c *Container) ConfigureApplicationPorts() {

}

func (c *Container) ConfigureConsolePorts() {

}

func (c *Container) ConfigureDebugPorts() {

}

func (c *Container) SetDiskLimit(limitInBytes uint64) (error) {
	_, err := c.client.LimitDisk(c.handle, limitInBytes)
	return err
}

func (c *Container) SetMemoryLimit(limitInBytes int) {

}

func (c *Container) ConfigureHomeDirectory() {

}
