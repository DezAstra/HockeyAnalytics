package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"hockeyAnalytics/internal/services"

	"github.com/gin-gonic/gin"
)

type ImportLogResponse struct {
	ID        uint     `json:"id"`
	CreatedAt string   `json:"created_at"`
	Source    string   `json:"source"`
	Season    string   `json:"season"`
	Status    string   `json:"status"`
	Message   string   `json:"message"`
	Imported  int      `json:"imported"`
	Updated   int      `json:"updated"`
	Skipped   int      `json:"skipped"`
	Errors    []string `json:"errors"`
}

func GetImportLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))

	logs, err := services.ListImportLogs(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]ImportLogResponse, 0, len(logs))
	for _, item := range logs {
		errors := make([]string, 0)
		if strings.TrimSpace(item.Errors) != "" {
			errors = strings.Split(item.Errors, "\n")
		}

		response = append(response, ImportLogResponse{
			ID:        item.ID,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			Source:    item.Source,
			Season:    item.Season,
			Status:    item.Status,
			Message:   item.Message,
			Imported:  item.Imported,
			Updated:   item.Updated,
			Skipped:   item.Skipped,
			Errors:    errors,
		})
	}

	c.JSON(http.StatusOK, response)
}
