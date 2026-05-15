package handlers

import (
	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
	"net/http"
	"sort"
	"strconv"

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
	OverallScore    float64 `json:"overall_score"`
}

func GetPlayersAnalytics(c *gin.Context) {

	var players []models.Player

	season := c.Query("season")
	position := c.Query("position")
	team := c.Query("team")
	sortBy := c.DefaultQuery("sort", "overall")
	limit := c.DefaultQuery("limit", "50")

	query :=
		database.DB.
			Preload("Team").
			Preload("Stats")

	if position != "" {
		query = query.Where(
			"position = ?",
			position,
		)
	}
	if team != "" {

		query = query.Joins(
			"JOIN teams ON teams.id = players.team_id",
		).Where(
			"teams.name = ?",
			team,
		)
	}

	query.Find(&players)

	var response []PlayerAnalyticsResponse

	for _, player := range players {

		for _, stats := range player.Stats {

			if season != "" &&
				stats.Season != season {

				continue
			}

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
				PlayerAnalyticsResponse{
					Player:          player.Name,
					Team:            player.Team.Name,
					Position:        player.Position,
					Season:          stats.Season,
					BaseScore:       base,
					NormalizedScore: normalized,
					ContextScore:    context,
					OverallScore:    overall,
				},
			)
		}
	}

	switch sortBy {

	case "base":

		sort.Slice(response,
			func(i, j int) bool {

				return response[i].BaseScore >
					response[j].BaseScore
			})

	case "normalized":

		sort.Slice(response,
			func(i, j int) bool {

				return response[i].NormalizedScore >
					response[j].NormalizedScore
			})

	case "context":

		sort.Slice(response,
			func(i, j int) bool {

				return response[i].ContextScore >
					response[j].ContextScore
			})

	default:

		sort.Slice(response,
			func(i, j int) bool {

				return response[i].OverallScore >
					response[j].OverallScore
			})
	}

	parsedLimit, err :=
		strconv.Atoi(limit)

	if err == nil &&
		parsedLimit < len(response) {

		response = response[:parsedLimit]
	}

	c.JSON(http.StatusOK, response)
}
