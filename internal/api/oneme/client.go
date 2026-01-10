package oneme

import (
	"encoding/json"
	"fmt"

	onemeMessages "github.com/icyfalc0n/max_calls_api/internal/api/oneme/messages"
)

type ApiClient struct {
	RawClient *RawApiClient
	seq       int
}

func NewOnemeApiClient() (ApiClient, error) {
	client, err := NewRawApiClient()
	if err != nil {
		return ApiClient{}, err
	}
	c := ApiClient{
		RawClient: client,
		seq:       0,
	}
	err = c.DoClientHello()
	if err != nil {
		return ApiClient{}, err
	}
	return c, nil
}

func SendMessage[T any](c *ApiClient, opcode int, payload T) (int, error) {
	messageSeq := c.seq
	msg := NewMessage(messageSeq, opcode, payload)
	if err := SendRawMessage(c.RawClient, msg); err != nil {
		return 0, err
	}
	c.seq++
	return messageSeq, nil
}

func ReceiveMessage[T any](c *ApiClient, seq *int) (T, error) {
	var result T
	raw, err := c.RawClient.ReceiveRawMessage(seq)
	if err != nil {
		return result, err
	}

	payload, ok := raw["payload"].(map[string]any)
	if !ok {
		return result, fmt.Errorf("invalid payload")
	}

	_, ok = payload["error"].(string)
	if ok {
		rawEncoded, _ := json.Marshal(raw)
		return result, fmt.Errorf("error message received: %s", rawEncoded)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (c *ApiClient) SendClientHello() (int, error) {
	return SendMessage(c, onemeMessages.ClientHelloOpcode(), onemeMessages.NewClientHello())
}

func (c *ApiClient) DoClientHello() error {
	seq, err := c.SendClientHello()
	if err != nil {
		return err
	}
	_, err = ReceiveMessage[struct{}](c, &seq)
	return err
}

func (c *ApiClient) SendVerificationRequest(phone string) (int, error) {
	req := onemeMessages.NewVerificationRequest(phone)
	return SendMessage(c, onemeMessages.VerificationRequestOpcode(), req)
}

func (c *ApiClient) DoVerificationRequest(phone string) (onemeMessages.VerificationToken, error) {
	seq, err := c.SendVerificationRequest(phone)
	if err != nil {
		return onemeMessages.VerificationToken{}, err
	}
	return ReceiveMessage[onemeMessages.VerificationToken](c, &seq)
}

func (c *ApiClient) SendCodeEnter(token, verifyCode string) (int, error) {
	req := onemeMessages.NewCodeEnter(token, verifyCode)
	return SendMessage(c, onemeMessages.CodeEnterOpcode(), req)
}

func (c *ApiClient) DoCodeEnter(token, verifyCode string) (onemeMessages.SuccessfulLogin, error) {
	seq, err := c.SendCodeEnter(token, verifyCode)
	if err != nil {
		return onemeMessages.SuccessfulLogin{}, err
	}
	return ReceiveMessage[onemeMessages.SuccessfulLogin](c, &seq)
}

func (c *ApiClient) SendChatSyncRequest(token string) (int, error) {
	req := onemeMessages.NewChatSyncRequest(token)
	return SendMessage(c, onemeMessages.ChatSyncRequestOpcode(), req)
}

func (c *ApiClient) DoChatSync(token string) (onemeMessages.ChatSyncResponse, error) {
	seq, err := c.SendChatSyncRequest(token)
	if err != nil {
		return onemeMessages.ChatSyncResponse{}, err
	}
	return ReceiveMessage[onemeMessages.ChatSyncResponse](c, &seq)
}

func (c *ApiClient) SendCallTokenRequest() (int, error) {
	req := onemeMessages.CallTokenRequest{}
	return SendMessage(c, onemeMessages.CallTokenRequestOpcode(), req)
}

func (c *ApiClient) DoCallTokenRequest() (string, error) {
	seq, err := c.SendCallTokenRequest()
	if err != nil {
		return "", err
	}
	resp, err := ReceiveMessage[onemeMessages.CallToken](c, &seq)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *ApiClient) WaitForIncomingCall() (onemeMessages.IncomingCall, error) {
	for {
		raw, err := c.RawClient.ReceiveRawMessage(nil)
		if err != nil {
			return onemeMessages.IncomingCall{}, err
		}
		opcodeVal, _ := raw["opcode"].(int)
		if opcodeVal != onemeMessages.IncomingCallOpcode() {
			continue
		}

		payloadMap, ok := raw["payload"].(map[string]any)
		if !ok {
			continue
		}

		data, err := json.Marshal(payloadMap)
		if err != nil {
			return onemeMessages.IncomingCall{}, err
		}

		var payload onemeMessages.IncomingCallJSON
		if err := json.Unmarshal(data, &payload); err != nil {
			return onemeMessages.IncomingCall{}, err
		}

		return onemeMessages.NewIncomingCall(payload)
	}
}
