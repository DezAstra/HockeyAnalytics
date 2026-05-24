package handlers

import (
	"net/http"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/services"
	"hockeyAnalytics/internal/utils"

	"github.com/gin-gonic/gin"
)

type PlayerAnalyticsResponse struct {
	PlayerID uint `json:"player_id"`

	NHLID int `json:"nhl_id"`

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

	GamesPlayed int `json:"games_played"`

	Goals int `json:"goals"`

	Assists int `json:"assists"`

	Points int `json:"points"`
}

// GetPlayersAnalytics godoc
// @Summary Проанализировать игроков
// @Description Возвращает баллы моделей по статистике игроков
// @Tags Аналитика
// @Produce json
// @Param season query string false "Сезон"
// @Success 200 {array} PlayerAnalyticsResponse
// @Router /analytics [get]
func GetPlayersAnalytics(c *gin.Context) {

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

	syncResult, err :=
		services.SyncSeasonIfMissingWithResult(
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

	if !syncResult.AlreadyExisted {
		status := "success"
		if len(syncResult.Errors) > 0 || syncResult.Skipped > 0 {
			status = "warning"
		}

		_ = services.CreateImportLog(services.ImportLogInput{
			Source:   "NHL API lazy",
			Season:   syncResult.DisplaySeason,
			Status:   status,
			Message:  "season lazy synced",
			Imported: syncResult.PlayersCreated + syncResult.StatsCreated,
			Updated:  syncResult.PlayersUpdated + syncResult.StatsUpdated,
			Skipped:  syncResult.Skipped,
			Errors:   syncResult.Errors,
		})
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

	var response []PlayerAnalyticsResponse

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
				PlayerAnalyticsResponse{
					PlayerID: item.Player.ID,

					NHLID: item.Player.NHLID,

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

					GamesPlayed: item.GamesPlayed,

					Goals: item.Goals,

					Assists: item.Assists,

					Points: item.Points,
				},
			)
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}
