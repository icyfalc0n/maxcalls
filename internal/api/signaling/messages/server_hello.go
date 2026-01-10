package messages

import "fmt"

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

func FindUserIDByExternalID(convo ServerHello, externalID string) (int64, error) {
	for _, participant := range convo.Conversation.Participants {
		if participant.ExternalID.Id == externalID {
			return participant.Id, nil
		}
	}
	return 0, fmt.Errorf("could not find user id %s", externalID)
}
