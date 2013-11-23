package commands

import (
	"bytes"

	. "launchpad.net/gocheck"
)

func (w *CSuite) TestBasicSay(c *C) {
	writer := &bytes.Buffer{}
	basic := BasicUI{
		Writer: writer,
	}

	basic.Say("a fancy string")
	c.Assert(writer.String(), Equals, "a fancy string\n")
}

func (w *CSuite) TestBasicError(c *C) {
	writer := &bytes.Buffer{}
	basic := BasicUI{
		Writer: writer,
	}

	basic.Error("a fancy error string")
	c.Assert(writer.String(), Equals, "a fancy error string\n")
}
