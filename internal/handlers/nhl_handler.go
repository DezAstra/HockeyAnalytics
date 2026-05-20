package handlers

import (
	"hockeyAnalytics/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNHLSeasonStats(
	c *gin.Context,
) {

	season :=
		c.DefaultQuery(
			"season",
			"20232024",
		)

	data, err :=
		services.FetchSeasonSummary(
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
		data,
	)
}
