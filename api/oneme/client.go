package oneme

import (
	"encoding/json"
	"fmt"

	"github.com/icyfalc0n/max_calls_api/api/oneme/messages"
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
	return SendMessage(c, messages.ClientHelloOpcode(), messages.NewClientHello())
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
	req := messages.NewVerificationRequest(phone)
	return SendMessage(c, messages.VerificationRequestOpcode(), req)
}

func (c *ApiClient) DoVerificationRequest(phone string) (messages.VerificationToken, error) {
	seq, err := c.SendVerificationRequest(phone)
	if err != nil {
		return messages.VerificationToken{}, err
	}
	return ReceiveMessage[messages.VerificationToken](c, &seq)
}

func (c *ApiClient) SendCodeEnter(token, verifyCode string) (int, error) {
	req := messages.NewCodeEnter(token, verifyCode)
	return SendMessage(c, messages.CodeEnterOpcode(), req)
}

func (c *ApiClient) DoCodeEnter(token, verifyCode string) (messages.SuccessfulLogin, error) {
	seq, err := c.SendCodeEnter(token, verifyCode)
	if err != nil {
		return messages.SuccessfulLogin{}, err
	}
	return ReceiveMessage[messages.SuccessfulLogin](c, &seq)
}

func (c *ApiClient) SendChatSyncRequest(token string) (int, error) {
	req := messages.NewChatSyncRequest(token)
	return SendMessage(c, messages.ChatSyncRequestOpcode(), req)
}

func (c *ApiClient) DoChatSync(token string) (messages.ChatSyncResponse, error) {
	seq, err := c.SendChatSyncRequest(token)
	if err != nil {
		return messages.ChatSyncResponse{}, err
	}
	return ReceiveMessage[messages.ChatSyncResponse](c, &seq)
}

func (c *ApiClient) SendCallTokenRequest() (int, error) {
	req := messages.CallTokenRequest{}
	return SendMessage(c, messages.CallTokenRequestOpcode(), req)
}

func (c *ApiClient) DoCallTokenRequest() (string, error) {
	seq, err := c.SendCallTokenRequest()
	if err != nil {
		return "", err
	}
	resp, err := ReceiveMessage[messages.CallToken](c, &seq)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *ApiClient) WaitForIncomingCall() (messages.IncomingCall, error) {
	for {
		raw, err := c.RawClient.ReceiveRawMessage(nil)
		if err != nil {
			return messages.IncomingCall{}, err
		}
		opcodeVal, ok := raw["opcode"].(float64)
		if !ok || int(opcodeVal) != messages.IncomingCallOpcode() {
			continue
		}

		payloadMap, ok := raw["payload"].(map[string]any)
		if !ok {
			continue
		}

		data, err := json.Marshal(payloadMap)
		if err != nil {
			return messages.IncomingCall{}, err
		}

		var payload messages.IncomingCallJSON
		if err := json.Unmarshal(data, &payload); err != nil {
			return messages.IncomingCall{}, err
		}

		return messages.NewIncomingCall(payload)
	}
}
