package commands

import (
	"fmt"

	"github.com/vito/gordon"
)

type SpawnCommand struct {
	client warden.Client
	ui     UI

	handle string
	script string
}

func NewSpawn(client warden.Client, ui UI, handle string, script string) SpawnCommand {
	return SpawnCommand{
		client: client,
		ui:     ui,

		handle: handle,
		script: script,
	}
}

func (command SpawnCommand) Run() {
	response, err := command.client.Spawn(command.handle, command.script)
	if err != nil {
		command.ui.Error(err.Error())
		return
	}

	jobId := response.GetJobId()
	command.ui.Say(fmt.Sprintf("%d", jobId))
}
