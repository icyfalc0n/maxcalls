package main

import (
	"context"
	"fmt"
	"slices"

	callsMessages "github.com/icyfalc0n/max_calls_api/api/calls/messages"
	onemeMessages "github.com/icyfalc0n/max_calls_api/api/oneme/messages"
	"github.com/icyfalc0n/max_calls_api/api/signaling"
	"github.com/icyfalc0n/max_calls_api/api/signaling/messages"
	"github.com/pion/ice/v4"
	"github.com/pion/stun/v3"
)

type IceConnector struct {
	SignalingClient signaling.SignalingClient
	IceAgent        *ice.Agent
}

func (c *IceConnector) Connect() (*ice.Conn, error) {
	err := c.sendLocalCredentials()
	if err != nil {
		return nil, err
	}
	credentials, err := c.receiveRemoteCredentials()
	if err != nil {
		return nil, err
	}

	err = c.gatherICECandidates()
	if err != nil {
		return nil, err
	}
	err = c.receiveRemoteICECandidates()
	if err != nil {
		return nil, err
	}

	return c.establishConnection(credentials)
}

func (c *IceConnector) sendLocalCredentials() error {
	ufrag, password, err := c.IceAgent.GetLocalUserCredentials()
	if err != nil {
		return err
	}
	fmt.Printf("Sending local user credentials. ufrag: %s, password: %s\n", ufrag, password)
	err = c.SignalingClient.SendSignal(messages.Credentials{UFrag: ufrag, Password: password})
	if err != nil {
		return err
	}

	return nil
}

func (c *IceConnector) receiveRemoteCredentials() (messages.Credentials, error) {
	var remoteCredentials messages.Credentials
	fmt.Println("Waiting for remote credentials...")
	err := c.SignalingClient.ReceiveSignal(&remoteCredentials)
	if err != nil {
		return messages.Credentials{}, err
	}
	fmt.Printf("Received remote credentials. ufrag: %s, password: %s\n", remoteCredentials.UFrag, remoteCredentials.Password)

	return remoteCredentials, nil
}

func (c *IceConnector) gatherICECandidates() error {
	iceCandidates := make(chan ice.Candidate)
	defer close(iceCandidates)
	c.IceAgent.OnCandidate(func(c ice.Candidate) {
		iceCandidates <- c
	})
	fmt.Println("Gathering ICE candidates...")
	err := c.IceAgent.GatherCandidates()
	if err != nil {
		return err
	}

	for candidate := range iceCandidates {
		if candidate == nil {
			fmt.Println("ICE gathering complete")
			err = c.SignalingClient.SendSignal(messages.NewCandidate{})
			if err != nil {
				return err
			}
			break
		}
		marshaledCanidate := candidate.Marshal()
		fmt.Printf("New local ICE candidate: %s\n", marshaledCanidate)

		err := c.SignalingClient.SendSignal(messages.NewCandidate{Candidate: marshaledCanidate})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *IceConnector) receiveRemoteICECandidates() error {
	for {
		var remoteCandidateMsg messages.NewCandidate
		err := c.SignalingClient.ReceiveSignal(&remoteCandidateMsg)
		if err != nil {
			return err
		}

		if remoteCandidateMsg.Candidate == "" {
			fmt.Println("Remote ICE gathering complete")
			break
		}

		fmt.Printf("New remote ICE candidate: %s\n", remoteCandidateMsg.Candidate)

		remoteCandidate, err := ice.UnmarshalCandidate(remoteCandidateMsg.Candidate)
		if err != nil {
			return err
		}
		err = c.IceAgent.AddRemoteCandidate(remoteCandidate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *IceConnector) establishConnection(credentials messages.Credentials) (*ice.Conn, error) {
	fmt.Println("Finding an ICE candidate pair...")
	conn, err := c.IceAgent.Dial(context.Background(), credentials.UFrag, credentials.Password)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connection established!")
	return conn, nil
}

func NewAgentFromOutgoing(startedConversationInfo callsMessages.StartedConversationInfo) (*ice.Agent, error) {
	stunIceURIs, err := parseStunServers(startedConversationInfo.StunServer.Urls)
	if err != nil {
		return nil, err
	}
	turnIceURIs, err := parseTurnServers(startedConversationInfo.TurnServer.Urls, startedConversationInfo.TurnServer.Username, startedConversationInfo.TurnServer.Password)
	if err != nil {
		return nil, err
	}
	iceServers := slices.Concat(stunIceURIs, turnIceURIs)

	return ice.NewAgentWithOptions(ice.WithUrls(iceServers))
}

func NewAgentFromIncoming(incomingCall onemeMessages.IncomingCall) (*ice.Agent, error) {
	stunIceURIs, err := parseStunServers([]string{incomingCall.Stun})
	if err != nil {
		return nil, err
	}
	turnIceURIs, err := parseTurnServers(incomingCall.Turn.Servers, incomingCall.Turn.User, incomingCall.Turn.Password)
	if err != nil {
		return nil, err
	}
	iceServers := slices.Concat(stunIceURIs, turnIceURIs)

	return ice.NewAgentWithOptions(ice.WithUrls(iceServers))
}

func parseStunServers(stunServers []string) ([]*stun.URI, error) {
	result := []*stun.URI{}

	for _, server := range stunServers {
		uri, err := stun.ParseURI(server)
		if err != nil {
			return []*stun.URI{}, err
		}
		result = append(result, uri)
	}

	return result, nil
}

func parseTurnServers(turnServers []string, username string, password string) ([]*stun.URI, error) {
	result := []*stun.URI{}

	for _, server := range turnServers {
		uri, err := stun.ParseURI(server)
		if err != nil {
			return []*stun.URI{}, err

		}
		uri.Username = username
		uri.Password = password
		result = append(result, uri)
	}

	return result, nil
}
