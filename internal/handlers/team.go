package handlers

import (
	"net/http"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"

	"github.com/gin-gonic/gin"
)

type TeamPlayerResponse struct {
	PlayerID uint `json:"player_id"`

	Player string `json:"player"`

	Position string `json:"position"`

	Team string `json:"team"`

	Goals int `json:"goals"`

	Assists int `json:"assists"`

	Points int `json:"points"`

	NormalizedScore float64 `json:"normalized_score"`

	ContextScore float64 `json:"context_score"`

	OverallScore float64 `json:"overall_score"`

	Archetype string `json:"archetype"`
}

type TeamResponse struct {
	Team string `json:"team"`

	Season string `json:"season"`

	Players []TeamPlayerResponse `json:"players"`
}

func GetTeam(c *gin.Context) {

	team :=
		c.Param("team")

	season :=
		c.Query("season")

	var stats []models.PlayerSeasonStats

	query :=
		database.DB.
			Preload("Player").
			Where("team = ?", team)

	if season != "" {

		query =
			query.Where(
				"season = ?",
				season,
			)
	}

	query.Find(&stats)

	var response []TeamPlayerResponse

	for _, s := range stats {

		normalized :=
			analytics.NormalizedModel(s)

		context :=
			analytics.ContextModel(
				s,
				s.Player.Position,
			)

		overall :=
			normalized*0.3 +
				context*0.4

		archetype :=
			analytics.DetectArchetype(
				s,
				s.Player.Position,
			)

		response = append(
			response,
			TeamPlayerResponse{
				PlayerID: s.Player.ID,

				Player: s.Player.Name,

				Position: s.Player.Position,

				Team: s.Team,

				Goals: s.Goals,

				Assists: s.Assists,

				Points: s.Points,

				NormalizedScore: normalized,

				ContextScore: context,

				OverallScore: overall,

				Archetype: archetype,
			},
		)
	}

	c.JSON(
		http.StatusOK,
		TeamResponse{
			Team: team,

			Season: season,

			Players: response,
		},
	)
}
