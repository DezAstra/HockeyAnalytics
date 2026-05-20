package handlers

import (
	"net/http"

	"hockeyAnalytics/internal/database"

	"github.com/gin-gonic/gin"
)

func GetSeasons(c *gin.Context) {

	var seasons []string

	database.DB.
		Table("player_season_stats").
		Distinct().
		Pluck("season", &seasons)

	c.JSON(
		http.StatusOK,
		seasons,
	)
}
