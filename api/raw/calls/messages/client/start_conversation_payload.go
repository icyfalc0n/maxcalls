package client

type StartConversationPayload struct {
	IsVideo bool `json:"is_video"`
}

func NewStartConversationPayload() StartConversationPayload {
	return StartConversationPayload{
		IsVideo: false,
	}
}
