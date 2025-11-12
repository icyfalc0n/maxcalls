package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/icyfalc0n/max_calls_api/api"
)

func Caller(onemeClient api.OnemeApiClient) {
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

	callsClient := api.CallsApiClient{}
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

	signalingServer, err := api.NewSignalingFromOutgoing(startedConversationInfo, calltakerExternalID)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			fmt.Printf("[Signaling.Calltaker] %v\n", <-signalingServer.ReceiveChannel)
		}
	}()

	for {
		message := "Hello calltaker!"
		fmt.Printf("[Signaling.Caller] %s\n", message)
		signalingServer.SendChannel <- message
		time.Sleep(5 * time.Second)
	}
}
