package handlers

import (
	"hockeyAnalytics/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SyncNHLSeason(
	c *gin.Context,
) {
	season :=
		c.DefaultQuery(
			"season",
			"20232024",
		)

	result, err :=
		services.ImportSeason(
			season,
		)

	if err != nil {
		_ = services.CreateImportLog(services.ImportLogInput{
			Source:   "NHL API",
			Season:   result.DisplaySeason,
			Status:   "error",
			Message:  err.Error(),
			Imported: result.PlayersCreated + result.StatsCreated,
			Updated:  result.PlayersUpdated + result.StatsUpdated,
			Skipped:  result.Skipped,
			Errors:   append(result.Errors, err.Error()),
		})

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":  err.Error(),
				"result": result,
			},
		)

		return
	}

	analyticsEngine.ClearCache()

	status := "success"
	if len(result.Errors) > 0 || result.Skipped > 0 {
		status = "warning"
	}

	_ = services.CreateImportLog(services.ImportLogInput{
		Source:   "NHL API",
		Season:   result.DisplaySeason,
		Status:   status,
		Message:  "season synced",
		Imported: result.PlayersCreated + result.StatsCreated,
		Updated:  result.PlayersUpdated + result.StatsUpdated,
		Skipped:  result.Skipped,
		Errors:   result.Errors,
	})

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "season synced",
			"result":  result,
		},
	)
}
