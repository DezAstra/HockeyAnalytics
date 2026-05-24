package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	NHLID    int `json:"nhl_id" gorm:"index"`
	Name     string
	Position string
	Stats    []PlayerSeasonStats
}
