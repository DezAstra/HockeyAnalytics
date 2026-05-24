package handlers

import (
	"hockeyAnalytics/internal/services"
	"hockeyAnalytics/internal/utils"
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

	apiSeason, err := utils.ToAPISeason(season)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	data, err :=
		services.FetchSeasonSummary(
			apiSeason,
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
