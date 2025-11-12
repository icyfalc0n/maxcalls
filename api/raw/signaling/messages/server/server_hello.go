package server

type ServerHello struct {
	Conversation Conversation `json:"conversation"`
}

type Conversation struct {
	Participants []Participant `json:"participants"`
}

type Participant struct {
	ExternalID ParticipantExternalID `json:"externalId"`
	Id         int64                 `json:"id"`
}

type ParticipantExternalID struct {
	Id string `json:"id"`
}

func FindUserIDByExternalID(convo ServerHello, externalID string) int64 {
	for _, participant := range convo.Conversation.Participants {
		if participant.ExternalID.Id == externalID {
			return participant.Id
		}
	}
	panic("User id not found in call")
}
