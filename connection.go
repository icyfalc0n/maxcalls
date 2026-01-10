package maxcalls

import (
	"github.com/icyfalc0n/maxcalls/internal/ice"
	pionIce "github.com/pion/ice/v4"
)

// Connection represents an established call connection with a remote peer.
// It wraps a Pion ICE connection and provides methods to manage the connection lifecycle.
type Connection struct {
	conn      *pionIce.Conn
	connector ice.IceConnector
}

// Close closes the connection and releases all associated resources.
// This includes closing the ICE connection and the signaling client.
// After calling Close, the Connection should not be used.
//
// Returns an error if closing the connection fails.
func (c *Connection) Close() error {
	return c.connector.Close()
}
