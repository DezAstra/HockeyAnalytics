package handlers

import (
	"net/http"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"

	"github.com/gin-gonic/gin"
)

type PlayerAnalyticsResponse struct {
	Player   string `json:"player"`
	Team     string `json:"team"`
	Position string `json:"position"`
	Season   string `json:"season"`

	BaseScore       float64 `json:"base_score"`
	NormalizedScore float64 `json:"normalized_score"`
	ContextScore    float64 `json:"context_score"`

	OverallScore float64 `json:"overall_score"`

	OverallPercentile float64 `json:"overall_percentile"`
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

	database.DB.
		Preload("Team").
		Preload("Stats").
		Find(&players)

	var tempData []analyticsTemp

	var allOverallScores []float64

	// 1 проход

	for _, player := range players {

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

		response = append(
			response,
			PlayerAnalyticsResponse{
				Player: item.Player.Name,

				Team: item.Player.Team.Name,

				Position: item.Player.Position,

				Season: item.Stats.Season,

				BaseScore: item.Base,

				NormalizedScore: item.Normalized,

				ContextScore: item.Context,

				OverallScore: item.Overall,

				OverallPercentile: percentile,
			},
		)
	}

	c.JSON(http.StatusOK, response)
}

type analyticsTemp struct {
	Player models.Player

	Stats models.PlayerSeasonStats

	Base float64

	Normalized float64

	Context float64

	Overall float64
}
