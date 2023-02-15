package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err != nil {
			c.JSON(-1, gin.H{
				"message": err.Error(),
			})
		}
	}
}

func InitializeRouter() *gin.Engine {
	r := gin.Default()

	r.Use(errorMiddleware())

	r.GET("/ping", handlePing)

	authRoutes := r.Group("/auth")
	authRoutes.GET("/login", handleAuthLogin)
	authRoutes.GET("/callback", handleAuthCallback)

	r.GET("/user", handleGetUser)

	playerRoutes := r.Group("/player")
	playerRoutes.GET("", handlePlayer)

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
	code := c.Query("code")
	token, err := auth.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Couldn't retrieve token: %v", err)
	}

	SetSession(token)

	err = db.SetToken(token)
	if err != nil {
		fmt.Printf("%v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleGetUser(c *gin.Context) {
	session, err := GetSession()
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	user, err := session.client.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func handlePlayer(c *gin.Context) {
	session, err := GetSession()
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	playing, err := session.client.PlayerCurrentlyPlaying(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if c.Query("full") != "" {
		c.JSON(http.StatusOK, playing)
		return
	}

	var track gin.H = nil

	if playing.Item != nil {
		track = gin.H{
			"name":   playing.Item.Name,
			"album":  playing.Item.Album.Name,
			"artist": playing.Item.Artists[0].Name,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"is_playing": playing.Playing,
		"track":      track,
	})
}
