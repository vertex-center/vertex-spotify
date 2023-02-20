package models

import (
	"time"

	"golang.org/x/oauth2"
)

type Session struct {
	Id           uint
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s Session) GetToken() (*oauth2.Token, error) {
	expiry, err := time.Parse(time.RFC3339, s.Expiry)
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken:  s.AccessToken,
		TokenType:    s.TokenType,
		RefreshToken: s.RefreshToken,
		Expiry:       expiry,
	}, nil
}
