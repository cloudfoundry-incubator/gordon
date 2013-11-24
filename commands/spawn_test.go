package commands

import (
	"errors"

	. "launchpad.net/gocheck"
)

func (w *CSuite) TestSpawn(c *C) {
	ui := &FakeUI{}
	client := &FakeClient{
		SpawnJobId: uint32(1234),
	}
	command := NewSpawn(client, ui, "a-cool-handle", "/bin/echo woah")

	command.Run()

	c.Assert(client.SpawnedHandle, Equals, "a-cool-handle")
	c.Assert(client.SpawnedScript, Equals, "/bin/echo woah")
	c.Assert(ui.Output, Equals, "1234\n")
	c.Assert(ui.ErrorOutput, Equals, "")
}

func (w *CSuite) TestSpawnWithError(c *C) {
	ui := &FakeUI{}
	client := &FakeClient{
		SpawnError: errors.New("oh no there was an error"),
	}
	command := NewSpawn(client, ui, "a-cool-handle", "/bin/echo woah")

	command.Run()

	c.Assert(ui.Output, Equals, "")
	c.Assert(ui.ErrorOutput, Equals, "oh no there was an error\n")
}
