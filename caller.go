package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/icyfalc0n/max_calls_api/api/calls"
	"github.com/icyfalc0n/max_calls_api/api/oneme"
	"github.com/icyfalc0n/max_calls_api/api/signaling"
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

	for {
		const message = "Hello calltaker!"
		fmt.Printf("[Signaling.Caller] %s\n", message)
		err = signalingClient.SendSignal(message)
		if err != nil {
			panic(err)
		}

		receivedMsg, err := signalingClient.ReceiveSignal()
		if err != nil {
			panic(err)
		}
		fmt.Printf("[Signaling.Calltaker] %v\n", receivedMsg)

		time.Sleep(5 * time.Second)
	}
}
