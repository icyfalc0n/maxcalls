package messages

import (
	"github.com/google/uuid"
)

type SessionData struct {
	Token         string `json:"auth_token"`
	ClientType    string `json:"client_type"`
	ClientVersion string `json:"client_version"`
	DeviceID      string `json:"device_id"`
	Version       int    `json:"version"`
}

func NewSessionData(token string) SessionData {
	return SessionData{
		Token:         token,
		ClientType:    "SDK_JS",
		ClientVersion: "1.1",
		DeviceID:      uuid.NewString(),
		Version:       3,
	}
}
