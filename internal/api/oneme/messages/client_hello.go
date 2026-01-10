package messages

import (
	"github.com/google/uuid"
)

type DeviceType string

const (
	DeviceTypeWeb DeviceType = "WEB"
)

type Locale string

const (
	LocaleRussian Locale = "ru"
)

type OSVersion string

const (
	OSVersionWindows OSVersion = "Windows"
)

type DeviceName string

const (
	DeviceNameChrome DeviceName = "Chrome"
)

type UserAgentHeader string

const (
	UserAgentHeaderChrome = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36"
)

const AppVersion string = "25.11.2"

type Screen string

const (
	ScreenPC1080p Screen = "1080x1920 1.0x"
)

type Timezone string

const (
	TimezoneMoscow Timezone = "Europe/Moscow"
)

type UserAgent struct {
	DeviceType      DeviceType      `json:"deviceType"`
	Locale          Locale          `json:"locale"`
	DeviceLocale    Locale          `json:"deviceLocale"`
	OSVersion       OSVersion       `json:"osVersion"`
	DeviceName      DeviceName      `json:"deviceName"`
	UserAgentHeader UserAgentHeader `json:"headerUserAgent"`
	AppVersion      string          `json:"appVersion"`
	Screen          Screen          `json:"screen"`
	Timezone        Timezone        `json:"timezone"`
}

func NewUserAgent() UserAgent {
	return UserAgent{
		DeviceType:      DeviceTypeWeb,
		Locale:          LocaleRussian,
		DeviceLocale:    LocaleRussian,
		OSVersion:       OSVersionWindows,
		DeviceName:      DeviceNameChrome,
		UserAgentHeader: UserAgentHeaderChrome,
		AppVersion:      AppVersion,
		Screen:          ScreenPC1080p,
		Timezone:        TimezoneMoscow,
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
