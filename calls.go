package maxcalls

import (
	"github.com/icyfalc0n/maxcalls/internal/api/calls"
	callsMessages "github.com/icyfalc0n/maxcalls/internal/api/calls/messages"
	"github.com/icyfalc0n/maxcalls/internal/api/oneme"
	"github.com/icyfalc0n/maxcalls/internal/api/signaling"
	"github.com/icyfalc0n/maxcalls/internal/ice"
)

type Calls struct {
	onemeClient oneme.ApiClient
	callsClient calls.ApiClient
	loginData   callsMessages.LoginData
}

func NewCalls(authToken string) (Calls, error) {
	onemeClient, err := oneme.NewOnemeApiClient()
	if err != nil {
		return Calls{}, err
	}

	_, err = onemeClient.DoChatSync(authToken)
	if err != nil {
		onemeClient.Close()
		return Calls{}, err
	}

	callsLoginToken, err := onemeClient.DoCallTokenRequest()
	if err != nil {
		onemeClient.Close()
		return Calls{}, err
	}

	callsClient := calls.ApiClient{}
	loginData, err := callsClient.Login(callsLoginToken)
	if err != nil {
		onemeClient.Close()
		return Calls{}, err
	}

	return Calls{
		onemeClient,
		callsClient,
		loginData,
	}, nil
}

func (c *Calls) Call(id string) (Connection, error) {
	startedConversationInfo, err := c.callsClient.StartConversation(
		c.loginData.SessionKey, id)
	if err != nil {
		return Connection{}, err
	}

	signalingClient, err := signaling.NewSignalingFromOutgoing(startedConversationInfo.Endpoint, id)
	if err != nil {
		return Connection{}, err
	}

	iceAgent, err := ice.NewAgentFromOutgoing(startedConversationInfo)
	if err != nil {
		signalingClient.Close()
		return Connection{}, err
	}

	connector := ice.IceConnector{
		SignalingClient: signalingClient,
		IceAgent:        iceAgent,
	}
	conn, err := connector.Connect()
	if err != nil {
		connector.Close()
		return Connection{}, err
	}

	return Connection{conn, connector}, nil
}

func (c *Calls) WaitForCall() (Connection, error) {
	incomingCall, err := c.onemeClient.WaitForIncomingCall()
	if err != nil {
		return Connection{}, err
	}

	signalingClient, err := signaling.NewSignalingFromIncoming(incomingCall, c.loginData)
	if err != nil {
		return Connection{}, err
	}

	iceAgent, err := ice.NewAgentFromIncoming(incomingCall)
	if err != nil {
		signalingClient.Close()
		return Connection{}, err
	}

	connector := ice.IceConnector{SignalingClient: signalingClient, IceAgent: iceAgent}
	conn, err := connector.Connect()
	if err != nil {
		connector.Close()
		return Connection{}, err
	}

	return Connection{conn, connector}, nil
}

func (c *Calls) ExternalUserId() string {
	return c.loginData.ExternalUserID
}

func (c *Calls) Close() error {
	return c.onemeClient.Close()
}
