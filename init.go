package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex-core-golang/pubsub"
	"github.com/vertex-center/vertex-spotify/auth"
	"github.com/vertex-center/vertex-spotify/database"
	"github.com/vertex-center/vertex-spotify/router"
	"github.com/vertex-center/vertex-spotify/session"
	"github.com/vertex-center/vertex-spotify/tracker"
)

var logger = console.New("vertex-spotify::init")
var environment Environment

type Environment struct {
	SpotifyID          string `env:"SPOTIFY_ID,required"`
	SpotifySecret      string `env:"SPOTIFY_SECRET,required"`
	SpotifyRedirectUri string `env:"SPOTIFY_REDIRECT_URI,required"`
	DbUser             string `env:"DB_USER"`
	DbPassword         string `env:"DB_PASSWORD"`
	DbName             string `env:"DB_NAME" envDefault:"spotifyservice"`
}

func main() {
	loadEnv()

	auth.Init(auth.Config{
		SpotifyID:          environment.SpotifyID,
		SpotifySecret:      environment.SpotifySecret,
		SpotifyRedirectUri: environment.SpotifyRedirectUri,
	})

	pubsub.InitPubSub()

	r := router.InitializeRouter()

	err := database.Connect(database.Config{
		User:     environment.DbUser,
		Password: environment.DbPassword,
		Name:     environment.DbName,
	})
	if err != nil {
		logger.Error(err)
		return
	}

	tracker.Start()

	token, err := database.GetToken()
	if err == nil {
		session.SetToken(token)
	} else {
		logger.Error(err)
	}

	err = r.Run(":6150")
	if err != nil {
		logger.Error(fmt.Errorf("error while starting server: %v", err))
		os.Exit(1)
		return
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		logger.Error(fmt.Errorf("error loading .env file: %v", err))
	}

	err = env.Parse(&environment)
	if err != nil {
		logger.Error(fmt.Errorf("failed to parse .env to Config: %v", err))
	}
}
