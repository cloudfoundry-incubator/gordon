package commands

import (
	"fmt"
	"testing"

	. "launchpad.net/gocheck"

	protocol "github.com/vito/gordon/protocol"
)

type CSuite struct{}

func Test(t *testing.T) { TestingT(t) }

func init() {
	Suite(&CSuite{})
}

type FakeClient struct {
	ListHandles []string
	ListError   error

	CreateHandle string
	CreateError  error

	DestroyedHandle string
	DestroyError    error

	SpawnedHandle string
	SpawnedScript string
	SpawnJobId    uint32
	SpawnError    error
}

func (client *FakeClient) List() (*protocol.ListResponse, error) {
	response := &protocol.ListResponse{
		Handles: client.ListHandles,
	}
	return response, client.ListError
}

func (client *FakeClient) Create() (*protocol.CreateResponse, error) {
	response := &protocol.CreateResponse{
		Handle: &client.CreateHandle,
	}
	return response, client.CreateError
}

func (client *FakeClient) Destroy(handle string) (*protocol.DestroyResponse, error) {
	client.DestroyedHandle = handle
	response := &protocol.DestroyResponse{}
	return response, client.DestroyError
}

func (client *FakeClient) Spawn(handle, script string) (*protocol.SpawnResponse, error) {
	client.SpawnedHandle = handle
	client.SpawnedScript = script
	response := &protocol.SpawnResponse{
		JobId: &client.SpawnJobId,
	}
	return response, client.SpawnError
}

type FakeUI struct {
	Output      string
	ErrorOutput string
}

func (ui *FakeUI) Say(message string) {
	ui.Output = fmt.Sprintf("%s%s\n", ui.Output, message)
}

func (ui *FakeUI) Error(message string) {
	ui.ErrorOutput = fmt.Sprintf("%s%s\n", ui.ErrorOutput, message)
}
