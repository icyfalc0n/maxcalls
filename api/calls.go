package api

import (
	"encoding/json"
	"github.com/google/uuid"
	raw_calls "github.com/icyfalc0n/max_calls_api/api/raw/calls"
	"github.com/icyfalc0n/max_calls_api/api/raw/calls/messages/client"
	"github.com/icyfalc0n/max_calls_api/api/raw/calls/messages/server"
)

type CallsApiClient struct {
	RawClient *raw_calls.ApiClient
}

func (c *CallsApiClient) Login(token string) (server.LoginData, error) {
	sessionData := client.NewSessionData(token)
	data, err := json.Marshal(sessionData)
	if err != nil {
		return server.LoginData{}, err
	}

	response, err := c.RawClient.RunMethod("auth.anonymLogin", map[string]string{
		"session_data": string(data),
	})
	if err != nil {
		return server.LoginData{}, err
	}

	var loginData server.LoginData
	if err := json.Unmarshal(response, &loginData); err != nil {
		return server.LoginData{}, err
	}

	return loginData, nil
}

func (c *CallsApiClient) StartConversation(sessionKey string, callTakerID string, conversationID uuid.UUID) (server.StartedConversationInfo, error) {
	payload := client.NewStartConversationPayload()
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return server.StartedConversationInfo{}, err
	}

	formData := map[string]string{
		"conversationId":  conversationID.String(),
		"isVideo":         "false",
		"protocolVersion": "5",
		"payload":         string(payloadData),
		"externalIds":     callTakerID,
		"session_key":     sessionKey,
	}

	response, err := c.RawClient.RunMethod("vchat.startConversation", formData)
	if err != nil {
		return server.StartedConversationInfo{}, err
	}

	var startedInfo server.StartedConversationInfo
	if err := json.Unmarshal(response, &startedInfo); err != nil {
		return server.StartedConversationInfo{}, err
	}

	return startedInfo, nil
}
