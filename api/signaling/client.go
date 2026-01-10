package signaling

import (
	"encoding/json"
	"strconv"
	"sync/atomic"

	callsMessages "github.com/icyfalc0n/max_calls_api/api/calls/messages"
	onemeMessages "github.com/icyfalc0n/max_calls_api/api/oneme/messages"
	"github.com/icyfalc0n/max_calls_api/api/signaling/messages"
)

type SignalingClient struct {
	rawClient     RawSignalingClient
	participantId int64
	sequence      *atomic.Int32
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

	initialSequence := atomic.Int32{}
	client := SignalingClient{rawClient, callerID, &initialSequence}
	err = client.rawClient.SendJSON(messages.NewAcceptCall(client.nextSequence()))
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

	initialSequence := atomic.Int32{}
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

func (c *SignalingClient) nextSequence() int {
	return int(c.sequence.Add(1))
}

func (c *SignalingClient) ReceiveSignal(v any) error {
	for {
		// It can be valid JSON but with type and notification fields omitted, so we're doing deserialization manually
		var msg map[string]any
		err := c.rawClient.ReceiveJSON(&msg)
		if err != nil {
			return err
		}

		if msg["type"].(string) != "notification" {
			continue
		}
		if msg["notification"].(string) != "transmitted-data" {
			continue
		}

		marshaled, err := json.Marshal(msg["data"])
		if err != nil {
			return err
		}
		err = json.Unmarshal(marshaled, v)
		if err != nil {
			return err
		}
		return nil
	}
}

func (c *SignalingClient) SendSignal(signal any) error {
	msg := messages.NewTransmitData(c.nextSequence(), c.participantId, signal)
	return c.rawClient.SendJSON(msg)
}
