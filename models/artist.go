package models

type Artist struct {
	ID        uint   `gorm:"primaryKey"`
	SpotifyID string `gorm:"unique"`
	Name      string
	Uri       string
	Url       string
	Albums    []*Album `gorm:"many2many:album_artists;"`
}
