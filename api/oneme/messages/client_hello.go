package messages

import (
	"github.com/google/uuid"
)

const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36"

type DeviceType string

const (
	DeviceTypeWEB DeviceType = "WEB"
)

type UserAgent struct {
	DeviceType      DeviceType `json:"deviceType"`
	Locale          string     `json:"locale"`
	DeviceLocale    string     `json:"deviceLocale"`
	OSVersion       string     `json:"osVersion"`
	DeviceName      string     `json:"deviceName"`
	HeaderUserAgent string     `json:"headerUserAgent"`
	AppVersion      string     `json:"appVersion"`
	Screen          string     `json:"screen"`
	Timezone        string     `json:"timezone"`
}

func NewUserAgent() UserAgent {
	return UserAgent{
		DeviceType:      DeviceTypeWEB,
		Locale:          "ru",
		DeviceLocale:    "ru",
		OSVersion:       "Windows",
		DeviceName:      "Chrome",
		HeaderUserAgent: USER_AGENT,
		AppVersion:      "25.11.2",
		Screen:          "1080x1920 1.0x",
		Timezone:        "Europe/Moscow",
	}
}

type ClientHello struct {
	UserAgent UserAgent `json:"userAgent"`
	DeviceID  string    `json:"deviceId"`
}

func NewClientHello() ClientHello {
	return ClientHello{
		UserAgent: NewUserAgent(),
		DeviceID:  uuid.NewString(),
	}
}

func ClientHelloOpcode() int {
	return 6
}
