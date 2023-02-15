package main

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"log"
)

var config Config
var auth *spotifyauth.Authenticator

type Config struct {
	SpotifyClientID     string `env:"SPOTIFY_CLIENT_ID,required"`
	SpotifyClientSecret string `env:"SPOTIFY_CLIENT_SECRET,required"`
	SpotifyRedirectUri  string `env:"SPOTIFY_REDIRECT_URI,required"`
}

func main() {
	loadEnv()

	r := InitializeRouter()

	startTicker()

	err := r.Run(":6150")
	if err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	err = env.Parse(&config)
	if err != nil {
		log.Fatalf("Failed to parse .env to Config: %v", err)
	}

	auth = spotifyauth.New(
		spotifyauth.WithClientID(config.SpotifyClientID),
		spotifyauth.WithClientSecret(config.SpotifyClientSecret),
		spotifyauth.WithRedirectURL(config.SpotifyRedirectUri),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeStreaming,
		),
	)
}
