package commands

import (
	"github.com/vito/gordon"
)

type CreateCommand struct {
	client warden.Client
	ui     UI
}

func NewCreate(client warden.Client, ui UI) CreateCommand {
	return CreateCommand{
		client: client,
		ui:     ui,
	}
}

func (command CreateCommand) Run() {
	response, err := command.client.Create()
	if err != nil {
		command.ui.Error(err.Error())
		return
	}

	handle := response.GetHandle()
	command.ui.Say(handle)
}
