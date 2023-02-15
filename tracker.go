package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/quentinguidee/microservice-core/pubsub"
	"github.com/zmb3/spotify/v2"
	"time"
)

type CurrentTrack struct {
	listeningTime time.Duration
	track         spotify.FullTrack
}

var currentTrack *CurrentTrack

func (t CurrentTrack) ToJSON() gin.H {
	return gin.H{
		"is_playing": true,
		"track": gin.H{
			"name":   t.track.Name,
			"album":  t.track.Album.Name,
			"artist": t.track.Artists[0].Name,
		},
	}
}

var ticker = time.NewTicker(1500 * time.Millisecond)

func startTicker() {
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				tick()
			case <-done:
				return
			}
		}
	}()
}

func tick() {
	session, err := GetSession()
	if err != nil {
		fmt.Printf("Failed to ping. User not (yet) logged in.\n")
		return
	}

	player, err := session.client.PlayerCurrentlyPlaying(context.Background())
	if err != nil {
		fmt.Printf("Failed to get 'player currently playing' info: %v\n", err)
		return
	}

	// The track changed
	if currentTrack != nil && currentTrack.track.ID != player.Item.ID {
		// TODO: Save the track
		fmt.Printf("Played %s during %s seconds\n", currentTrack.track.Name, currentTrack.listeningTime)

		currentTrack = nil
	}

	if currentTrack == nil {
		if player.Item != nil {
			currentTrack = &CurrentTrack{
				listeningTime: 0,
				track:         *player.Item,
			}

			message, err := json.Marshal(currentTrack.ToJSON())
			if err != nil {
				fmt.Printf("Failed to parse currentTrack info: %v\n", err)
				return
			}

			pubsub.Pub("SPOTIFY_PLAYER_CHANGE", message)
		}
	} else if player.Playing {
		currentTrack.listeningTime += 1 * time.Second
	}
}
