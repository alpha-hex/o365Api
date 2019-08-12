package o365Api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Token interface {
	GetUserBearerToken() (TokenResponse, error)
}

type TokenRequest struct {
	Client_ID     string
	Client_Secret string
	Tenant_ID     string
	UserName      string
	UserPassword  string
}

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

func (t TokenRequest) GetUserBearerToken() (TokenResponse, error) {
	if len(t.Client_ID) == 0 || len(t.Client_Secret) == 0 || len(t.Tenant_ID) == 0 || len(t.UserName) == 0 || len(t.UserPassword) == 0 {
		return TokenResponse{}, errors.New("TokenRequest is not valid")
	}

	url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", t.tenant_id)

	payload := strings.NewReader(fmt.Sprintf("grant_type=password&client_id=%s&client_secret=%s&scope=https://graph.microsoft.com/.default&userName=%s&password=%s",
		t.Client_ID, t.Client_Secret, t.UserName, t.UserPassword))
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
