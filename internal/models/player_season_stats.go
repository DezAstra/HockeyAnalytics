package models

import "gorm.io/gorm"

type PlayerSeasonStats struct {
	gorm.Model

	PlayerID uint `gorm:"index;uniqueIndex:idx_player_season"`
	Player   Player

	Season string `gorm:"index;uniqueIndex:idx_player_season"`
	Team   string `gorm:"index"`

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
