package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"

	"github.com/gin-gonic/gin"
)

type PlayerHistoryResponse struct {
	Season string `json:"season"`

	BaseScore float64 `json:"base_score"`

	NormalizedScore float64 `json:"normalized_score"`

	ContextScore float64 `json:"context_score"`

	OverallScore float64 `json:"overall_score"`
}

// GetPlayerHistory godoc
// @Summary История игрока
// @Description Возвращает статистику и аналитические оценки игрока по сезонам
// @Tags Аналитика
// @Produce json
//
// @Param id path int true "ID игрока"
//
// @Success 200 {array} PlayerHistoryResponse
//
// @Failure 404 {object} map[string]interface{}
//
// @Router /analytics/player/{id}/history [get]
func GetPlayerHistory(c *gin.Context) {

	playerID :=
		c.Param("id")

	id, err :=
		strconv.Atoi(playerID)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid player id",
			},
		)

		return
	}

	var player models.Player

	result :=
		database.DB.
			Preload("Stats").
			First(&player, id)

	if result.Error != nil {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "player not found",
			},
		)

		return
	}

	var response []PlayerHistoryResponse

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

		overall := normalized*0.3 + context*0.4

		response = append(
			response,
			PlayerHistoryResponse{
				Season: stats.Season,

				BaseScore: base,

				NormalizedScore: normalized,

				ContextScore: context,

				OverallScore: overall,
			},
		)
	}

	sort.Slice(
		response,
		func(i, j int) bool {

			return response[i].Season <
				response[j].Season
		},
	)

	c.JSON(http.StatusOK, response)
}
