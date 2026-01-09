package messages

type AcceptCall struct {
	Command       string        `json:"command"`
	Sequence      int           `json:"sequence"`
	MediaSettings MediaSettings `json:"mediaSettings"`
}

type MediaSettings struct {
	IsAudioEnabled             bool `json:"isAudioEnabled"`
	IsVideoEnabled             bool `json:"isVideoEnabled"`
	IsScreenSharingEnabled     bool `json:"isScreenSharingEnabled"`
	IsFastScreenSharingEnabled bool `json:"isFastScreenSharingEnabled"`
	IsAudioSharingEnabled      bool `json:"isAudioSharingEnabled"`
	IsAnimojiEnabled           bool `json:"isAnimojiEnabled"`
}

func NewAcceptCall(sequence int) AcceptCall {
	return AcceptCall{
		Command:  "accept-call",
		Sequence: sequence,
		MediaSettings: MediaSettings{
			IsAudioEnabled:             true,
			IsVideoEnabled:             false,
			IsScreenSharingEnabled:     false,
			IsFastScreenSharingEnabled: false,
			IsAudioSharingEnabled:      false,
			IsAnimojiEnabled:           false,
		},
	}
}
