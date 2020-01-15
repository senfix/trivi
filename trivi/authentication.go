package trivi

import "net/http"

func (t *triviApi) Authenticate(body AuthenticationRequest) (err error, response AuthenticationResponse) {
	t.auth = nil
	err = t.execute(http.MethodPost, TRIVI_AUTH, body, &response)
	t.auth = &response
	return
}

type AuthenticationRequest struct {
	AppID     string `json:"appId"`
	AppSecret string `json:"appSecret"`
}

type AuthenticationResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
