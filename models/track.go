package models

type Track struct {
	ID         uint   `gorm:"primaryKey"`
	SpotifyID  string `gorm:"unique"`
	Name       string
	Duration   int
	Explicit   bool
	Uri        string
	Url        string
	Type       string
	Album      Album
	AlbumID    uint
	Popularity int
}
