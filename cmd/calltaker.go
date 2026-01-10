package main

import (
	"fmt"
	"os"

	"github.com/icyfalc0n/max_calls_api/internal/api/calls"
	"github.com/icyfalc0n/max_calls_api/internal/api/oneme"
	"github.com/icyfalc0n/max_calls_api/internal/api/signaling"
	"github.com/icyfalc0n/max_calls_api/internal/ice"
)

func Calltaker(onemeClient oneme.ApiClient) {
	onemeAuthTokenBytes, err := os.ReadFile("token_calltaker")
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

	fmt.Printf("Call taker ID: %s\n", loginData.ExternalUserID)
	fmt.Println("Waiting for incoming calls....")

	incomingCall, err := onemeClient.WaitForIncomingCall()
	if err != nil {
		panic(err)
	}
	signalingClient, err := signaling.NewSignalingFromIncoming(incomingCall, loginData)
	if err != nil {
		panic(err)
	}

	iceAgent, err := ice.NewAgentFromIncoming(incomingCall)
	if err != nil {
		panic(err)
	}

	iceConnector := ice.IceConnector{SignalingClient: signalingClient, IceAgent: iceAgent}
	_, err = iceConnector.Connect()
	if err != nil {
		panic(err)
	}
}
