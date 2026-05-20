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

	err :=
		services.ImportSeason(
			season,
		)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "season synced",
		},
	)
}
