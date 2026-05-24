package services

import (
	"strings"

	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
)

type ImportLogInput struct {
	Source   string
	Season   string
	Status   string
	Message  string
	Imported int
	Updated  int
	Skipped  int
	Errors   []string
}

func CreateImportLog(input ImportLogInput) error {
	log := models.ImportLog{
		Source:   input.Source,
		Season:   input.Season,
		Status:   input.Status,
		Message:  input.Message,
		Imported: input.Imported,
		Updated:  input.Updated,
		Skipped:  input.Skipped,
		Errors:   strings.Join(input.Errors, "\n"),
	}

	return database.DB.Create(&log).Error
}

func ListImportLogs(limit int) ([]models.ImportLog, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}

	var logs []models.ImportLog
	err := database.DB.
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error

	return logs, err
}
