package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quentinguidee/microservice-core/pubsub"
	"github.com/vertex-center/vertex-spotify/session"
	"github.com/zmb3/spotify/v2"
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
	sess, err := session.GetSession()
	if err != nil {
		fmt.Printf("Failed to ping. User not (yet) logged in.\n")
		return
	}

	player, err := sess.Client.PlayerCurrentlyPlaying(context.Background())
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

	if currentTrack == nil && player.Playing {
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

		return
	}

	if player.Playing {
		currentTrack.listeningTime += 1 * time.Second
	} else if currentTrack != nil {
		currentTrack = nil
		pubsub.Pub("SPOTIFY_PLAYER_CHANGE", []byte(`{"is_playing": false}`))
	}
}
