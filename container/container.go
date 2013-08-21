package container

import (
	warden "github.com/cloudfoundry/gordon"
)

type Container struct {
	client warden.ConnectedWardenClient
	handle string
}

type ContainerCreator interface {
	Create([]*BindMount) error
	SetDiskLimit(limitInBytes uint64) error
	SetMemoryLimit(limitInBytes uint64) error
	Handle() string
}

type BindMount struct {
	SrcPath string `json:"src_path"`
	DstPath string `json:"dst_path"`
	Mode    string `json:"mode"`
}

func NewContainer(client warden.ConnectedWardenClient) *Container {
	return &Container{client: client}
}

func (c *Container) Create(pathsToBind []*BindMount) error {
	var bindMountRequests []*warden.CreateRequest_BindMount
	ro := warden.CreateRequest_BindMount_RO
	for _, bindMount := range pathsToBind {
		bindMountRequests = append(bindMountRequests,
			&warden.CreateRequest_BindMount{
				SrcPath: &bindMount.SrcPath,
				DstPath: &bindMount.DstPath,
				Mode:    &ro,
			})
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

func (c *Container) SetDiskLimit(limitInBytes uint64) error {
	_, err := c.client.LimitDisk(c.handle, limitInBytes)
	return err
}

func (c *Container) SetMemoryLimit(limitInBytes uint64) error {
	_, err := c.client.LimitMemory(c.handle, limitInBytes)
	return err
}

func (c *Container) ConfigureHomeDirectory() {

}

func (c *Container) Handle() string {
	return c.handle
}
