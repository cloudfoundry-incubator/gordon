package warden

import (
	"time"

	protocol "github.com/vito/gordon/protocol"
)

type Client interface {
	List() (*protocol.ListResponse, error)

	Create() (*protocol.CreateResponse, error)
	Destroy(handle string) (*protocol.DestroyResponse, error)

	Spawn(handle, script string) (*protocol.SpawnResponse, error)
	Link(handle string, jobId uint32) (*protocol.LinkResponse, error)
}

type WardenClient struct {
	SocketPath string

	connectionProvider ConnectionProvider
	connection         chan *Connection
}

func NewClient(cp ConnectionProvider) *WardenClient {
	return &WardenClient{
		connectionProvider: cp,
		connection:         make(chan *Connection),
	}
}

func (c *WardenClient) Connect() error {
	conn, err := c.connectionProvider.ProvideConnection()
	if err != nil {
		return err
	}

	go c.serveConnections(conn)

	return nil
}

func (c *WardenClient) Create() (*protocol.CreateResponse, error) {
	return (<-c.connection).Create()
}

func (c *WardenClient) Destroy(handle string) (*protocol.DestroyResponse, error) {
	return (<-c.connection).Destroy(handle)
}

func (c *WardenClient) Spawn(handle, script string) (*protocol.SpawnResponse, error) {
	return (<-c.connection).Spawn(handle, script)
}

func (c *WardenClient) NetIn(handle string) (*protocol.NetInResponse, error) {
	return (<-c.connection).NetIn(handle)
}

func (c *WardenClient) LimitMemory(handle string, limit uint64) (*protocol.LimitMemoryResponse, error) {
	return (<-c.connection).LimitMemory(handle, limit)
}

func (c *WardenClient) GetMemoryLimit(handle string) (uint64, error) {
	return (<-c.connection).GetMemoryLimit(handle)
}

func (c *WardenClient) LimitDisk(handle string, limit uint64) (*protocol.LimitDiskResponse, error) {
	return (<-c.connection).LimitDisk(handle, limit)
}

func (c *WardenClient) GetDiskLimit(handle string) (uint64, error) {
	return (<-c.connection).GetDiskLimit(handle)
}

func (c *WardenClient) List() (*protocol.ListResponse, error) {
	return (<-c.connection).List()
}

func (c *WardenClient) Link(handle string, jobId uint32) (*protocol.LinkResponse, error) {
	return (<-c.connection).Link(handle, jobId)
}

func (c *WardenClient) Info(handle string) (*protocol.InfoResponse, error) {
	return (<-c.connection).Info(handle)
}

func (c *WardenClient) CopyIn(handle, src, dst string) (*protocol.CopyInResponse, error) {
	return c.acquireConnection().CopyIn(handle, src, dst)
}

func (c *WardenClient) Stream(handle string, jobId uint32) (chan *protocol.StreamResponse, error) {
	return c.acquireConnection().Stream(handle, jobId)
}

func (c *WardenClient) Run(handle, script string) (*protocol.RunResponse, error) {
	return c.acquireConnection().Run(handle, script)
}

func (c *WardenClient) serveConnections(conn *Connection) {
	for stop := false; !stop; {
		select {
		case <-conn.disconnected:
			stop = true
			break

		case c.connection <- conn:
		}
	}

	go c.serveConnections(c.acquireConnection())
}

func (c *WardenClient) acquireConnection() *Connection {
	for {
		conn, err := c.connectionProvider.ProvideConnection()
		if err == nil {
			return conn
		}

		time.Sleep(500 * time.Millisecond)
	}
}
