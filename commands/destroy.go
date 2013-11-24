package commands

import (
	"github.com/vito/gordon"
)

type DestroyCommand struct {
	client warden.Client
	ui     UI

	handle string
}

func NewDestroy(client warden.Client, ui UI, handle string) DestroyCommand {
	return DestroyCommand{
		client: client,
		ui:     ui,

		handle: handle,
	}
}

func (command DestroyCommand) Run() {
	_, err := command.client.Destroy(command.handle)
	if err != nil {
		command.ui.Error(err.Error())
		return
	}
}
