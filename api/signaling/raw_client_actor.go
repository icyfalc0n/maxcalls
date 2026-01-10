package signaling

import (
	"slices"

	"github.com/gorilla/websocket"
)

type RawClientActor struct {
	conn             *websocket.Conn
	incomingMessages chan<- IncomingMessage
	outgoingMessages <-chan OutgoingMessage
}

func (actor *RawClientActor) Start() {
	for {
		select {
		case outgoingMessage := <-actor.outgoingMessages:
			err := actor.conn.WriteMessage(websocket.TextMessage, outgoingMessage.Bytes)
			outgoingMessage.ErrChan <- err
			if err != nil {
				continue
			}
		default:
			_, msg, err := actor.conn.ReadMessage()
			if err != nil {
				actor.incomingMessages <- IncomingMessage{Err: err}
				continue
			}

			isPing, err := answerPing(actor.conn, msg)
			if err != nil {
				actor.incomingMessages <- IncomingMessage{Err: err}
				continue
			}
			if isPing {
				continue
			}

			actor.incomingMessages <- IncomingMessage{Bytes: msg}
		}
	}
}

func answerPing(conn *websocket.Conn, msg []byte) (bool, error) {
	if slices.Equal(msg, []byte("ping")) {
		return true, conn.WriteMessage(websocket.TextMessage, []byte("pong"))
	}

	return false, nil
}
