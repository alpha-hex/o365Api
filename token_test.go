package o365Api

import (
	"testing"
)

func TestGetUserBearerToken(t *testing.T) {
	tokenRequest := TokenRequest{
		client_id: "",
		client_secret: "",
		tenant_id: "",
		userName: "",
		userPassword: "",
	}

	bearerToken , err := Token.GetUserBearerToken(tokenRequest)
	if err != nil {
		t.Error(err)
	}

	if len(bearerToken.AccessToken) == 0 {
		t.Error("Access token is empty")
	}
}