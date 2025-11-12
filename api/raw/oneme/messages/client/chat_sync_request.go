package client

type ChatSyncRequest struct {
	Token        string `json:"token"`
	Interactive  bool   `json:"interactive"`
	ChatsCount   int    `json:"chatsCount"`
	ChatsSync    int    `json:"chatsSync"`
	ContactsSync int    `json:"contactsSync"`
	PresenceSync int    `json:"presenceSync"`
	DraftsSync   int    `json:"draftsSync"`
}

func NewChatSyncRequest(token string) ChatSyncRequest {
	return ChatSyncRequest{
		Token:        token,
		Interactive:  false,
		ChatsCount:   40,
		ChatsSync:    0,
		ContactsSync: 0,
		PresenceSync: 0,
		DraftsSync:   0,
	}
}

func ChatSyncRequestOpcode() int {
	return 19
}
