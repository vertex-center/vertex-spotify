package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	vertexrouter "github.com/vertex-center/vertex-core-golang/router"
	"github.com/vertex-center/vertex-spotify/auth"
	"github.com/vertex-center/vertex-spotify/database"
	"github.com/vertex-center/vertex-spotify/session"
	"github.com/vertex-center/vertex-spotify/tracker"
)

func InitializeRouter() *gin.Engine {
	r := vertexrouter.CreateRouter()

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
	authError := c.Query("error")
	if authError != "" {
		err := fmt.Errorf("error: %v", authError)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	code := c.Query("code")
	token, err := auth.Exchange(code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.SetToken(token)

	err = database.SetToken(token)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func handleGetUser(c *gin.Context) {
	client, err := session.GetClient()
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	user, err := client.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func handlePlayer(c *gin.Context) {
	client, err := session.GetClient()
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if c.Query("full") != "" {
		playing, err := client.PlayerCurrentlyPlaying(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, playing)
		}
		return
	}

	c.JSON(http.StatusOK, tracker.GetCurrentTrack().ToJSON())
}
