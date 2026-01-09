package signaling

import (
	"strconv"

	callsMessages "github.com/icyfalc0n/max_calls_api/api/calls/messages"
	onemeMessages "github.com/icyfalc0n/max_calls_api/api/oneme/messages"
	"github.com/icyfalc0n/max_calls_api/api/signaling/messages"
)

type SignalingClient struct {
	rawClient     RawSignalingClient
	participantId int64
	sequence      *int
}

func NewSignalingFromIncoming(incomingCall onemeMessages.IncomingCall, loginData callsMessages.LoginData) (SignalingClient, error) {
	rawClient, err := NewRawSignalingFromIncoming(incomingCall, loginData)
	if err != nil {
		return SignalingClient{}, err
	}

	callerExternalID := incomingCall.CallerID
	callerID, err := readUserID(rawClient, strconv.Itoa(callerExternalID))
	if err != nil {
		return SignalingClient{}, err
	}

	initialSequence := 1
	client := SignalingClient{rawClient, callerID, &initialSequence}
	err = client.sendMessage(messages.NewAcceptCall(*client.sequence))
	if err != nil {
		return SignalingClient{}, err
	}

	return client, nil
}

func NewSignalingFromOutgoing(signalingServerEndpoint string, calltakerExternalID string) (SignalingClient, error) {
	rawClient, err := NewRawSignalingFromOutgoing(signalingServerEndpoint)
	if err != nil {
		return SignalingClient{}, err
	}

	calltakerID, err := readUserID(rawClient, calltakerExternalID)
	if err != nil {
		return SignalingClient{}, err
	}

	initialSequence := 1
	return SignalingClient{rawClient, calltakerID, &initialSequence}, nil
}

func readUserID(rawClient RawSignalingClient, externalUserID string) (int64, error) {
	var serverHello messages.ServerHello
	err := rawClient.ReceiveJSON(&serverHello)
	if err != nil {
		return 0, err
	}

	return messages.FindUserIDByExternalID(serverHello, externalUserID), nil
}

func (c *SignalingClient) sendMessage(v any) error {
	err := c.rawClient.SendJSON(v)
	if err != nil {
		return err
	}
	*c.sequence += 1

	return nil
}

func (c *SignalingClient) ReceiveSignal() (any, error) {
	for {
		// It can be valid JSON but with type and notification fields omitted, so we're doing deserialization manually
		var msg map[string]any
		err := c.rawClient.ReceiveJSON(&msg)
		if err != nil {
			return nil, err
		}

		if msg["type"].(string) != "notification" {
			continue
		}
		if msg["notification"].(string) != "transmitted-data" {
			continue
		}

		return msg["data"], nil
	}
}

func (c *SignalingClient) SendSignal(signal any) error {
	msg := messages.NewTransmitData(*c.sequence, c.participantId, signal)
	return c.sendMessage(msg)
}
