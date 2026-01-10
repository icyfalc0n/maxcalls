package max_calls_api

import (
	"github.com/icyfalc0n/max_calls_api/internal/ice"
	pionIce "github.com/pion/ice/v4"
)

type Connection struct {
	conn      *pionIce.Conn
	connector ice.IceConnector
}

func (c *Connection) Close() error {
	return c.connector.Close()
}
