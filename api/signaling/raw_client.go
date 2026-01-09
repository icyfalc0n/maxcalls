package signaling

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/gorilla/websocket"
	callsMessages "github.com/icyfalc0n/max_calls_api/api/calls/messages"
	onemeMessages "github.com/icyfalc0n/max_calls_api/api/oneme/messages"
)

type OutgoingMessage struct {
	Bytes   []byte
	ErrChan chan<- error
}

type IncomingMessage struct {
	Bytes []byte
	Err   error
}

type RawSignalingClient struct {
	receiveChannel <-chan IncomingMessage
	sendChannel    chan<- OutgoingMessage
}

func (c *RawSignalingClient) Send(message []byte) error {
	errChan := make(chan error)
	c.sendChannel <- OutgoingMessage{Bytes: message, ErrChan: errChan}
	return <-errChan
}

func (c *RawSignalingClient) SendJSON(v any) error {
	marshaled, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.Send(marshaled)
}

func (c *RawSignalingClient) Receive() ([]byte, error) {
	incomingMessage := <-c.receiveChannel
	return incomingMessage.Bytes, incomingMessage.Err
}

func (c *RawSignalingClient) ReceiveJSON(v any) error {
	msg, err := c.Receive()
	if err != nil {
		return err
	}
	return json.Unmarshal(msg, v)
}

func newFromEndpoint(endpoint string) (RawSignalingClient, error) {
	header := http.Header{}
	header.Set("Origin", "https://web.max.ru")

	query := url.Values{}
	query.Set("platform", "WEB")
	query.Set("appVersion", "1.1")
	query.Set("version", "5")
	query.Set("device", "browser")
	query.Set("capabilities", "603F")
	query.Set("clientType", "ONE_ME")
	query.Set("tgt", "start")

	endpoint = fmt.Sprintf("%s&%s", endpoint, query.Encode())
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, header)
	if err != nil {
		return RawSignalingClient{}, err
	}

	receiveChannel := make(chan IncomingMessage, 10)
	sendChannel := make(chan OutgoingMessage, 10)

	go startRawClientActor(conn, receiveChannel, sendChannel)

	return RawSignalingClient{receiveChannel, sendChannel}, nil
}

func startRawClientActor(conn *websocket.Conn, receiveChannel chan<- IncomingMessage, sendChannel <-chan OutgoingMessage) {
	for {
		select {
		case outgoingMessage := <-sendChannel:
			err := conn.WriteMessage(websocket.TextMessage, outgoingMessage.Bytes)
			outgoingMessage.ErrChan <- err
			if err != nil {
				continue
			}
			fmt.Printf("[Signaling.Client] %s\n", outgoingMessage.Bytes)
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				receiveChannel <- IncomingMessage{Err: err}
				continue
			}

			isPing, err := answerPing(conn, msg)
			if err != nil {
				receiveChannel <- IncomingMessage{Err: err}
				continue
			}
			if isPing {
				continue
			}

			fmt.Printf("[Signaling.Server] %s\n", msg)
			receiveChannel <- IncomingMessage{Bytes: msg}
		}
	}
}

func answerPing(conn *websocket.Conn, msg []byte) (bool, error) {
	if slices.Equal(msg, []byte("ping")) {
		return true, conn.WriteMessage(websocket.TextMessage, []byte("pong"))
	}

	return false, nil
}

func NewRawSignalingFromIncoming(incomingCall onemeMessages.IncomingCall, loginData callsMessages.LoginData) (RawSignalingClient, error) {
	query := url.Values{}
	query.Set("userId", loginData.UID)
	query.Set("entityType", "USER")
	query.Set("conversationId", incomingCall.ConversationID)
	query.Set("token", incomingCall.Signaling.Token)
	endpoint := fmt.Sprintf("%s?%s", incomingCall.Signaling.URL, query.Encode())

	return newFromEndpoint(endpoint)
}

func NewRawSignalingFromOutgoing(signalingServerEndpoint string) (RawSignalingClient, error) {
	return newFromEndpoint(signalingServerEndpoint)
}
