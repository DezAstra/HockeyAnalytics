package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/services"
	"hockeyAnalytics/internal/utils"

	"github.com/gin-gonic/gin"
)

type LeaderboardPlayerResponse struct {
	PlayerID uint `json:"player_id"`

	Player string `json:"player"`

	Team string `json:"team"`

	Position string `json:"position"`

	Season string `json:"season"`

	BaseScore float64 `json:"base_score"`

	NormalizedScore float64 `json:"normalized_score"`

	ContextScore float64 `json:"context_score"`

	OverallScore float64 `json:"overall_score"`

	OverallPercentile float64 `json:"overall_percentile"`

	Archetype string `json:"archetype"`
}

// GetLeaderboard godoc
// @Summary Лидерборд игроков
// @Description Возвращает топ игроков по общей аналитической оценке
// @Tags Аналитика
// @Produce json
// @Param season query string true "Сезон"
// @Param limit query int false "Лимит"
// @Success 200 {array} LeaderboardPlayerResponse
// @Failure 400 {object} map[string]interface{}
// @Router /analytics/leaderboard [get]
func GetLeaderboard(c *gin.Context) {

	season :=
		c.Query("season")

	if season == "" {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "season is required",
			},
		)

		return
	}

	limit := 10

	limitQuery :=
		c.Query("limit")

	if limitQuery != "" {

		parsedLimit, err :=
			strconv.Atoi(
				limitQuery,
			)

		if err == nil &&
			parsedLimit > 0 {

			limit =
				parsedLimit
		}
	}

	displaySeason, err :=
		utils.ToDisplaySeason(
			season,
		)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err =
		services.SyncSeasonIfMissing(
			displaySeason,
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

	seasonAnalytics, err :=
		analyticsEngine.GetSeasonAnalytics(
			c.Request.Context(),
			displaySeason,
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

	var stats []models.PlayerSeasonStats

	err =
		database.DB.
			Preload("Player").
			Where(
				"season = ?",
				displaySeason,
			).
			Find(&stats).
			Error

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	allOverallScores :=
		make(
			[]float64,
			0,
			len(seasonAnalytics),
		)

	for _, result := range seasonAnalytics {

		allOverallScores =
			append(
				allOverallScores,
				result.Overall,
			)
	}

	var response []LeaderboardPlayerResponse

	for _, item := range stats {

		result, ok :=
			seasonAnalytics[item.Player.ID]

		if !ok {
			continue
		}

		archetype :=
			analytics.DetectArchetype(
				item,
				item.Player.Position,
			)

		percentile :=
			analytics.CalculatePercentile(
				result.Overall,
				allOverallScores,
			)

		response =
			append(
				response,
				LeaderboardPlayerResponse{
					PlayerID: item.Player.ID,

					Player: item.Player.Name,

					Team: item.Team,

					Position: item.Player.Position,

					Season: item.Season,

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
		response,
		func(i, j int) bool {

			return response[i].OverallScore >
				response[j].OverallScore
		},
	)

	if len(response) > limit {

		response =
			response[:limit]
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}
