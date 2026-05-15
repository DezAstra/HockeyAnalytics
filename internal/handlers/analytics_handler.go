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

	BaseScore       float64 `json:"base_score"`
	NormalizedScore float64 `json:"normalized_score"`
	ContextScore    float64 `json:"context_score"`
}

func GetPlayersAnalytics(c *gin.Context) {

	var players []models.Player

	database.DB.
		Preload("Team").
		Preload("Stats").
		Find(&players)

	var response []PlayerAnalyticsResponse

	for _, player := range players {

		base :=
			analytics.BaseStatModel(player.Stats)

		normalized :=
			analytics.NormalizedModel(player.Stats)

		context :=
			analytics.ContextModel(
				player.Stats,
				player.Position,
			)

		response = append(
			response,
			PlayerAnalyticsResponse{
				Player:   player.Name,
				Team:     player.Team.Name,
				Position: player.Position,

				BaseScore:       base,
				NormalizedScore: normalized,
				ContextScore:    context,
			},
		)
	}

	c.JSON(http.StatusOK, response)
}
