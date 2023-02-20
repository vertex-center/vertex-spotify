package main

import (
	"log"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	"github.com/quentinguidee/microservice-core/pubsub"
	"github.com/vertex-center/vertex-spotify/database"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var environment Environment
var auth *spotifyauth.Authenticator

type Environment struct {
	SpotifyClientID     string `env:"SPOTIFY_CLIENT_ID,required"`
	SpotifyClientSecret string `env:"SPOTIFY_CLIENT_SECRET,required"`
	SpotifyRedirectUri  string `env:"SPOTIFY_REDIRECT_URI,required"`
	DbUser              string `env:"DB_USER"`
	DbPassword          string `env:"DB_PASSWORD"`
	DbName              string `env:"DB_NAME" envDefault:"spotifyservice"`
}

func main() {
	loadEnv()

	pubsub.InitPubSub()

	r := InitializeRouter()

	err := database.Connect(database.Config{
		User:     environment.DbUser,
		Password: environment.DbPassword,
		Name:     environment.DbName,
	})
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	startTicker()

	token, err := database.GetToken()
	if err == nil {
		SetSession(token)
	} else {
		println(err.Error())
	}

	err = r.Run(":6150")
	if err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	err = env.Parse(&environment)
	if err != nil {
		log.Fatalf("Failed to parse .env to Config: %v", err)
	}

	auth = spotifyauth.New(
		spotifyauth.WithClientID(environment.SpotifyClientID),
		spotifyauth.WithClientSecret(environment.SpotifyClientSecret),
		spotifyauth.WithRedirectURL(environment.SpotifyRedirectUri),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeStreaming,
		),
	)
}
