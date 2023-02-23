package auth

import (
	"context"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type Config struct {
	SpotifyID          string
	SpotifySecret      string
	SpotifyRedirectUri string
}

var auth *spotifyauth.Authenticator

func Init(config Config) {
	auth = spotifyauth.New(
		spotifyauth.WithClientID(config.SpotifyID),
		spotifyauth.WithClientSecret(config.SpotifySecret),
		spotifyauth.WithRedirectURL(config.SpotifyRedirectUri),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeStreaming,
		),
	)
}

func URL() string {
	return auth.AuthURL("STATE") // TODO: Randomize
}

func Exchange(code string) (*oauth2.Token, error) {
	return auth.Exchange(context.Background(), code)
}
