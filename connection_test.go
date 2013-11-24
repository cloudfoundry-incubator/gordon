package warden

import (
	"bytes"
	"math"

	"code.google.com/p/goprotobuf/proto"
	. "launchpad.net/gocheck"

	protocol "github.com/vito/gordon/protocol"
)

func (w *WSuite) TestConnectionCreating(c *C) {
	conn := &fakeConn{
		ReadBuffer: messages(&protocol.CreateResponse{
			Handle: proto.String("foohandle"),
		}),

		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	resp, err := connection.Create()
	c.Assert(err, IsNil)

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(messages(&protocol.CreateRequest{}).Bytes()),
	)

	c.Assert(resp.GetHandle(), Equals, "foohandle")
}

func (w *WSuite) TestConnectionDestroying(c *C) {
	conn := &fakeConn{
		ReadBuffer:  messages(&protocol.DestroyResponse{}),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	_, err := connection.Destroy("foo")
	c.Assert(err, IsNil)

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(messages(&protocol.DestroyRequest{Handle: proto.String("foo")}).Bytes()),
	)
}

func (w *WSuite) TestMemoryLimiting(c *C) {
	conn := &fakeConn{
		ReadBuffer:  messages(&protocol.LimitMemoryResponse{LimitInBytes: proto.Uint64(40)}),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	res, err := connection.LimitMemory("foo", 42)
	c.Assert(err, IsNil)

	c.Assert(res.GetLimitInBytes(), Equals, uint64(40))

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(
			messages(
				&protocol.LimitMemoryRequest{
					Handle:       proto.String("foo"),
					LimitInBytes: proto.Uint64(42),
				},
			).Bytes(),
		),
	)
}

func (w *WSuite) TestGettingMemoryLimit(c *C) {
	conn := &fakeConn{
		ReadBuffer:  messages(&protocol.LimitMemoryResponse{LimitInBytes: proto.Uint64(40)}),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	memoryLimit, err := connection.GetMemoryLimit("foo")
	c.Assert(err, IsNil)
	c.Assert(memoryLimit, Equals, uint64(40))

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(
			messages(
				&protocol.LimitMemoryRequest{
					Handle: proto.String("foo"),
				},
			).Bytes(),
		),
	)
}

func (w *WSuite) TestGettingMemoryLimitThatLooksFishy(c *C) {
	conn := &fakeConn{
		ReadBuffer:  messages(&protocol.LimitMemoryResponse{LimitInBytes: proto.Uint64(math.MaxInt64)}),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	memoryLimit, err := connection.GetMemoryLimit("foo")
	c.Assert(err, IsNil)
	c.Assert(memoryLimit, Equals, uint64(0))

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(
			messages(
				&protocol.LimitMemoryRequest{
					Handle: proto.String("foo"),
				},
			).Bytes(),
		),
	)
}

func (w *WSuite) TestDiskLimiting(c *C) {
	conn := &fakeConn{
		ReadBuffer:  messages(&protocol.LimitDiskResponse{ByteLimit: proto.Uint64(40)}),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	res, err := connection.LimitDisk("foo", 42)
	c.Assert(err, IsNil)

	c.Assert(res.GetByteLimit(), Equals, uint64(40))

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(
			messages(
				&protocol.LimitDiskRequest{
					Handle:    proto.String("foo"),
					ByteLimit: proto.Uint64(42),
				},
			).Bytes(),
		),
	)
}

func (w *WSuite) TestGettingDiskLimit(c *C) {
	conn := &fakeConn{
		ReadBuffer:  messages(&protocol.LimitDiskResponse{ByteLimit: proto.Uint64(40)}),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	diskLimit, err := connection.GetDiskLimit("foo")
	c.Assert(err, IsNil)
	c.Assert(diskLimit, Equals, uint64(40))

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(
			messages(
				&protocol.LimitDiskRequest{
					Handle: proto.String("foo"),
				},
			).Bytes(),
		),
	)
}

func (w *WSuite) TestConnectionSpawn(c *C) {
	conn := &fakeConn{
		ReadBuffer:  messages(&protocol.SpawnResponse{JobId: proto.Uint32(42)}),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	resp, err := connection.Spawn("foo-handle", "echo hi")
	c.Assert(err, IsNil)

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(messages(&protocol.SpawnRequest{
			Handle: proto.String("foo-handle"),
			Script: proto.String("echo hi"),
		}).Bytes()),
	)

	c.Assert(resp.GetJobId(), Equals, uint32(42))
}

func (w *WSuite) TestConnectionNetIn(c *C) {
	conn := &fakeConn{
		ReadBuffer: messages(
			&protocol.NetInResponse{
				HostPort:      proto.Uint32(7331),
				ContainerPort: proto.Uint32(7331),
			},
		),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	resp, err := connection.NetIn("foo-handle")
	c.Assert(err, IsNil)

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(messages(&protocol.NetInRequest{
			Handle: proto.String("foo-handle"),
		}).Bytes()),
	)

	c.Assert(resp.GetHostPort(), Equals, uint32(7331))
	c.Assert(resp.GetContainerPort(), Equals, uint32(7331))
}

