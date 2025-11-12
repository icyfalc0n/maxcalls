package client

type TransmitData struct {
	Command         string `json:"command"`
	Sequence        int    `json:"sequence"`
	ParticipantID   int64  `json:"participantId"`
	Data            any    `json:"data"`
	ParticipantType string `json:"participantType"`
}

func NewTransmitData(sequence int, participantID int64, data any) TransmitData {
	return TransmitData{
		Command:         "transmit-data",
		Sequence:        sequence,
		ParticipantID:   participantID,
		Data:            data,
		ParticipantType: "USER",
	}
}
