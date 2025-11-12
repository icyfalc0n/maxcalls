package api

import (
	"encoding/json"
	"slices"
	"strconv"

	calls_server "github.com/icyfalc0n/max_calls_api/api/raw/calls/messages/server"
	oneme_server "github.com/icyfalc0n/max_calls_api/api/raw/oneme/messages/server"
	raw_singaling "github.com/icyfalc0n/max_calls_api/api/raw/signaling"
	signalingClient "github.com/icyfalc0n/max_calls_api/api/raw/signaling/messages/client"
	signalingServer "github.com/icyfalc0n/max_calls_api/api/raw/signaling/messages/server"
)

type SignalingClient struct {
	SendChannel    chan<- any
	ReceiveChannel <-chan any
}

func NewSignalingFromIncoming(incomingCall oneme_server.IncomingCall, loginData calls_server.LoginData) (SignalingClient, error) {
	rawClient, err := raw_singaling.NewSignalingFromIncoming(incomingCall, loginData)
	if err != nil {
		return SignalingClient{}, err
	}

	callerExternalID := incomingCall.CallerID
	callerID, err := readUserID(rawClient, strconv.Itoa(callerExternalID))
	if err != nil {
		return SignalingClient{}, err
	}

	acceptCallMsg := signalingClient.NewAcceptCall(1)
	rawClient.WriteJSON(acceptCallMsg)

	sendChannel := make(chan any, 10)
	receiveChannel := make(chan any, 10)
	go startChannelConverter(rawClient, receiveChannel, sendChannel, callerID)

	return SignalingClient{sendChannel, receiveChannel}, nil
}

func NewSignalingFromOutgoing(startedConversationInfo calls_server.StartedConversationInfo, calltakerExternalID string) (SignalingClient, error) {
	rawClient, err := raw_singaling.NewSignalingFromOutgoing(startedConversationInfo)
	if err != nil {
		return SignalingClient{}, err
	}

	calltakerID, err := readUserID(rawClient, calltakerExternalID)
	if err != nil {
		return SignalingClient{}, err
	}

	sendChannel := make(chan any, 10)
	receiveChannel := make(chan any, 10)
	go startChannelConverter(rawClient, receiveChannel, sendChannel, calltakerID)

	return SignalingClient{sendChannel, receiveChannel}, nil
}

func readUserID(rawClient raw_singaling.ApiClient, externalUserID string) (int64, error) {
	msg := rawClient.Read()

	var serverHello signalingServer.ServerHello
	err := json.Unmarshal(msg, &serverHello)
	if err != nil {
		return 0, err
	}

	return signalingServer.FindUserIDByExternalID(serverHello, externalUserID), nil
}

func startChannelConverter(rawClient raw_singaling.ApiClient, receiveChannel chan<- any, sendChannel <-chan any, callerID int64) {
	sequence := 2
	for {
		select {
		case receivedMessage := <-rawClient.ReceiveChannel:
			if slices.Equal(receivedMessage, []byte("ping")) {
				rawClient.Write([]byte("pong"))
				continue
			}

			var decodedMessage map[string]any
			err := json.Unmarshal(receivedMessage, &decodedMessage)
			if err != nil {
				panic(err)
			}

			if decodedMessage["type"].(string) != "notification" {
				continue
			}
			if decodedMessage["notification"].(string) != "transmitted-data" {
				continue
			}

			receiveChannel <- decodedMessage["data"]

		case messageData := <-sendChannel:
			message := signalingClient.NewTransmitData(sequence, callerID, messageData)
			err := rawClient.WriteJSON(message)
			if err != nil {
				panic(err)
			}
			sequence += 1

		}
	}
}
