package oneme

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

const endpoint = "wss://ws-api.oneme.ru/websocket"
const origin = "https://web.max.ru"

type RawApiClient struct {
	conn              *websocket.Conn
	undispatchedQueue []map[string]any
}

func NewRawApiClient() (*RawApiClient, error) {
	header := http.Header{}
	header.Set("Origin", origin)

	conn, _, err := websocket.DefaultDialer.Dial(endpoint, header)
	if err != nil {
		return nil, err
	}
	return &RawApiClient{
		conn:              conn,
		undispatchedQueue: make([]map[string]any, 0),
	}, nil
}

func SendRawMessage[T any](c *RawApiClient, message Message[T]) error {
	data, err := json.Marshal(message)
	fmt.Printf("[OneMe Client] %s\n", data)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *RawApiClient) ReceiveRawMessage(seq *int) (map[string]any, error) {
	if seq == nil && len(c.undispatchedQueue) > 0 {
		msg := c.undispatchedQueue[0]
		c.undispatchedQueue = c.undispatchedQueue[1:]
		return msg, nil
	}

	for {
		_, data, err := c.conn.ReadMessage()
		fmt.Printf("[OneMe Server] %s\n", data)
		if err != nil {
			return nil, err
		}

		var message map[string]any
		if err := json.Unmarshal(data, &message); err != nil {
			return nil, err
		}

		if seqVal, ok := message["seq"].(float64); ok && (seq == nil || int(seqVal) == *seq) {
			return message, nil
		}
		c.undispatchedQueue = append(c.undispatchedQueue, message)
	}
}

func (c *RawApiClient) Close() error {
	return c.conn.Close()
}
