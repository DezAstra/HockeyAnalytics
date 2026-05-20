package handlers

import (
	"net/http"

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

	OverallScore float64 `json:"overall_score"`
	Archetype    string  `json:"archetype"`
}

type CareerResponse struct {
	Player string `json:"player"`

	Position string `json:"position"`

	Career []CareerSeasonResponse `json:"career"`
}

// GetPlayerCareer godoc
// @Summary Карьера игрока
// @Description Возвращает историю сезонов игрока
// @Tags Игроки
// @Produce json
// @Param id path int true "Player ID"
// @Success 200 {object} CareerResponse
// @Router /players/{id}/career [get]
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

		overall :=
			analytics.NormalizedModel(stats)*0.3 +
				analytics.ContextModel(
					stats,
					player.Position,
				)*0.4

		archetype :=
			analytics.DetectArchetype(
				stats,
				player.Position,
			)

		seasons = append(
			seasons,
			CareerSeasonResponse{
				Season:      stats.Season,
				Team:        stats.Team,
				GamesPlayed: stats.GamesPlayed,

				Goals:   stats.Goals,
				Assists: stats.Assists,
				Points:  stats.Points,

				OverallScore: overall,
				Archetype:    archetype,
			},
		)
	}

	c.JSON(
		http.StatusOK,
		CareerResponse{
			Player: player.Name,

			Position: player.Position,

			Career: seasons,
		},
	)
}
