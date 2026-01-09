package messages

type VerificationRequest struct {
	Phone    string `json:"phone"`
	Type     string `json:"type"`
	Language string `json:"language"`
}

func NewVerificationRequest(phone string) VerificationRequest {
	return VerificationRequest{
		Phone:    phone,
		Type:     "START_AUTH",
		Language: "ru",
	}
}

func VerificationRequestOpcode() int {
	return 17
}
