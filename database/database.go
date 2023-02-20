package database

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
)

var db Database

type Database struct {
	Instance *sqlx.DB
}

func Connect(config Config) {
	dbx, err := sqlx.Connect("postgres", config.ConnectionString())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	db.Instance = dbx

	//_, err = db.Instance.Exec("DROP TABLE sessions")
	_, err = db.Instance.Query(`
		CREATE TABLE IF NOT EXISTS sessions (
		    AccessToken VARCHAR,
		    TokenType VARCHAR,
		    RefreshToken VARCHAR,
		    Expiry VARCHAR
	    )
	`)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func GetToken() (*oauth2.Token, error) {
	var token oauth2.Token
	var expiry string

	row := db.Instance.QueryRow("SELECT * FROM sessions;")
	if row == nil {
		return nil, errors.New("session not found")
	}

	err := row.Scan(&token.AccessToken, &token.TokenType, &token.RefreshToken, &expiry)
	if err != nil {
		return nil, err
	}

	exp, err := time.Parse(time.RFC3339, expiry)
	token.Expiry = exp

	fmt.Printf("%+v\n", token)
	return &token, nil
}

func SetToken(token *oauth2.Token) error {
	tx := db.Instance.MustBegin()

	_ = tx.MustExec("TRUNCATE TABLE sessions")
	_ = tx.MustExec(
		`INSERT INTO sessions (AccessToken, TokenType, RefreshToken, Expiry) VALUES ($1, $2, $3, $4);`,
		token.AccessToken,
		token.TokenType,
		token.RefreshToken,
		token.Expiry.Format(time.RFC3339),
	)

	return tx.Commit()
}
