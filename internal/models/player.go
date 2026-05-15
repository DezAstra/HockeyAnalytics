package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model

	Name     string
	Position string

	TeamID uint
	Team   Team

	Stats PlayerStats
}
