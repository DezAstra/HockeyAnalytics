package handlers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"hockeyAnalytics/internal/database"

	"github.com/gin-gonic/gin"
)

const firstSupportedSeasonStart = 2010

func GetSeasons(c *gin.Context) {
	seasonSet := map[string]bool{}

	var dbSeasons []string

	if err := database.DB.
		Table("player_season_stats").
		Distinct().
		Pluck("season", &dbSeasons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, season := range dbSeasons {
		season = strings.TrimSpace(season)
		if season != "" {
			seasonSet[season] = true
		}
	}

	for _, season := range supportedNHLSeasons() {
		seasonSet[season] = true
	}

	seasons := make([]string, 0, len(seasonSet))
	for season := range seasonSet {
		seasons = append(seasons, season)
	}

	sort.Slice(seasons, func(i, j int) bool {
		return seasonValue(seasons[i]) > seasonValue(seasons[j])
	})

	c.JSON(http.StatusOK, seasons)
}

func supportedNHLSeasons() []string {
	now := time.Now()
	currentStartYear := now.Year()

	// NHL season starts in autumn. Until then the active/latest completed
	// season is the one that started in the previous calendar year.
	if now.Month() < time.September {
		currentStartYear--
	}

	seasons := make([]string, 0, currentStartYear-firstSupportedSeasonStart+1)

	for year := currentStartYear; year >= firstSupportedSeasonStart; year-- {
		seasons = append(seasons, formatShortSeason(year))
	}

	return seasons
}

func formatShortSeason(startYear int) string {
	return twoDigitYear(startYear) + "/" + twoDigitYear(startYear+1)
}

func twoDigitYear(year int) string {
	return strconv.Itoa(year%100 + 100)[1:]
}

func seasonValue(season string) int {
	parts := strings.Split(season, "/")
	if len(parts) != 2 {
		if len(season) == 8 {
			value, _ := strconv.Atoi(season[:4])
			return value
		}
		return 0
	}

	year, _ := strconv.Atoi(parts[0])
	if year >= 90 {
		return 1900 + year
	}

	return 2000 + year
}
