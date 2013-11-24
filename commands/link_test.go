package commands

import (
	"errors"

	. "launchpad.net/gocheck"
)

func (w *CSuite) TestLink(c *C) {
	ui := &FakeUI{}
	client := &FakeClient{
		LinkStdout:     "stdout output",
		LinkStderr:     "stderr output",
		LinkExitStatus: 0,
	}
	command := NewLink(client, ui, "a-linked-handle", 1)

	command.Run()

	c.Assert(client.LinkedHandle, Equals, "a-linked-handle")
	c.Assert(client.LinkedJobId, Equals, uint32(1))
	c.Assert(ui.Output, Equals, "status: 0\n\nstdout:\nstdout output\n\nstderr:\nstderr output\n")
	c.Assert(ui.ErrorOutput, Equals, "")
}

func (w *CSuite) TestLinkWithError(c *C) {
	ui := &FakeUI{}
	client := &FakeClient{
		LinkError: errors.New("oh no there was an error"),
	}
	command := NewLink(client, ui, "a-linked-handle", 1)

	command.Run()

	c.Assert(ui.Output, Equals, "")
	c.Assert(ui.ErrorOutput, Equals, "oh no there was an error\n")
}