func (w *WSuite) TestConnectionList(c *C) {
	conn := &fakeConn{
		ReadBuffer: messages(
			&protocol.ListResponse{
				Handles: []string{"container1", "container2", "container3"},
			},
		),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	resp, err := connection.List()
	c.Assert(err, IsNil)

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(messages(&protocol.ListRequest{}).Bytes()),
	)

	c.Assert(resp.GetHandles(), DeepEquals, []string{"container1", "container2", "container3"})
}

func (w *WSuite) TestConnectionLink(c *C) {
	conn := &fakeConn{
		ReadBuffer: messages(
			&protocol.LinkResponse{
				ExitStatus: proto.Uint32(0),
				Stdout:     proto.String("stdout output"),
				Stderr:     proto.String("stderr output"),
				Info: &protocol.InfoResponse{
					State: proto.String("active"),
				}},
		),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	resp, err := connection.Link("handle", 1)
	c.Assert(err, IsNil)

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(messages(&protocol.LinkRequest{
			Handle: proto.String("handle"),
			JobId:  proto.Uint32(1),
		}).Bytes()),
	)

	c.Assert(resp.GetExitStatus(), Equals, uint32(0))
	c.Assert(resp.GetStdout(), Equals, "stdout output")
	c.Assert(resp.GetStderr(), Equals, "stderr output")
	c.Assert(resp.GetInfo().GetState(), Equals, "active")
}

func (w *WSuite) TestConnectionInfo(c *C) {
	conn := &fakeConn{
		ReadBuffer: messages(
			&protocol.InfoResponse{
				State: proto.String("active"),
			},
		),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	resp, err := connection.Info("handle")
	c.Assert(err, IsNil)

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(messages(&protocol.InfoRequest{
			Handle: proto.String("handle"),
		}).Bytes()),
	)

	c.Assert(resp.GetState(), Equals, "active")
}

func (w *WSuite) TestConnectionCopyIn(c *C) {
	conn := &fakeConn{
		ReadBuffer:  messages(&protocol.CopyInResponse{}),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	_, err := connection.CopyIn("foo-handle", "/foo", "/bar")
	c.Assert(err, IsNil)

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(messages(&protocol.CopyInRequest{
			Handle:  proto.String("foo-handle"),
			SrcPath: proto.String("/foo"),
			DstPath: proto.String("/bar"),
		}).Bytes()),
	)
}

func (w *WSuite) TestConnectionRun(c *C) {
	conn := &fakeConn{
		ReadBuffer:  messages(&protocol.RunResponse{ExitStatus: proto.Uint32(137)}),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	resp, err := connection.Run("foo-handle", "echo hi")
	c.Assert(err, IsNil)

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(messages(&protocol.RunRequest{
			Handle: proto.String("foo-handle"),
			Script: proto.String("echo hi"),
		}).Bytes()),
	)

	c.Assert(resp.GetExitStatus(), Equals, uint32(137))
}

func (w *WSuite) TestConnectionStream(c *C) {
	conn := &fakeConn{
		ReadBuffer: messages(
			&protocol.StreamResponse{Name: proto.String("stdout"), Data: proto.String("1")},
			&protocol.StreamResponse{Name: proto.String("stderr"), Data: proto.String("2")},
			&protocol.StreamResponse{ExitStatus: proto.Uint32(3)},
		),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	resp, err := connection.Stream("foo-handle", 42)
	c.Assert(err, IsNil)

	c.Assert(
		string(conn.WriteBuffer.Bytes()),
		Equals,
		string(messages(&protocol.StreamRequest{
			Handle: proto.String("foo-handle"),
			JobId:  proto.Uint32(42),
		}).Bytes()),
	)

	res1 := <-resp
	c.Assert(res1.GetName(), Equals, "stdout")
	c.Assert(res1.GetData(), Equals, "1")

	res2 := <-resp
	c.Assert(res2.GetName(), Equals, "stderr")
	c.Assert(res2.GetData(), Equals, "2")

	res3, ok := <-resp
	c.Assert(res3.GetExitStatus(), Equals, uint32(3))
	c.Assert(ok, Equals, true)
}

func (w *WSuite) TestConnectionError(c *C) {
	conn := &fakeConn{
		ReadBuffer:  messages(&protocol.ErrorResponse{Message: proto.String("boo")}),
		WriteBuffer: bytes.NewBuffer([]byte{}),
	}

	connection := NewConnection(conn)

	resp, err := connection.Run("foo-handle", "echo hi")
	c.Assert(resp, IsNil)
	c.Assert(err, Not(IsNil))

	c.Assert(err.Error(), Equals, "boo")
}
