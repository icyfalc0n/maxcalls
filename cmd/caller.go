package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/icyfalc0n/max_calls_api/internal/api/calls"
	"github.com/icyfalc0n/max_calls_api/internal/api/oneme"
	"github.com/icyfalc0n/max_calls_api/internal/api/signaling"
	"github.com/icyfalc0n/max_calls_api/internal/ice"
)

func Caller(onemeClient oneme.ApiClient) {
	onemeAuthTokenBytes, err := os.ReadFile("token_caller")
	if err != nil {
		panic(err)
	}
	onemeAuthToken := string(onemeAuthTokenBytes)

	_, err = onemeClient.DoChatSync(onemeAuthToken)
	if err != nil {
		panic(err)
	}

	callsLoginToken, err := onemeClient.DoCallTokenRequest()
	if err != nil {
		panic(err)
	}

	callsClient := calls.ApiClient{}
	loginData, err := callsClient.Login(callsLoginToken)
	if err != nil {
		panic(err)
	}

	reader := StdinReader{Reader: bufio.NewReader(os.Stdin)}
	fmt.Printf("Caller ID: %s\n", loginData.ExternalUserID)
	fmt.Printf("Call taker ID: ")
	calltakerExternalID := reader.Read()

	startedConversationInfo, err := callsClient.StartConversation(loginData.SessionKey, calltakerExternalID, uuid.New())
	if err != nil {
		panic(err)
	}

	signalingClient, err := signaling.NewSignalingFromOutgoing(startedConversationInfo.Endpoint, calltakerExternalID)
	if err != nil {
		panic(err)
	}

	iceAgent, err := ice.NewAgentFromOutgoing(startedConversationInfo)
	if err != nil {
		panic(err)
	}
	defer iceAgent.Close()

	iceConnector := ice.IceConnector{SignalingClient: signalingClient, IceAgent: iceAgent}
	_, err = iceConnector.Connect()
	if err != nil {
		panic(err)
	}
}
