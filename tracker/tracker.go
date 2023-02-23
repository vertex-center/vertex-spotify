package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex-core-golang/pubsub"
	"github.com/vertex-center/vertex-spotify/database"
	"github.com/vertex-center/vertex-spotify/models"
	"github.com/vertex-center/vertex-spotify/session"
	"github.com/zmb3/spotify/v2"
)

var logger = console.New("vertex-spotify::tracker")

type CurrentTrack struct {
	listeningTime time.Duration
	track         spotify.FullTrack
}

var currentTrack *CurrentTrack

func GetCurrentTrack() *CurrentTrack {
	return currentTrack
}

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

func Start() {
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
	client, err := session.GetClient()
	if err != nil {
		logger.Error(err)
		return
	}

	player, err := client.PlayerCurrentlyPlaying(context.Background())
	if err != nil {
		logger.Error(fmt.Errorf("failed to get 'player currently playing' info: %v", err))
		return
	}

	spotifyPlaying := player.Playing
	vertexPlaying := currentTrack != nil

	if !vertexPlaying && !spotifyPlaying {
		// Nothing happens, do nothing.
		return
	}

	if vertexPlaying && !spotifyPlaying {
		currentTrack = nil
		logger.Log("Spotify paused.")

		err := pubsub.Pub("SPOTIFY_PLAYER_CHANGE", []byte(`{"is_playing": false}`))
		if err != nil {
			logger.Warn(fmt.Sprintf("couldn't publish 'SPOTIFY_PLAYER_CHANGE' event: %v", err))
		}

		return
	}

	if spotifyPlaying {
		if !vertexPlaying || currentTrack.track.ID != player.Item.ID {
			// play->play: If the track changed, save the track
			// pause->play: Save the track

			if !vertexPlaying {
				logger.Log("Spotify play.")
			} else if vertexPlaying && currentTrack.track.ID != player.Item.ID {
				if currentTrack.listeningTime.Seconds() >= 5 {
					logger.Log(fmt.Sprintf("saving '%s' (%s seconds)...", currentTrack.track.Name, currentTrack.listeningTime))
					err := saveListening()
					if err != nil {
						logger.Error(fmt.Errorf("track changed but failed to save: %v", err))
					} else {
						logger.Log("track changed and saved successfully")
					}
				} else {
					logger.Log(fmt.Sprintf("track '%s' skipped and not saved (%s < 5s)", currentTrack.track.Name, currentTrack.listeningTime))
				}
			}

			currentTrack = &CurrentTrack{
				listeningTime: 0,
				track:         *player.Item,
			}

			err = pubPlayerChange()
			if err != nil {
				logger.Warn(fmt.Sprintf("couldn't publish 'SPOTIFY_PLAYER_CHANGE' event: %v", err))
			}
		} else {
			currentTrack.listeningTime += 1 * time.Second
		}

		return
	}
}

func pubPlayerChange() error {
	message, err := json.Marshal(currentTrack.ToJSON())
	if err != nil {
		return fmt.Errorf("Failed to parse currentTrack info: %v\n", err)
	}

	return pubsub.Pub("SPOTIFY_PLAYER_CHANGE", message)
}

func saveListening() error {
	t := currentTrack.track
	a := currentTrack.track.Album

	var artists []*models.Artist
	for _, artist := range a.Artists {
		artists = append(artists, &models.Artist{
			SpotifyID: string(artist.ID),
			Name:      artist.Name,
			Uri:       string(artist.URI),
			Url:       artist.ExternalURLs["spotify"],
		})
	}

	var images []models.AlbumImage
	for _, image := range a.Images {
		images = append(images, models.AlbumImage{
			Height: image.Height,
			Width:  image.Width,
			Url:    image.URL,
		})
	}

	album := models.Album{
		SpotifyID:            string(a.ID),
		Name:                 a.Name,
		Artists:              artists,
		Group:                a.AlbumGroup,
		Type:                 a.AlbumType,
		Uri:                  string(a.URI),
		Url:                  a.ExternalURLs["spotify"],
		ReleaseDate:          a.ReleaseDate,
		ReleaseDatePrecision: a.ReleaseDatePrecision,
		Images:               images,
	}

	track := models.Track{
		SpotifyID:  string(t.ID),
		Name:       t.Name,
		Duration:   t.Duration,
		Explicit:   t.Explicit,
		Uri:        string(t.URI),
		Url:        t.ExternalURLs["spotify"],
		Type:       t.Type,
		Album:      album,
		Popularity: t.Popularity,
	}

	listening := models.Listening{
		Duration: currentTrack.listeningTime,
		Track:    track,
	}

	currentTrack = nil

	return database.SaveListening(listening)
}
