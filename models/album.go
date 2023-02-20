package models

type Album struct {
	ID                   uint   `gorm:"primaryKey"`
	SpotifyID            string `gorm:"unique"`
	Name                 string
	Artists              []*Artist `gorm:"many2many:album_artists;"`
	Group                string
	Type                 string
	Uri                  string
	Url                  string
	ReleaseDate          string
	ReleaseDatePrecision string
	Images               []AlbumImage
}

type AlbumImage struct {
	ID      uint `gorm:"primaryKey"`
	Height  int
	Width   int
	Url     string
	AlbumID uint
}
