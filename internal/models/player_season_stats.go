package models

import "gorm.io/gorm"

type PlayerSeasonStats struct {
	gorm.Model

	PlayerID uint
	Player   Player

	Season string

	GamesPlayed int

	Goals     int
	Assists   int
	PlusMinus int

	PenaltyMinutes int

	EvenStrengthGoals int
	PowerPlayGoals    int
	ShortHandedGoals  int

	Shots int

	BlockedShots int
	Hits         int

	FaceoffsWon  int
	FaceoffsLost int

	TimeOfIce *float64
}
