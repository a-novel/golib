package grpcmocks

import (
	"time"

	"golang.org/x/oauth2"
)

type IDTokenStub struct {
	oauth2.TokenSource
}

var DefaultToken = &oauth2.Token{
	AccessToken:  "secret-access-token",
	TokenType:    "custom-token-type",
	RefreshToken: "refresh-token",
	Expiry:       time.Now().Add(time.Hour),
	ExpiresIn:    999999999,
}

func (i *IDTokenStub) Token() (*oauth2.Token, error) {
	return DefaultToken, nil
}
