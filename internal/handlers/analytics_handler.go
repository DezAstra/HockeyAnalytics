package handlers

import (
	"net/http"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/services"

	"github.com/gin-gonic/gin"
)

type PlayerAnalyticsResponse struct {
	PlayerID uint   `json:"player_id"`
	Player   string `json:"player"`
	Team     string `json:"team"`
	Position string `json:"position"`
	Season   string `json:"season"`

	BaseScore         float64 `json:"base_score"`
	NormalizedScore   float64 `json:"normalized_score"`
	ContextScore      float64 `json:"context_score"`
	OverallScore      float64 `json:"overall_score"`
	OverallPercentile float64 `json:"overall_percentile"`
	Archetype         string  `json:"archetype"`
}

// GetPlayersAnalytics godoc
// @Summary Проанализировать игроков
// @Description Возвращает баллы моделей по статистике игроков
// @Tags Аналитика
// @Produce json
// @Success 200 {array} PlayerAnalyticsResponse
// @Router /analytics [get]
func GetPlayersAnalytics(c *gin.Context) {

	var players []models.Player

	season :=
		c.Query("season")

	err := services.SyncSeasonIfMissing(
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

	query := database.DB

	if season != "" {

		query = query.Preload(
			"Stats",
			"season = ?",
			season,
		)

	} else {

		query = query.Preload("Stats")
	}

	query.Find(&players)

	var tempData []analyticsTemp

	var allOverallScores []float64

	// 1 проход

	for _, player := range players {

		if len(player.Stats) == 0 {
			continue
		}

		for _, stats := range player.Stats {

			base :=
				analytics.BaseStatModel(stats)

			normalized :=
				analytics.NormalizedModel(stats)

			context :=
				analytics.ContextModel(
					stats,
					player.Position,
				)

			overall :=
				normalized*0.3 +
					context*0.4

			tempData = append(
				tempData,
				analyticsTemp{
					Player: player,

					Stats: stats,

					Base: base,

					Normalized: normalized,

					Context: context,

					Overall: overall,
				},
			)

			allOverallScores = append(
				allOverallScores,
				overall,
			)
		}
	}

	// 2 проход

	var response []PlayerAnalyticsResponse

	for _, item := range tempData {

		percentile :=
			analytics.CalculatePercentile(
				item.Overall,
				allOverallScores,
			)

		archetype :=
			analytics.DetectArchetype(
				item.Stats,
				item.Player.Position,
			)

		response = append(
			response,
			PlayerAnalyticsResponse{
				PlayerID: item.Player.ID,
				Player:   item.Player.Name,

				Team:     item.Stats.Team,
				Position: item.Player.Position,
				Season:   item.Stats.Season,

				BaseScore:         item.Base,
				NormalizedScore:   item.Normalized,
				ContextScore:      item.Context,
				OverallScore:      item.Overall,
				OverallPercentile: percentile,
				Archetype:         archetype,
			},
		)
	}

	c.JSON(http.StatusOK, response)
}

type analyticsTemp struct {
	Player models.Player
	Stats  models.PlayerSeasonStats

	Base       float64
	Normalized float64
	Context    float64
	Overall    float64
}
