package oneme

type MessagingSide int

const (
	MessagingSideClient MessagingSide = 0
	MessagingSideServer MessagingSide = 1
)

type Message[T any] struct {
	SequenceNumber int           `json:"seq"`
	Opcode         int           `json:"opcode"`
	Payload        T             `json:"payload"`
	Version        int           `json:"ver"`
	Side           MessagingSide `json:"cmd"`
}

func NewMessage[T any](sequenceNumber int, opcode int, payload T) Message[T] {
	return Message[T]{
		SequenceNumber: sequenceNumber,
		Opcode:         opcode,
		Payload:        payload,
		Version:        11,
		Side:           MessagingSideClient,
	}
}
