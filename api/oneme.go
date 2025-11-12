package api

import (
	"encoding/json"
	"fmt"

	raw_oneme "github.com/icyfalc0n/max_calls_api/api/raw/oneme"
	"github.com/icyfalc0n/max_calls_api/api/raw/oneme/messages/client"
	"github.com/icyfalc0n/max_calls_api/api/raw/oneme/messages/server"
)

type OnemeApiClient struct {
	RawClient *raw_oneme.ApiClient
	seq       int
}

func NewOnemeApiClient() (OnemeApiClient, error) {
	client, err := raw_oneme.NewApiClient()
	if err != nil {
		return OnemeApiClient{}, err
	}
	c := OnemeApiClient{
		RawClient: client,
		seq:       0,
	}
	err = c.DoClientHello()
	if err != nil {
		return OnemeApiClient{}, err
	}
	return c, nil
}

func SendMessage[T any](c *OnemeApiClient, opcode int, payload T) (int, error) {
	messageSeq := c.seq
	msg := raw_oneme.NewMessage(messageSeq, opcode, payload)
	if err := raw_oneme.SendRawMessage(c.RawClient, msg); err != nil {
		return 0, err
	}
	c.seq++
	return messageSeq, nil
}

func ReceiveMessage[T any](c *OnemeApiClient, seq *int) (T, error) {
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

func (c *OnemeApiClient) SendClientHello() (int, error) {
	return SendMessage(c, client.ClientHelloOpcode(), client.NewClientHello())
}

func (c *OnemeApiClient) DoClientHello() error {
	seq, err := c.SendClientHello()
	if err != nil {
		return err
	}
	_, err = ReceiveMessage[struct{}](c, &seq)
	return err
}

func (c *OnemeApiClient) SendVerificationRequest(phone string) (int, error) {
	req := client.NewVerificationRequest(phone)
	return SendMessage(c, client.VerificationRequestOpcode(), req)
}

func (c *OnemeApiClient) DoVerificationRequest(phone string) (server.VerificationToken, error) {
	seq, err := c.SendVerificationRequest(phone)
	if err != nil {
		return server.VerificationToken{}, err
	}
	return ReceiveMessage[server.VerificationToken](c, &seq)
}

func (c *OnemeApiClient) SendCodeEnter(token, verifyCode string) (int, error) {
	req := client.NewCodeEnter(token, verifyCode)
	return SendMessage(c, client.CodeEnterOpcode(), req)
}

func (c *OnemeApiClient) DoCodeEnter(token, verifyCode string) (server.SuccessfulLogin, error) {
	seq, err := c.SendCodeEnter(token, verifyCode)
	if err != nil {
		return server.SuccessfulLogin{}, err
	}
	return ReceiveMessage[server.SuccessfulLogin](c, &seq)
}

func (c *OnemeApiClient) SendChatSyncRequest(token string) (int, error) {
	req := client.NewChatSyncRequest(token)
	return SendMessage(c, client.ChatSyncRequestOpcode(), req)
}

func (c *OnemeApiClient) DoChatSync(token string) (server.ChatSyncResponse, error) {
	seq, err := c.SendChatSyncRequest(token)
	if err != nil {
		return server.ChatSyncResponse{}, err
	}
	return ReceiveMessage[server.ChatSyncResponse](c, &seq)
}

func (c *OnemeApiClient) SendCallTokenRequest() (int, error) {
	req := client.CallTokenRequest{}
	return SendMessage(c, client.CallTokenRequestOpcode(), req)
}

func (c *OnemeApiClient) DoCallTokenRequest() (string, error) {
	seq, err := c.SendCallTokenRequest()
	if err != nil {
		return "", err
	}
	resp, err := ReceiveMessage[server.CallToken](c, &seq)
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (c *OnemeApiClient) WaitForIncomingCall() (server.IncomingCall, error) {
	const INCOMING_CALL_OPCODE = 137
	for {
		raw, err := c.RawClient.ReceiveRawMessage(nil)
		if err != nil {
			return server.IncomingCall{}, err
		}
		opcodeVal, ok := raw["opcode"].(float64)
		if !ok || int(opcodeVal) != INCOMING_CALL_OPCODE {
			continue
		}

		payloadMap, ok := raw["payload"].(map[string]any)
		if !ok {
			continue
		}

		data, err := json.Marshal(payloadMap)
		if err != nil {
			return server.IncomingCall{}, err
		}

		var payload server.IncomingCallJSON
		if err := json.Unmarshal(data, &payload); err != nil {
			return server.IncomingCall{}, err
		}

		return server.NewIncomingCall(payload)
	}
}
