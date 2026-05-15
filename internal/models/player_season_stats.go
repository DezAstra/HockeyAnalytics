package models

import "gorm.io/gorm"

type PlayerSeasonStats struct {
	gorm.Model

	PlayerID uint
	Player   Player

	Season string

	GamesPlayed    int
	Goals          int
	Assists        int
	PlusMinus      int
	PenaltyMinutes int

	TimeOfIce *float64
}
