package calls

import (
	"encoding/json"

	"github.com/google/uuid"
	callsMessages "github.com/icyfalc0n/max_calls_api/internal/api/calls/messages"
)

type ApiClient struct {
	RawClient *RawApiClient
}

func (c *ApiClient) Login(token string) (callsMessages.LoginData, error) {
	sessionData := callsMessages.NewSessionData(token)
	data, err := json.Marshal(sessionData)
	if err != nil {
		return callsMessages.LoginData{}, err
	}

	response, err := c.RawClient.CallMethod("auth.anonymLogin", map[string]string{
		"session_data": string(data),
	})
	if err != nil {
		return callsMessages.LoginData{}, err
	}

	var loginData callsMessages.LoginData
	if err := json.Unmarshal(response, &loginData); err != nil {
		return callsMessages.LoginData{}, err
	}

	return loginData, nil
}

func (c *ApiClient) StartConversation(sessionKey string, callTakerID string, conversationID uuid.UUID) (callsMessages.StartedConversationInfo, error) {
	payload := callsMessages.NewStartConversationPayload()
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return callsMessages.StartedConversationInfo{}, err
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
		return callsMessages.StartedConversationInfo{}, err
	}

	var startedInfo callsMessages.StartedConversationInfo
	if err := json.Unmarshal(response, &startedInfo); err != nil {
		return callsMessages.StartedConversationInfo{}, err
	}

	return startedInfo, nil
}
