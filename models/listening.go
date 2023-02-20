package models

import "time"

type Listening struct {
	ID        uint `gorm:"primaryKey"`
	Duration  time.Duration
	Track     Track
	TrackID   uint
	CreatedAt time.Time
}
