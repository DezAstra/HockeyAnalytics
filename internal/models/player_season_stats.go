package models

import "gorm.io/gorm"

type PlayerSeasonStats struct {
	gorm.Model

	PlayerID uint
	Player   Player

	Season string
	Team   string

	GamesPlayed int

	Goals     int
	Assists   int
	PlusMinus int
	Points    int

	PenaltyMinutes int

	EvenStrengthGoals int
	PowerPlayGoals    int
	ShortHandedGoals  int

	Shots int

	BlockedShots int
	Hits         int

	FaceoffsWon    int
	FaceoffsLost   int
	FaceoffPercent *float64 `gorm:"type:decimal(10,2)"`
	TimeOfIce      *float64
}
