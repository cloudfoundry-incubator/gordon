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
}

func (client *FakeClient) List() (*protocol.ListResponse, error) {
	response := &protocol.ListResponse{
		Handles: client.ListHandles,
	}
	return response, client.ListError
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
