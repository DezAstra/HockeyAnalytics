package handlers

import (
	"net/http"
	"strconv"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/utils"

	"github.com/gin-gonic/gin"
)

type PlayerComparisonData struct {
	PlayerID uint `json:"player_id"`

	NHLID int `json:"nhl_id"`

	Player string `json:"player"`

	Team string `json:"team"`

	Position string `json:"position"`

	Season string `json:"season"`

	Goals int `json:"goals"`

	Assists int `json:"assists"`

	Points int `json:"points"`

	Hits int `json:"hits"`

	Blocks int `json:"blocks"`

	BaseScore float64 `json:"base_score"`

	NormalizedScore float64 `json:"normalized_score"`

	ContextScore float64 `json:"context_score"`

	OverallScore float64 `json:"overall_score"`
}

type ComparisonResponse struct {
	Player1 PlayerComparisonData `json:"player1"`

	Player2 PlayerComparisonData `json:"player2"`
}

func buildComparisonData(
	player models.Player,
	stats models.PlayerSeasonStats,
	overall float64,
) PlayerComparisonData {

	base :=
		analytics.BaseStatModel(
			stats,
		)

	normalized :=
		analytics.NormalizedModel(
			stats,
		)

	context :=
		analytics.ContextModel(
			stats,
			player.Position,
		)

	return PlayerComparisonData{
		PlayerID: player.ID,

		NHLID: player.NHLID,

		Player: player.Name,

		Team: stats.Team,

		Position: player.Position,

		Season: stats.Season,

		Goals: stats.Goals,

		Assists: stats.Assists,

		Points: stats.Points,

		Hits: stats.Hits,

		Blocks: stats.BlockedShots,

		BaseScore: base,

		NormalizedScore: normalized,

		ContextScore: context,

		OverallScore: overall,
	}
}

// ComparePlayers godoc
// @Summary Сравнение игроков
// @Description Сравнивает двух игроков по статистике и аналитическим моделям
// @Tags Аналитика
// @Produce json
//
// @Param player1 query int true "ID первого игрока"
// @Param player2 query int true "ID второго игрока"
// @Param season query string false "Сезон"
//
// @Success 200 {object} ComparisonResponse
//
// @Failure 400 {object} map[string]interface{}
//
// @Router /analytics/compare [get]
func ComparePlayers(c *gin.Context) {

	player1Query :=
		c.Query("player1")

	player2Query :=
		c.Query("player2")

	season :=
		c.Query("season")

	if season != "" {
		displaySeason, err := utils.ToDisplaySeason(season)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		season = displaySeason
	}

	player1ID, err1 :=
		strconv.Atoi(
			player1Query,
		)

	player2ID, err2 :=
		strconv.Atoi(
			player2Query,
		)

	if err1 != nil ||
		err2 != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid player ids",
			},
		)

		return
	}

	var player1 models.Player
	var player2 models.Player

	err :=
		database.DB.
			Preload("Stats").
			First(
				&player1,
				player1ID,
			).
			Error

	if err != nil {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "player1 not found",
			},
		)

		return
	}

	err =
		database.DB.
			Preload("Stats").
			First(
				&player2,
				player2ID,
			).
			Error

	if err != nil {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "player2 not found",
			},
		)

		return
	}

	if len(player1.Stats) == 0 ||
		len(player2.Stats) == 0 {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "player stats not found",
			},
		)

		return
	}

	var player1Stats models.PlayerSeasonStats
	var player2Stats models.PlayerSeasonStats

	if season == "" {

		player1Stats =
			player1.Stats[0]

		player2Stats =
			player2.Stats[0]

		for _, stats := range player1.Stats {

			if seasonValue(stats.Season) >
				seasonValue(
					player1Stats.Season,
				) {

				player1Stats =
					stats
			}
		}

		for _, stats := range player2.Stats {

			if seasonValue(stats.Season) >
				seasonValue(
					player2Stats.Season,
				) {

				player2Stats =
					stats
			}
		}

	} else {

		player1Found := false

		player2Found := false

		for _, stats := range player1.Stats {

			if stats.Season == season {

				player1Stats =
					stats

				player1Found =
					true
			}
		}

		for _, stats := range player2.Stats {

			if stats.Season == season {

				player2Stats =
					stats

				player2Found =
					true
			}
		}

		if !player1Found ||
			!player2Found {

			c.JSON(
				http.StatusNotFound,
				gin.H{
					"error": "one of players has no stats for selected season",
				},
			)

			return
		}
	}

	seasonAnalytics, err :=
		analyticsEngine.GetSeasonAnalytics(
			c.Request.Context(),
			player1Stats.Season,
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

	player1Overall := 0.0

	player2Overall := 0.0

	if result, ok :=
		seasonAnalytics[player1.ID]; ok {

		player1Overall =
			result.Overall
	}

	if result, ok :=
		seasonAnalytics[player2.ID]; ok {

		player2Overall =
			result.Overall
	}

	response :=
		ComparisonResponse{
			Player1: buildComparisonData(
				player1,
				player1Stats,
				player1Overall,
			),

			Player2: buildComparisonData(
				player2,
				player2Stats,
				player2Overall,
			),
		}

	c.JSON(
		http.StatusOK,
		response,
	)
}
