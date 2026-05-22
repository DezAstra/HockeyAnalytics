package handlers

import (
	"net/http"
	"strconv"
	"strings"

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

func seasonValue(
	season string,
) int {

	parts :=
		strings.Split(
			season,
			"/",
		)

	if len(parts) != 2 {
		return 0
	}

	year, _ :=
		strconv.Atoi(
			parts[0],
		)

	return year
}
