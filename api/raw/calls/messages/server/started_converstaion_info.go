package server

type TurnServer struct {
	Urls     []string `json:"urls"`
	Username string   `json:"username"`
	Password string   `json:"credential"`
}

type StunServer struct {
	Urls []string `json:"urls"`
}

type StartedConversationInfo struct {
	TurnServer TurnServer `json:"turn_server"`
	StunServer StunServer `json:"stun_server"`
	Endpoint   string     `json:"endpoint"`
}
