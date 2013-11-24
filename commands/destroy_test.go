package commands

import (
	"errors"

	. "launchpad.net/gocheck"
)

func (w *CSuite) TestDestroy(c *C) {
	ui := &FakeUI{}
	client := &FakeClient{}
	command := NewDestroy(client, ui, "a-destroyed-handle")

	command.Run()

	c.Assert(client.DestroyedHandle, Equals, "a-destroyed-handle")
	c.Assert(ui.Output, Equals, "")
	c.Assert(ui.ErrorOutput, Equals, "")
}

func (w *CSuite) TestDestroyWithError(c *C) {
	ui := &FakeUI{}
	client := &FakeClient{
		DestroyError: errors.New("oh no there was an error"),
	}
	command := NewDestroy(client, ui, "a-destroyed-handle")

	command.Run()

	c.Assert(ui.Output, Equals, "")
	c.Assert(ui.ErrorOutput, Equals, "oh no there was an error\n")
}
