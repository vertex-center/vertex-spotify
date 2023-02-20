package session

import (
	"context"
	"errors"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

var session *Session = nil

type Session struct {
	token  oauth2.Token
	Client *spotify.Client
}

func GetSession() (*Session, error) {
	if session == nil {
		return nil, errors.New("user not logged in to Spotify")
	}
	return session, nil
}

func SetSession(token *oauth2.Token) {
	httpClient := spotifyauth.New().Client(context.Background(), token)

	session = &Session{
		token:  *token,
		Client: spotify.New(httpClient),
	}
}
