package signaling

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"

	calls_server "github.com/icyfalc0n/max_calls_api/api/raw/calls/messages/server"

	oneme_server "github.com/icyfalc0n/max_calls_api/api/raw/oneme/messages/server"
)

type ApiClient struct {
	ReceiveChannel <-chan []byte
	SendChannel    chan<- []byte
}

func (c *ApiClient) Write(message []byte) {
	c.SendChannel <- message
}

func (c *ApiClient) WriteJSON(v any) error {
	marshaled, err := json.Marshal(v)
	if err != nil {
		return err
	}

	c.Write(marshaled)
	return nil
}

func (c *ApiClient) Read() []byte {
	return <-c.ReceiveChannel
}

func newFromEndpoint(endpoint string) (ApiClient, error) {
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
		return ApiClient{}, err
	}

	receiveChannel := make(chan []byte, 10)
	sendChannel := make(chan []byte, 10)

	go startChannelConverter(conn, receiveChannel, sendChannel)

	return ApiClient{receiveChannel, sendChannel}, nil
}

func startChannelConverter(conn *websocket.Conn, receiveChannel chan<- []byte, sendChannel <-chan []byte) {
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

func NewSignalingFromIncoming(incomingCall oneme_server.IncomingCall, loginData calls_server.LoginData) (ApiClient, error) {
	query := url.Values{}
	query.Set("userId", loginData.UID)
	query.Set("entityType", "USER")
	query.Set("conversationId", incomingCall.ConversationID)
	query.Set("token", incomingCall.Signaling.Token)
	endpoint := fmt.Sprintf("%s?%s", incomingCall.Signaling.URL, query.Encode())

	return newFromEndpoint(endpoint)
}

func NewSignalingFromOutgoing(startedConversationInfo calls_server.StartedConversationInfo) (ApiClient, error) {
	return newFromEndpoint(startedConversationInfo.Endpoint)
}
