package client

type CodeEnter struct {
	Token         string `json:"token"`
	VerifyCode    string `json:"verifyCode"`
	AuthTokenType string `json:"authTokenType"`
}

func NewCodeEnter(token string, verifyCode string) CodeEnter {
	return CodeEnter{
		Token:         token,
		VerifyCode:    verifyCode,
		AuthTokenType: "CHECK_CODE",
	}
}

func CodeEnterOpcode() int {
	return 18
}
