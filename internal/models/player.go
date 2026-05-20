package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	NHLID    int `json:"nhl_id"`
	Name     string
	Position string
	Stats    []PlayerSeasonStats
}
