package server

type LoginData struct {
	UID              string `json:"uid"`
	SessionKey       string `json:"session_key"`
	SessionSecretKey string `json:"session_secret_key"`
	APIServer        string `json:"api_server"`
	ExternalUserID   string `json:"external_user_id"`
}
