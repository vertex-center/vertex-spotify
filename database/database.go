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
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: config.DSN(),
	}))

	instance.db = db

	err = instance.db.AutoMigrate(&models.Session{})
	if err != nil {
		return err
	}

	return nil
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
		Id:           1,
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.Format(time.RFC3339),
	}

	result := instance.db.Where(&models.Session{Id: 1}).Updates(&session)
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
