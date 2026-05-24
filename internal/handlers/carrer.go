package handlers

import (
	"net/http"
	"sort"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"

	"github.com/gin-gonic/gin"
)

type CareerSeasonResponse struct {
	Season string `json:"season"`
	Team   string `json:"team"`

	GamesPlayed int `json:"games_played"`
	Goals       int `json:"goals"`
	Assists     int `json:"assists"`
	Points      int `json:"points"`

	BaseScore         float64 `json:"base_score"`
	NormalizedScore   float64 `json:"normalized_score"`
	ContextScore      float64 `json:"context_score"`
	OverallScore      float64 `json:"overall_score"`
	OverallPercentile float64 `json:"overall_percentile"`
	Archetype         string  `json:"archetype"`
}

type CareerResponse struct {
	PlayerID uint `json:"player_id"`

	NHLID int `json:"nhl_id"`

	Player string `json:"player"`

	Position string `json:"position"`

	Career []CareerSeasonResponse `json:"career"`
}

func GetPlayerCareer(c *gin.Context) {

	id := c.Param("id")

	var player models.Player

	err :=
		database.DB.
			Preload("Stats").
			First(&player, id).Error

	if err != nil {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "игрок не найден",
			},
		)

		return
	}

	var seasons []CareerSeasonResponse

	for _, stats := range player.Stats {

		seasonAnalytics,
			err :=
			analyticsEngine.GetSeasonAnalytics(
				c.Request.Context(),
				stats.Season,
			)

		if err != nil {
			continue
		}

		result,
			ok :=
			seasonAnalytics[player.ID]

		if !ok {
			continue
		}

		allOverallScores := make([]float64, 0, len(seasonAnalytics))
		for _, analyticsResult := range seasonAnalytics {
			allOverallScores = append(allOverallScores, analyticsResult.Overall)
		}

		percentile := analytics.CalculatePercentile(result.Overall, allOverallScores)

		archetype :=
			analytics.DetectArchetype(
				stats,
				player.Position,
			)

		seasons =
			append(
				seasons,
				CareerSeasonResponse{
					Season: stats.Season,

					Team: stats.Team,

					GamesPlayed: stats.GamesPlayed,

					Goals: stats.Goals,

					Assists: stats.Assists,

					Points: stats.Points,

					BaseScore: result.BaseScore,

					NormalizedScore: result.NormalizedScore,

					ContextScore: result.ContextScore,

					OverallScore: result.Overall,

					OverallPercentile: percentile,

					Archetype: archetype,
				},
			)
	}

	sort.Slice(
		seasons,
		func(i, j int) bool {

			return seasonValue(
				seasons[i].Season,
			) <
				seasonValue(
					seasons[j].Season,
				)
		},
	)

	c.JSON(
		http.StatusOK,
		CareerResponse{
			PlayerID: player.ID,

			NHLID: player.NHLID,

			Player: player.Name,

			Position: player.Position,

			Career: seasons,
		},
	)
}
