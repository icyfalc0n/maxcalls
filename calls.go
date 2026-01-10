package maxcalls

import (
	"github.com/icyfalc0n/maxcalls/internal/api/calls"
	callsMessages "github.com/icyfalc0n/maxcalls/internal/api/calls/messages"
	"github.com/icyfalc0n/maxcalls/internal/api/oneme"
	"github.com/icyfalc0n/maxcalls/internal/api/signaling"
	"github.com/icyfalc0n/maxcalls/internal/ice"
	"github.com/icyfalc0n/maxcalls/logger"
)

// Calls represents a client for managing calls through the MAX messenger API.
// It maintains connections to the OneMe API for authentication and call notifications,
// and the Calls API for call session management.
type Calls struct {
	onemeClient oneme.ApiClient
	callsClient calls.ApiClient
	loginData   callsMessages.LoginData
	logger      logger.Logger
}

// NewCalls creates a new Calls client and authenticates with the MAX messenger API.
// It performs the following steps:
//  1. Establishes a connection to the OneMe API
//  2. Synchronizes chat state using the provided auth token
//  3. Requests a call token from OneMe
//  4. Logs into the Calls API using the call token
//
// The authToken parameter should be a valid authentication token for the MAX messenger.
//
// Returns a Calls instance ready to make or receive calls, or an error if authentication fails.
// The caller should call Close() when done to clean up resources.
func NewCalls(logger logger.Logger, authToken string) (Calls, error) {
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
		logger,
	}, nil
}

// Call initiates an outgoing call to the specified user ID.
// It performs the following steps:
//  1. Starts a conversation session with the Calls API
//  2. Establishes a signaling connection for the call
//  3. Creates an ICE agent for peer-to-peer connection
//  4. Connects using ICE to establish the media connection
//
// The id parameter should be the user ID of the person to call.
//
// Returns a Connection that can be used to send and receive media data,
// or an error if the call cannot be established. The caller should call
// Connection.Close() when done with the call.
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
		Logger:          c.logger,
	}
	conn, err := connector.Connect()
	if err != nil {
		connector.Close()
		return Connection{}, err
	}

	return Connection{conn, connector}, nil
}

// WaitForCall blocks until an incoming call is received, then accepts and establishes the connection.
// It performs the following steps:
//  1. Waits for an incoming call notification from the OneMe API
//  2. Establishes a signaling connection for the incoming call
//  3. Creates an ICE agent for peer-to-peer connection
//  4. Connects using ICE to establish the media connection
//
// Returns a Connection that can be used to send and receive media data,
// or an error if the call cannot be established. The caller should call
// Connection.Close() when done with the call.
//
// This method blocks until a call is received or an error occurs.
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

	connector := ice.IceConnector{SignalingClient: signalingClient, IceAgent: iceAgent, Logger: c.logger}
	conn, err := connector.Connect()
	if err != nil {
		connector.Close()
		return Connection{}, err
	}

	return Connection{conn, connector}, nil
}

// ExternalUserId returns the external user ID associated with the authenticated session.
// This ID can be used to identify the current user in the MAX messenger system.
func (c *Calls) ExternalUserId() string {
	return c.loginData.ExternalUserID
}

// Close closes the underlying OneMe API connection. After calling Close, the Calls instance
// should not be used for making or receiving calls.
//
// Returns an error if closing the connections fails.
func (c *Calls) Close() error {
	return c.onemeClient.Close()
}
