package commands

import (
	"errors"

	. "launchpad.net/gocheck"
)

func (w *CSuite) TestListWithNoContainers(c *C) {
	ui := &FakeUI{}
	client := &FakeClient{
		ListHandles: []string{},
	}
	command := NewList(client, ui)

	command.Run()

	c.Assert(ui.Output, Equals, "")
	c.Assert(ui.ErrorOutput, Equals, "")
}

func (w *CSuite) TestListWithContainers(c *C) {
	ui := &FakeUI{}
	client := &FakeClient{
		ListHandles: []string{"some", "awesome", "handles"},
	}
	command := NewList(client, ui)

	command.Run()

	c.Assert(ui.Output, Equals, "some\nawesome\nhandles\n")
	c.Assert(ui.ErrorOutput, Equals, "")
}

func (w *CSuite) TestListWithError(c *C) {
	ui := &FakeUI{}
	client := &FakeClient{
		ListError: errors.New("oh no there was an error"),
	}
	command := NewList(client, ui)

	command.Run()

	c.Assert(ui.Output, Equals, "")
	c.Assert(ui.ErrorOutput, Equals, "oh no there was an error\n")
}
