package warden

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"sync"

	"code.google.com/p/goprotobuf/proto"

	protocol "github.com/vito/gordon/protocol"
)

type Connection struct {
	conn      net.Conn
	read      *bufio.Reader
	writeLock sync.Mutex
	readLock  sync.Mutex

	disconnected chan bool
}

type WardenError struct {
	Message   string
	Data      string
	Backtrace []string
}

func (e *WardenError) Error() string {
	return e.Message
}

func Connect(socket_path string) (*Connection, error) {
	conn, err := net.Dial("unix", socket_path)
	if err != nil {
		return nil, err
	}

	return NewConnection(conn), nil
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
		read: bufio.NewReader(conn),

		// buffer size of 1 so that read and write errors
		// can both send without blocking
		disconnected: make(chan bool, 1),
	}
}

func (c *Connection) Close() {
	c.conn.Close()
}

func (c *Connection) Create() (*protocol.CreateResponse, error) {
	res, err := c.roundTrip(&protocol.CreateRequest{}, &protocol.CreateResponse{})
	if err != nil {
		return nil, err
	}

	return res.(*protocol.CreateResponse), nil
}

func (c *Connection) Destroy(handle string) (*protocol.DestroyResponse, error) {
	res, err := c.roundTrip(
		&protocol.DestroyRequest{Handle: proto.String(handle)},
		&protocol.DestroyResponse{},
	)

	if err != nil {
		return nil, err
	}

	return res.(*protocol.DestroyResponse), nil
}

func (c *Connection) Spawn(handle, script string) (*protocol.SpawnResponse, error) {
	res, err := c.roundTrip(
		&protocol.SpawnRequest{
			Handle: proto.String(handle),
			Script: proto.String(script),
		},
		&protocol.SpawnResponse{},
	)

	if err != nil {
		return nil, err
	}

	return res.(*protocol.SpawnResponse), nil
}

func (c *Connection) Run(handle, script string) (*protocol.RunResponse, error) {
	res, err := c.roundTrip(
		&protocol.RunRequest{
			Handle: proto.String(handle),
			Script: proto.String(script),
		},
		&protocol.RunResponse{},
	)

	if err != nil {
		return nil, err
	}

	return res.(*protocol.RunResponse), nil
}

func (c *Connection) Stream(handle string, jobId uint32) (chan *protocol.StreamResponse, error) {
	err := c.sendMessage(
		&protocol.StreamRequest{
			Handle: proto.String(handle),
			JobId:  proto.Uint32(jobId),
		},
	)

	if err != nil {
		return nil, err
	}

	responses := make(chan *protocol.StreamResponse)

	go func() {
		for {
			resMsg, err := c.readResponse(&protocol.StreamResponse{})
			if err != nil {
				close(responses)
				break
			}

			response := resMsg.(*protocol.StreamResponse)

			responses <- response

			if response.ExitStatus != nil {
				close(responses)
				break
			}
		}
	}()

	return responses, nil
}

func (c *Connection) NetIn(handle string) (*protocol.NetInResponse, error) {
	res, err := c.roundTrip(
		&protocol.NetInRequest{Handle: proto.String(handle)},
		&protocol.NetInResponse{},
	)

	if err != nil {
		return nil, err
	}

	return res.(*protocol.NetInResponse), nil
}

func (c *Connection) LimitMemory(handle string, limit uint64) (*protocol.LimitMemoryResponse, error) {
	res, err := c.roundTrip(
		&protocol.LimitMemoryRequest{
			Handle:       proto.String(handle),
			LimitInBytes: proto.Uint64(limit),
		},
		&protocol.LimitMemoryResponse{},
	)

	if err != nil {
		return nil, err
	}

	return res.(*protocol.LimitMemoryResponse), nil
}

func (c *Connection) GetMemoryLimit(handle string) (uint64, error) {
	res, err := c.roundTrip(
		&protocol.LimitMemoryRequest{
			Handle: proto.String(handle),
		},
		&protocol.LimitMemoryResponse{},
	)

	if err != nil {
		return 0, err
	}

	limit := res.(*protocol.LimitMemoryResponse).GetLimitInBytes()
	if limit == math.MaxInt64 { // PROBABLY NOT A LIMIT
		return 0, nil
	}

	return limit, nil
}

func (c *Connection) LimitDisk(handle string, limit uint64) (*protocol.LimitDiskResponse, error) {
	res, err := c.roundTrip(
		&protocol.LimitDiskRequest{
			Handle:    proto.String(handle),
			ByteLimit: proto.Uint64(limit),
		},
		&protocol.LimitDiskResponse{},
	)

	if err != nil {
		return nil, err
	}

	return res.(*protocol.LimitDiskResponse), nil
}

func (c *Connection) GetDiskLimit(handle string) (uint64, error) {
	res, err := c.roundTrip(
		&protocol.LimitDiskRequest{
			Handle: proto.String(handle),
		},
		&protocol.LimitDiskResponse{},
	)

	if err != nil {
		return 0, err
	}

	return res.(*protocol.LimitDiskResponse).GetByteLimit(), nil
}

