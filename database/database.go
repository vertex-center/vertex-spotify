package database

import (
	"time"

	"github.com/vertex-center/vertex-spotify/models"
	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var instance Database

type Database struct {
	db *gorm.DB
}

func Connect(config Config) error {
	var err error
	instance.db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: config.DSN(),
	}))
	if err != nil {
		return err
	}

	return instance.db.AutoMigrate(
		&models.Album{},
		&models.AlbumImage{},
		&models.Artist{},
		&models.Listening{},
		&models.Session{},
		&models.Track{},
	)
}

func GetToken() (*oauth2.Token, error) {
	var session models.Session
	result := instance.db.Take(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return session.GetToken()
}

func SetToken(token *oauth2.Token) error {
	session := models.Session{
		ID:           1,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.Format(time.RFC3339),
	}

	result := instance.db.Where(&models.Session{ID: 1}).Updates(&session)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		result = instance.db.Create(&session)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func SaveListening(listening models.Listening) error {
	return instance.db.Transaction(func(tx *gorm.DB) error {
		artists := listening.Track.Album.Artists
		for _, artist := range artists {
			result := tx.Where(&models.Artist{SpotifyID: artist.SpotifyID}).FirstOrCreate(&artist)
			if result.Error != nil {
				return result.Error
			}
		}

		album := &listening.Track.Album
		result := tx.Omit("Images").Where(&models.Album{SpotifyID: album.SpotifyID}).FirstOrCreate(&album)
		if result.Error != nil {
			return result.Error
		}

		albumImages := listening.Track.Album.Images
		for i := range albumImages {
			image := &listening.Track.Album.Images[i]
			result := tx.Where(&models.AlbumImage{
				Height:  image.Height,
				Width:   image.Width,
				AlbumID: album.ID,
			}).FirstOrCreate(&image)
			if result.Error != nil {
				return result.Error
			}
		}

		track := &listening.Track
		result = tx.Where(&models.Track{SpotifyID: track.SpotifyID}).FirstOrCreate(&track)
		if result.Error != nil {
			return result.Error
		}

		return tx.Create(&listening).Error
	})
}
