package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quentinguidee/microservice-core/router"
)

func InitializeRouter() *gin.Engine {
	r := router.CreateRouter()

	authRoutes := r.Group("/auth")
	authRoutes.GET("/login", handleAuthLogin)
	authRoutes.GET("/callback", handleAuthCallback)

	r.GET("/user", handleGetUser)

	playerRoutes := r.Group("/player")
	playerRoutes.GET("", handlePlayer)

	return r
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

	if c.Query("full") != "" {
		playing, err := session.client.PlayerCurrentlyPlaying(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, playing)
		return
	}

	c.JSON(http.StatusOK, currentTrack.ToJSON())
}
