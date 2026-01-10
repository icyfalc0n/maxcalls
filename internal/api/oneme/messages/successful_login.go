package messages

type Login struct {
	Token string `json:"token"`
}

type TokenAttributes struct {
	Login Login `json:"LOGIN"`
}

type SuccessfulLogin struct {
	TokenAttributes TokenAttributes `json:"tokenAttrs"`
}
