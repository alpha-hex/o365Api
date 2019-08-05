package o365Api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Token interface {
	GetUserBearerToken() (TokenResponse, error)
}

type TokenRequest struct {
	client_id		string
	client_secret	string
	tenant_id		string
	userName		string
	userPassword	string
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

func (t TokenRequest) GetUserBearerToken() (TokenResponse, error) {
	url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", t.tenant_id)

	payload := strings.NewReader(fmt.Sprintf("grant_type=password&client_id=%s&client_secret=%s&scope=https://graph.microsoft.com/.default&userName=%s&password=%s",
		t.client_id, t.client_secret, t.userName, t.userPassword))
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TokenResponse{}, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resp TokenResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return TokenResponse{}, err
	}

	return resp, nil
}