package session

import (
	"context"
	"errors"
	"fmt"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var session *Session = nil

type Session struct {
	Client *spotify.Client
}

func GetToken() (*oauth2.Token, error) {
	token, err := session.Client.Token()
	if err != nil {
		return nil, fmt.Errorf("user not logged in to Spotify: %v", err)
	}
	return token, nil
}

func GetClient() (*spotify.Client, error) {
	if session.Client == nil {
		return nil, errors.New("client is null")
	}
	return session.Client, nil
}

func SetSession(token *oauth2.Token) {
	httpClient := spotifyauth.New().Client(context.Background(), token)

	session = &Session{
		Client: spotify.New(httpClient),
	}
}
