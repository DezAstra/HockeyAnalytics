package models

import "gorm.io/gorm"

type PlayerStats struct {
	gorm.Model

	PlayerID uint

	GamesPlayed    int
	Goals          int
	Assists        int
	PlusMinus      int
	PenaltyMinutes int
	TimeOfIce      *float64
}
