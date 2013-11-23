package warden

import (
	"time"

	protocol "github.com/vito/gordon/protocol"
)

type Client struct {
	SocketPath string

	connectionProvider ConnectionProvider
	connection         chan *Connection
}

func NewClient(cp ConnectionProvider) *Client {
	return &Client{
		connectionProvider: cp,
		connection:         make(chan *Connection),
	}
}

func (c *Client) Connect() error {
	conn, err := c.connectionProvider.ProvideConnection()
	if err != nil {
		return err
	}

	go c.serveConnections(conn)

	return nil
}

func (c *Client) Create() (*protocol.CreateResponse, error) {
	return (<-c.connection).Create()
}

func (c *Client) Destroy(handle string) (*protocol.DestroyResponse, error) {
	return (<-c.connection).Destroy(handle)
}

func (c *Client) Spawn(handle, script string) (*protocol.SpawnResponse, error) {
	return (<-c.connection).Spawn(handle, script)
}

func (c *Client) NetIn(handle string) (*protocol.NetInResponse, error) {
	return (<-c.connection).NetIn(handle)
}

func (c *Client) LimitMemory(handle string, limit uint64) (*protocol.LimitMemoryResponse, error) {
	return (<-c.connection).LimitMemory(handle, limit)
}

func (c *Client) GetMemoryLimit(handle string) (uint64, error) {
	return (<-c.connection).GetMemoryLimit(handle)
}

func (c *Client) LimitDisk(handle string, limit uint64) (*protocol.LimitDiskResponse, error) {
	return (<-c.connection).LimitDisk(handle, limit)
}

func (c *Client) GetDiskLimit(handle string) (uint64, error) {
	return (<-c.connection).GetDiskLimit(handle)
}

func (c *Client) List() (*protocol.ListResponse, error) {
	return (<-c.connection).List()
}

func (c *Client) Info(handle string) (*protocol.InfoResponse, error) {
	return (<-c.connection).Info(handle)
}

func (c *Client) CopyIn(handle, src, dst string) (*protocol.CopyInResponse, error) {
	return c.acquireConnection().CopyIn(handle, src, dst)
}

func (c *Client) Stream(handle string, jobId uint32) (chan *protocol.StreamResponse, error) {
	return c.acquireConnection().Stream(handle, jobId)
}

func (c *Client) Run(handle, script string) (*protocol.RunResponse, error) {
	return c.acquireConnection().Run(handle, script)
}

func (c *Client) serveConnections(conn *Connection) {
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

func (c *Client) acquireConnection() *Connection {
	for {
		conn, err := c.connectionProvider.ProvideConnection()
		if err == nil {
			return conn
		}

		time.Sleep(500 * time.Millisecond)
	}
}
