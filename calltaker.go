package main

import (
	"fmt"
	"os"

	"github.com/icyfalc0n/max_calls_api/api/calls"
	"github.com/icyfalc0n/max_calls_api/api/oneme"
	"github.com/icyfalc0n/max_calls_api/api/signaling"
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
	signalingClient, err := signaling.NewSignalingFromIncoming(incomingCall, loginData)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
		}
	}()

	for {
		receivedMessage, err := signalingClient.ReceiveSignal()
		if err != nil {
			panic(err)
		}
		fmt.Printf("[Signaling.Caller] %v\n", receivedMessage)

		const message = "Hello caller!"
		err = signalingClient.SendSignal(message)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[Signaling.Calltaker] %s\n", message)
	}
}