func (c *Connection) CopyIn(handle, src, dst string) (*protocol.CopyInResponse, error) {
	res, err := c.roundTrip(
		&protocol.CopyInRequest{
			Handle:  proto.String(handle),
			SrcPath: proto.String(src),
			DstPath: proto.String(dst),
		},
		&protocol.CopyInResponse{},
	)

	if err != nil {
		return nil, err
	}

	return res.(*protocol.CopyInResponse), nil
}

func (c *Connection) List() (*protocol.ListResponse, error) {
	res, err := c.roundTrip(&protocol.ListRequest{}, &protocol.ListResponse{})
	if err != nil {
		return nil, err
	}

	return res.(*protocol.ListResponse), nil
}

func (c *Connection) Info(handle string) (*protocol.InfoResponse, error) {
	res, err := c.roundTrip(
		&protocol.InfoRequest{
			Handle: proto.String(handle),
		},
		&protocol.InfoResponse{},
	)
	if err != nil {
		return nil, err
	}

	return res.(*protocol.InfoResponse), nil
}

func (c *Connection) roundTrip(request proto.Message, response proto.Message) (proto.Message, error) {
	err := c.sendMessage(request)
	if err != nil {
		return nil, err
	}

	resp, err := c.readResponse(response)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Connection) sendMessage(req proto.Message) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	request, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	msg := &protocol.Message{
		Type:    protocol.Message_Type(message2type(req)).Enum(),
		Payload: request,
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(
		[]byte(
			fmt.Sprintf(
				"%d\r\n%s\r\n",
				len(data),
				data,
			),
		),
	)

	if err != nil {
		c.disconnected <- true
		return err
	}

	return nil
}

func (c *Connection) readResponse(response proto.Message) (proto.Message, error) {
	payload, err := c.readPayload()
	if err != nil {
		c.disconnected <- true
		return nil, err
	}

	message := &protocol.Message{}
	err = proto.Unmarshal(payload, message)
	if err != nil {
		return nil, err
	}

	// error response from server
	if message.GetType() == protocol.Message_Type(1) {
		errorResponse := &protocol.ErrorResponse{}
		err = proto.Unmarshal(message.Payload, errorResponse)
		if err != nil {
			return nil, errors.New("error unmarshalling error!")
		}

		return nil, &WardenError{
			Message:   errorResponse.GetMessage(),
			Data:      errorResponse.GetData(),
			Backtrace: errorResponse.GetBacktrace(),
		}
	}

	response_type := protocol.Message_Type(message2type(response))
	if message.GetType() != response_type {
		return nil, errors.New(
			fmt.Sprintf(
				"expected message type %s, got %s\n",
				response_type.String(),
				message.GetType().String(),
			),
		)
	}

	err = proto.Unmarshal(message.GetPayload(), response)
	return response, err
}

func (c *Connection) readPayload() ([]byte, error) {
	c.readLock.Lock()
	defer c.readLock.Unlock()

	msgHeader, err := c.read.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	msgLen, err := strconv.ParseUint(string(msgHeader[0:len(msgHeader)-2]), 10, 0)
	if err != nil {
		return nil, err
	}

	payload, err := readNBytes(int(msgLen), c.read)
	if err != nil {
		return nil, err
	}

	_, err = readNBytes(2, c.read) // CRLN
	if err != nil {
		return nil, err
	}

	return payload, err
}

func message2type(msg proto.Message) int32 {
	switch msg.(type) {
	case *protocol.ErrorResponse:
		return 1

	case *protocol.CreateRequest, *protocol.CreateResponse:
		return 11
	case *protocol.StopRequest, *protocol.StopResponse:
		return 12
	case *protocol.DestroyRequest, *protocol.DestroyResponse:
		return 13
	case *protocol.InfoRequest, *protocol.InfoResponse:
		return 14

	case *protocol.SpawnRequest, *protocol.SpawnResponse:
		return 21
	case *protocol.LinkRequest, *protocol.LinkResponse:
		return 22
	case *protocol.RunRequest, *protocol.RunResponse:
		return 23
	case *protocol.StreamRequest, *protocol.StreamResponse:
		return 24

	case *protocol.NetInRequest, *protocol.NetInResponse:
		return 31
	case *protocol.NetOutRequest, *protocol.NetOutResponse:
		return 32

	case *protocol.CopyInRequest, *protocol.CopyInResponse:
		return 41
	case *protocol.CopyOutRequest, *protocol.CopyOutResponse:
		return 42

	case *protocol.LimitMemoryRequest, *protocol.LimitMemoryResponse:
		return 51
	case *protocol.LimitDiskRequest, *protocol.LimitDiskResponse:
		return 52
	case *protocol.LimitBandwidthRequest, *protocol.LimitBandwidthResponse:
		return 53

	case *protocol.PingRequest, *protocol.PingResponse:
		return 91
	case *protocol.ListRequest, *protocol.ListResponse:
		return 92
	case *protocol.EchoRequest, *protocol.EchoResponse:
		return 93
	}

	panic("unknown message type")
}

func readNBytes(payloadLen int, io *bufio.Reader) ([]byte, error) {
	payload := make([]byte, payloadLen)

	for readCount := 0; readCount < payloadLen; {
		n, err := io.Read(payload[readCount:])
		if err != nil {
			return nil, err
		}

		readCount += n
	}

	return payload, nil
}
