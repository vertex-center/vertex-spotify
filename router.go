package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quentinguidee/microservice-core/router"
	"github.com/vertex-center/vertex-spotify/auth"
	"github.com/vertex-center/vertex-spotify/database"
	"github.com/vertex-center/vertex-spotify/session"
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
	url := auth.URL()

	c.Redirect(http.StatusFound, url)
}

func handleAuthCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := auth.Exchange(code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.SetSession(token)

	err = database.SetToken(token)
	if err != nil {
		fmt.Printf("%v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleGetUser(c *gin.Context) {
	sess, err := session.GetSession()
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	user, err := sess.Client.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func handlePlayer(c *gin.Context) {
	sess, err := session.GetSession()
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if c.Query("full") != "" {
		playing, err := sess.Client.PlayerCurrentlyPlaying(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, playing)
		}
		return
	}

	c.JSON(http.StatusOK, currentTrack.ToJSON())
}
