package commands

import (
	"github.com/vito/gordon"
)

type ListCommand struct {
	client warden.Client
	ui     UI
}

func NewList(client warden.Client, ui UI) ListCommand {
	return ListCommand{
		client: client,
		ui:     ui,
	}
}

func (command ListCommand) Run() {
	response, err := command.client.List()
	if err != nil {
		command.ui.Error(err.Error())
		return
	}

	for _, handle := range response.GetHandles() {
		command.ui.Say(handle)
	}
}
