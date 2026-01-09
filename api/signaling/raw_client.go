package signaling

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	calls_messages "github.com/icyfalc0n/max_calls_api/api/calls/messages"
	oneme_messages "github.com/icyfalc0n/max_calls_api/api/oneme/messages"
)

type RawApiClient struct {
	ReceiveChannel <-chan []byte
	SendChannel    chan<- []byte
}

func (c *RawApiClient) Write(message []byte) {
	c.SendChannel <- message
}

func (c *RawApiClient) WriteJSON(v any) error {
	marshaled, err := json.Marshal(v)
	if err != nil {
		return err
	}

	c.Write(marshaled)
	return nil
}

func (c *RawApiClient) Read() []byte {
	return <-c.ReceiveChannel
}

func newFromEndpoint(endpoint string) (RawApiClient, error) {
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
		return RawApiClient{}, err
	}

	receiveChannel := make(chan []byte, 10)
	sendChannel := make(chan []byte, 10)

	go startRawChannelConverter(conn, receiveChannel, sendChannel)

	return RawApiClient{receiveChannel, sendChannel}, nil
}

func startRawChannelConverter(conn *websocket.Conn, receiveChannel chan<- []byte, sendChannel <-chan []byte) {
	for {
		select {
		case messageToSend := <-sendChannel:
			err := conn.WriteMessage(websocket.TextMessage, messageToSend)
			if err != nil {
				panic(err)
			}
			fmt.Printf("[Signaling.Client] %s\n", messageToSend)
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				panic(err)
			}

			fmt.Printf("[Signaling.Server] %s\n", msg)
			receiveChannel <- msg
		}
	}
}

func NewRawSignalingFromIncoming(incomingCall oneme_messages.IncomingCall, loginData calls_messages.LoginData) (RawApiClient, error) {
	query := url.Values{}
	query.Set("userId", loginData.UID)
	query.Set("entityType", "USER")
	query.Set("conversationId", incomingCall.ConversationID)
	query.Set("token", incomingCall.Signaling.Token)
	endpoint := fmt.Sprintf("%s?%s", incomingCall.Signaling.URL, query.Encode())

	return newFromEndpoint(endpoint)
}

func NewRawSignalingFromOutgoing(startedConversationInfo calls_messages.StartedConversationInfo) (RawApiClient, error) {
	return newFromEndpoint(startedConversationInfo.Endpoint)
}
