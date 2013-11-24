package commands

import (
	"fmt"

	"github.com/vito/gordon"
)

type LinkCommand struct {
	client warden.Client
	ui     UI

	handle string
	jobId  uint32
}

func NewLink(client warden.Client, ui UI, handle string, jobId uint32) LinkCommand {
	return LinkCommand{
		client: client,
		ui:     ui,

		handle: handle,
		jobId:  jobId,
	}
}

func (command LinkCommand) Run() {
	response, err := command.client.Link(command.handle, command.jobId)
	if err != nil {
		command.ui.Error(err.Error())
		return
	}

	output := fmt.Sprintf("status: %d\n\nstdout:\n%s\n\nstderr:\n%s",
		response.GetExitStatus(),
		response.GetStdout(),
		response.GetStderr())
	command.ui.Say(output)
}
