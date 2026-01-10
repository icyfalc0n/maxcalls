package ice

import (
	"context"

	"github.com/icyfalc0n/maxcalls/internal/api/signaling"
	signalingMessages "github.com/icyfalc0n/maxcalls/internal/api/signaling/messages"
	"github.com/icyfalc0n/maxcalls/logger"
	"github.com/pion/ice/v4"
)

type IceConnector struct {
	SignalingClient signaling.SignalingClient
	IceAgent        *ice.Agent
	Logger          logger.Logger
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
	c.Logger.Debugf("Sending local user credentials. ufrag: %s, password: %s\n", ufrag, password)
	err = c.SignalingClient.SendSignal(signalingMessages.Credentials{UFrag: ufrag, Password: password})
	if err != nil {
		return err
	}

	return nil
}

func (c *IceConnector) receiveRemoteCredentials() (signalingMessages.Credentials, error) {
	var remoteCredentials signalingMessages.Credentials
	c.Logger.Debugf("Waiting for remote credentials...")
	err := c.SignalingClient.ReceiveSignal(&remoteCredentials)
	if err != nil {
		return signalingMessages.Credentials{}, err
	}
	c.Logger.Debugf("Received remote credentials. ufrag: %s, password: %s\n", remoteCredentials.UFrag, remoteCredentials.Password)

	return remoteCredentials, nil
}

func (c *IceConnector) gatherICECandidates() error {
	iceCandidates := make(chan ice.Candidate)
	defer close(iceCandidates)
	err := c.IceAgent.OnCandidate(func(c ice.Candidate) {
		iceCandidates <- c
	})
	if err != nil {
		return err
	}

	c.Logger.Debugf("Gathering ICE candidates...")
	err = c.IceAgent.GatherCandidates()
	if err != nil {
		return err
	}

	for candidate := range iceCandidates {
		if candidate == nil {
			c.Logger.Debugf("ICE gathering complete")
			err = c.SignalingClient.SendSignal(signalingMessages.NewCandidate{})
			if err != nil {
				return err
			}
			break
		}
		marshaledCandidate := candidate.Marshal()
		c.Logger.Debugf("New local ICE candidate: %s\n", marshaledCandidate)

		err := c.SignalingClient.SendSignal(signalingMessages.NewCandidate{Candidate: marshaledCandidate})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *IceConnector) receiveRemoteICECandidates() error {
	for {
		var remoteCandidateMsg signalingMessages.NewCandidate
		err := c.SignalingClient.ReceiveSignal(&remoteCandidateMsg)
		if err != nil {
			return err
		}

		if remoteCandidateMsg.Candidate == "" {
			c.Logger.Debugf("Remote ICE gathering complete")
			break
		}

		c.Logger.Debugf("New remote ICE candidate: %s\n", remoteCandidateMsg.Candidate)

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

func (c *IceConnector) establishConnection(credentials signalingMessages.Credentials) (*ice.Conn, error) {
	c.Logger.Debugf("Finding an ICE candidate pair...")
	conn, err := c.IceAgent.Dial(context.Background(), credentials.UFrag, credentials.Password)
	if err != nil {
		return nil, err
	}

	c.Logger.Debugf("Connection established!")
	return conn, nil
}

func (c *IceConnector) Close() error {
	err := c.IceAgent.Close()
	if err != nil {
		c.SignalingClient.Close()
		return err
	}

	return c.SignalingClient.Close()
}
