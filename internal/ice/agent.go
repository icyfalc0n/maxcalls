package ice

import (
	"slices"

	"github.com/icyfalc0n/max_calls_api/internal/api/calls/messages"
	messages2 "github.com/icyfalc0n/max_calls_api/internal/api/oneme/messages"
	ice2 "github.com/pion/ice/v4"
)

func NewAgentFromOutgoing(startedConversationInfo messages.StartedConversationInfo) (*ice2.Agent, error) {
	stunIceURIs, err := parseStunServers(startedConversationInfo.StunServer.Urls)
	if err != nil {
		return nil, err
	}
	turnIceURIs, err := parseTurnServers(startedConversationInfo.TurnServer.Urls, startedConversationInfo.TurnServer.Username, startedConversationInfo.TurnServer.Password)
	if err != nil {
		return nil, err
	}
	iceServers := slices.Concat(stunIceURIs, turnIceURIs)

	return ice2.NewAgentWithOptions(ice2.WithUrls(iceServers))
}

func NewAgentFromIncoming(incomingCall messages2.IncomingCall) (*ice2.Agent, error) {
	stunIceURIs, err := parseStunServers([]string{incomingCall.Stun})
	if err != nil {
		return nil, err
	}
	turnIceURIs, err := parseTurnServers(incomingCall.Turn.Servers, incomingCall.Turn.User, incomingCall.Turn.Password)
	if err != nil {
		return nil, err
	}
	iceServers := slices.Concat(stunIceURIs, turnIceURIs)

	return ice2.NewAgentWithOptions(ice2.WithUrls(iceServers))
}
