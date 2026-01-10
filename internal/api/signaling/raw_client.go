package signaling

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	callsMessages "github.com/icyfalc0n/maxcalls/internal/api/calls/messages"
	onemeMessages "github.com/icyfalc0n/maxcalls/internal/api/oneme/messages"
)

const origin = "https://web.max.ru"

type OutgoingMessage struct {
	Bytes   []byte
	ErrChan chan<- error
}

type IncomingMessage struct {
	Bytes []byte
	Err   error
}

type ActorCloseRequest struct {
	ErrChan chan<- error
}

type RawSignalingClient struct {
	incomingMessages <-chan IncomingMessage
	outgoingMessages chan<- OutgoingMessage
	closeChan        chan<- ActorCloseRequest
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

func newFromEndpoint(endpoint string) (RawSignalingClient, error) {
	header := http.Header{}
	header.Set("Origin", origin)

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

	incomingMessages := make(chan IncomingMessage, 10)
	outgoingMessages := make(chan OutgoingMessage, 10)
	closeChan := make(chan ActorCloseRequest)

	actor := RawClientActor{conn, incomingMessages, outgoingMessages, closeChan}
	go actor.Start()

	return RawSignalingClient{incomingMessages, outgoingMessages, closeChan}, nil
}

func (c *RawSignalingClient) Send(message []byte) error {
	errChan := make(chan error)
	defer close(errChan)
	c.outgoingMessages <- OutgoingMessage{Bytes: message, ErrChan: errChan}
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
	incomingMessage := <-c.incomingMessages
	return incomingMessage.Bytes, incomingMessage.Err
}

func (c *RawSignalingClient) ReceiveJSON(v any) error {
	msg, err := c.Receive()
	if err != nil {
		return err
	}
	return json.Unmarshal(msg, v)
}

func (c *RawSignalingClient) Close() error {
	errChan := make(chan error)
	defer close(errChan)
	defer close(c.closeChan)
	defer close(c.outgoingMessages)

	closeRequest := ActorCloseRequest{ErrChan: errChan}
	c.closeChan <- closeRequest
	return <-errChan
}
