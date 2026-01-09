package messages

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/pierrec/lz4/v4"
)

type IncomingCallJSON struct {
	Vcp            string `json:"vcp"`
	CallerID       int    `json:"callerId"`
	ConversationID string `json:"conversationId"`
}

type VcpDecoded struct {
	SignalingToken  string `json:"tkn"`
	SignalingServer string `json:"wse"`
	StunServer      string `json:"stne"`
	TurnServers     string `json:"trne"`
	TurnUser        string `json:"trnu"`
	TurnPassword    string `json:"trnp"`
}

type TurnServer struct {
	Servers  []string
	User     string
	Password string
}

type SignalingServer struct {
	Token string
	URL   string
}

type IncomingCall struct {
	Turn           TurnServer
	Signaling      SignalingServer
	Stun           string
	CallerID       int
	ConversationID string
}

func decodeVcp(vcp string) (VcpDecoded, error) {
	parts := strings.SplitN(vcp, ":", 2)
	if len(parts) != 2 {
		return VcpDecoded{}, errors.New("invalid vcp format")
	}
	uncompressedSizeStr := parts[0]
	compressedBase64 := parts[1]

	compressed, err := base64.StdEncoding.DecodeString(compressedBase64)
	if err != nil {
		return VcpDecoded{}, err
	}

	uncompressedSize := 0
	_, err = fmt.Sscanf(uncompressedSizeStr, "%d", &uncompressedSize)
	if err != nil {
		return VcpDecoded{}, err
	}

	decompressed := make([]byte, uncompressedSize)
	if _, err := lz4.UncompressBlock(compressed, decompressed); err != nil {
		return VcpDecoded{}, err
	}

	var decoded VcpDecoded
	if err := json.Unmarshal(decompressed, &decoded); err != nil {
		return VcpDecoded{}, err
	}

	return decoded, nil
}

func NewIncomingCall(raw IncomingCallJSON) (IncomingCall, error) {
	decodedVcp, err := decodeVcp(raw.Vcp)
	if err != nil {
		return IncomingCall{}, err
	}

	turn := TurnServer{
		Servers:  strings.Split(decodedVcp.TurnServers, ","),
		User:     decodedVcp.TurnUser,
		Password: decodedVcp.TurnPassword,
	}

	signaling := SignalingServer{
		Token: decodedVcp.SignalingToken,
		URL:   decodedVcp.SignalingServer,
	}

	return IncomingCall{
		Turn:           turn,
		Signaling:      signaling,
		Stun:           decodedVcp.StunServer,
		CallerID:       raw.CallerID,
		ConversationID: raw.ConversationID,
	}, nil
}
