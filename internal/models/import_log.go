package models

import "gorm.io/gorm"

type ImportLog struct {
	gorm.Model

	Source   string `json:"source" gorm:"index"`
	Season   string `json:"season" gorm:"index"`
	Status   string `json:"status" gorm:"index"`
	Message  string `json:"message"`
	Imported int    `json:"imported"`
	Updated  int    `json:"updated"`
	Skipped  int    `json:"skipped"`
	Errors   string `json:"errors" gorm:"type:text"`
}
