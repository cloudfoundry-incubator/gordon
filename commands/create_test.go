package commands

import (
	"errors"

	. "launchpad.net/gocheck"
)

func (w *CSuite) TestCreate(c *C) {
	ui := &FakeUI{}
	client := &FakeClient{
		CreateHandle: "container-name",
	}
	command := NewCreate(client, ui)

	command.Run()

	c.Assert(ui.Output, Equals, "container-name\n")
	c.Assert(ui.ErrorOutput, Equals, "")
}

func (w *CSuite) TestCreateWithError(c *C) {
	ui := &FakeUI{}
	client := &FakeClient{
		CreateError: errors.New("oh no there was an error"),
	}
	command := NewCreate(client, ui)

	command.Run()

	c.Assert(ui.Output, Equals, "")
	c.Assert(ui.ErrorOutput, Equals, "oh no there was an error\n")
}
