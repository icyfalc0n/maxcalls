package maxcalls

import (
	"github.com/icyfalc0n/maxcalls/internal/ice"
	pionIce "github.com/pion/ice/v4"
)

type Connection struct {
	conn      *pionIce.Conn
	connector ice.IceConnector
}

func (c *Connection) Close() error {
	return c.connector.Close()
}
