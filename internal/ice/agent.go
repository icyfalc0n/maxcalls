package ice

import (
	"slices"

	"github.com/icyfalc0n/maxcalls/internal/api/calls/messages"
	onemeMessages "github.com/icyfalc0n/maxcalls/internal/api/oneme/messages"
	pionIce "github.com/pion/ice/v4"
)

func NewAgentFromOutgoing(startedConversationInfo messages.StartedConversationInfo) (*pionIce.Agent, error) {
	stunIceURIs, err := parseStunServers(startedConversationInfo.StunServer.Urls)
	if err != nil {
		return nil, err
	}
	turnIceURIs, err := parseTurnServers(startedConversationInfo.TurnServer.Urls, startedConversationInfo.TurnServer.Username, startedConversationInfo.TurnServer.Password)
	if err != nil {
		return nil, err
	}
	iceServers := slices.Concat(stunIceURIs, turnIceURIs)

	return pionIce.NewAgentWithOptions(pionIce.WithUrls(iceServers))
}

func NewAgentFromIncoming(incomingCall onemeMessages.IncomingCall) (*pionIce.Agent, error) {
	stunIceURIs, err := parseStunServers([]string{incomingCall.Stun})
	if err != nil {
		return nil, err
	}
	turnIceURIs, err := parseTurnServers(incomingCall.Turn.Servers, incomingCall.Turn.User, incomingCall.Turn.Password)
	if err != nil {
		return nil, err
	}
	iceServers := slices.Concat(stunIceURIs, turnIceURIs)

	return pionIce.NewAgentWithOptions(pionIce.WithUrls(iceServers))
}
