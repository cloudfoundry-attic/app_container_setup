package container

import (
	"encoding/json"
	warden "github.com/cloudfoundry/gordon"
)

type CommandLineJson struct {
	DiskLimitInBytes   uint64       `json:"disk_limit_in_bytes",required`
	MemoryLimitInBytes uint64       `json:"memory_limit_in_bytes"`
	BindMounts         []*BindMount `json:"bind_mounts"`
	WardenSocketPath   string       `json:"warden_socket_path"`
}

type State struct {
	Container       ContainerCreator
	CommandLineJson *CommandLineJson
}

func NewState(container ContainerCreator, commandLineJson *CommandLineJson) *State {
	return &State{Container: container, CommandLineJson: commandLineJson}
}

func Main(inputJson string) (*State, error) {
	commandLineJson, err := parseInput(inputJson)

	connectionInfo := &warden.ConnectionInfo{commandLineJson.WardenSocketPath}
	container := NewContainer(warden.NewClient(connectionInfo))

	state := &State{CommandLineJson: commandLineJson, Container: container}
	//	state.Perform
	return state, err
}

func parseInput(inputJson string) (*CommandLineJson, error) {
	var input CommandLineJson
	err := json.Unmarshal([]byte(inputJson), &input)

	return &input, err
}

func (s *State) Perform() {
	s.Container.Create(s.CommandLineJson.BindMounts)
	s.Container.SetDiskLimit(s.CommandLineJson.DiskLimitInBytes)
	s.Container.SetMemoryLimit(s.CommandLineJson.MemoryLimitInBytes)
}

func (c *CommandLineJson) IsValid() bool {
	// TODO: logic
	return false
}
