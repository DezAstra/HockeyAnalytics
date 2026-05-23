package handlers

import (
	"net/http"
	"strconv"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"

	"github.com/gin-gonic/gin"
)

type PlayerComparisonData struct {
	Player   string `json:"player"`
	Team     string `json:"team"`
	Position string `json:"position"`
	Season   string `json:"season"`

	Goals   int `json:"goals"`
	Assists int `json:"assists"`
	Points  int `json:"points"`

	Hits   int `json:"hits"`
	Blocks int `json:"blocks"`

	BaseScore       float64 `json:"base_score"`
	NormalizedScore float64 `json:"normalized_score"`
	ContextScore    float64 `json:"context_score"`
	OverallScore    float64 `json:"overall_score"`
}

type ComparisonResponse struct {
	Player1 PlayerComparisonData `json:"player1"`

	Player2 PlayerComparisonData `json:"player2"`
}

func buildComparisonData(
	player models.Player,
	stats models.PlayerSeasonStats,
) PlayerComparisonData {

	base :=
		analytics.BaseStatModel(stats)

	normalized :=
		analytics.NormalizedModel(stats)

	context :=
		analytics.ContextModel(
			stats,
			player.Position,
		)

	overall := analytics.CalculateOverallScore(
		normalized,
		analytics.CalculateDistribution([]float64{normalized}),
		context,
		analytics.CalculateDistribution([]float64{context}),
	)

	return PlayerComparisonData{
		Player: player.Name,

		Team: stats.Team,

		Position: player.Position,

		Season: stats.Season,

		Goals: stats.Goals,

		Assists: stats.Assists,

		Points: stats.Goals + stats.Assists,

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

	player1ID, err1 :=
		strconv.Atoi(player1Query)

	player2ID, err2 :=
		strconv.Atoi(player2Query)

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

	database.DB.
		Preload("Stats").
		First(&player1, player1ID)

	database.DB.
		Preload("Stats").
		First(&player2, player2ID)

	var player1Stats models.PlayerSeasonStats
	var player2Stats models.PlayerSeasonStats

	// latest season fallback

	if season == "" {

		player1Stats =
			player1.Stats[len(player1.Stats)-1]

		player2Stats =
			player2.Stats[len(player2.Stats)-1]

	} else {

		for _, stats := range player1.Stats {

			if stats.Season == season {

				player1Stats = stats
			}
		}

		for _, stats := range player2.Stats {

			if stats.Season == season {

				player2Stats = stats
			}
		}
	}

	response := ComparisonResponse{
		Player1: buildComparisonData(
			player1,
			player1Stats,
		),

		Player2: buildComparisonData(
			player2,
			player2Stats,
		),
	}

	c.JSON(http.StatusOK, response)
}
