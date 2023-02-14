package main

import (
	"errors"
	"golang.org/x/oauth2"
)

var session *Session = nil

type Session struct {
	token oauth2.Token
}

func GetSession() (*Session, error) {
	if session == nil {
		return nil, errors.New("user not logged in to Spotify")
	}
	return session, nil
}

func SetSession(token *oauth2.Token) {
	session = &Session{
		token: *token,
	}
}
