package main

import (
	"github.com/gin-gonic/gin"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"net/http"
)

func InitializeRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", handlePing)

	authRoutes := r.Group("/auth")
	authRoutes.GET("/login", handleAuthLogin)
	authRoutes.GET("/callback", handleAuthCallback)

	return r
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func handleAuthLogin(c *gin.Context) {
	url := auth.AuthURL("STATE") // TODO: Randomize

	c.Redirect(http.StatusFound, url)
}

func handleAuthCallback(c *gin.Context) {
	credentialsConfig := &clientcredentials.Config{
		ClientID:     config.SpotifyClientID,
		ClientSecret: config.SpotifyClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := credentialsConfig.Token(c)
	if err != nil {
		log.Fatalf("Couldn't retrieve token: %v", err)
	} else {
		log.Printf("Token: %s", token)
	}

	SetSession(token)

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
