package calls

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/icyfalc0n/max_calls_api/api/calls/messages"
)

type CallsApiClient struct {
	RawClient *RawApiClient
}

func (c *CallsApiClient) Login(token string) (messages.LoginData, error) {
	sessionData := messages.NewSessionData(token)
	data, err := json.Marshal(sessionData)
	if err != nil {
		return messages.LoginData{}, err
	}

	response, err := c.RawClient.CallMethod("auth.anonymLogin", map[string]string{
		"session_data": string(data),
	})
	if err != nil {
		return messages.LoginData{}, err
	}

	var loginData messages.LoginData
	if err := json.Unmarshal(response, &loginData); err != nil {
		return messages.LoginData{}, err
	}

	return loginData, nil
}

func (c *CallsApiClient) StartConversation(sessionKey string, callTakerID string, conversationID uuid.UUID) (messages.StartedConversationInfo, error) {
	payload := messages.NewStartConversationPayload()
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return messages.StartedConversationInfo{}, err
	}

	formData := map[string]string{
		"conversationId":  conversationID.String(),
		"isVideo":         "false",
		"protocolVersion": "5",
		"payload":         string(payloadData),
		"externalIds":     callTakerID,
		"session_key":     sessionKey,
	}

	response, err := c.RawClient.CallMethod("vchat.startConversation", formData)
	if err != nil {
		return messages.StartedConversationInfo{}, err
	}

	var startedInfo messages.StartedConversationInfo
	if err := json.Unmarshal(response, &startedInfo); err != nil {
		return messages.StartedConversationInfo{}, err
	}

	return startedInfo, nil
}
