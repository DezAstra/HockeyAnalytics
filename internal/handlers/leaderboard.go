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

type LeaderboardResponse struct {
	Player   string `json:"player"`
	Team     string `json:"team"`
	Position string `json:"position"`
	Season   string `json:"season"`

	Value float64 `json:"value"`
}

// GetLeaderboard godoc
// @Summary Таблица лидеров
// @Description Возвращает список лучших игроков по выбранному показателю
// @Tags Аналитика
// @Produce json
//
// @Param stat path string true "Статистические показатели"
// @Param season query string false "Сезон"
// @Param position query string false "Позиция игрока"
// @Param team query string false "Команда"
// @Param limit query int false "Количество игроков"
// @Param min_gp query int false "Минимум сыгранных матчей"
// @Param order query string false "Сортировка asc/desc"
//
// @Success 200 {array} LeaderboardResponse
//
// @Failure 400 {object} map[string]interface{}
//
// @Router /analytics/leaders/{stat} [get]
func GetLeaderboard(c *gin.Context) {

	stat := c.Param("stat")
	season := c.Query("season")
	position := c.Query("position")
	team := c.Query("team")
	order := c.DefaultQuery("order", "desc")
	limitQuery := c.DefaultQuery("limit", "20")
	minGPQuery := c.DefaultQuery("min_gp", "0")
	limit, _ := strconv.Atoi(limitQuery)
	minGP, _ := strconv.Atoi(minGPQuery)

	var players []models.Player

	query :=
		database.DB.
			Preload("Stats")

	// фильтр по позиции

	if position != "" {

		query = query.Where(
			"position = ?",
			position,
		)
	}

	query.Find(&players)

	var response []LeaderboardResponse

	for _, player := range players {

		for _, stats := range player.Stats {

			// фильтр по сезону

			if season != "" &&
				stats.Season != season {
				continue
			}

			if team != "" &&
				stats.Team != team {
				continue
			}

			// фильтр по мин. матчам

			if stats.GamesPlayed < minGP {
				continue
			}

			value := 0.0

			switch stat {

			// простая статистика

			case "goals":
				value =
					float64(stats.Goals)

			case "assists":
				value =
					float64(stats.Assists)

			case "shots":
				value =
					float64(stats.Shots)

			case "blocks":
				value =
					float64(stats.BlockedShots)

			case "hits":
				value =
					float64(stats.Hits)

			case "pim":
				value =
					float64(stats.PenaltyMinutes)

			case "plusminus":
				value =
					float64(stats.PlusMinus)

			// продвинутые статы

			case "faceoff_percent":
				total :=
					stats.FaceoffsWon +
						stats.FaceoffsLost

				if total < 100 {
					continue
				}

				value =
					analytics.CalculateFaceoffPercent(
						stats.FaceoffsWon,
						stats.FaceoffsLost,
					)

			case "normalized":
				value =
					analytics.NormalizedModel(
						stats,
					)

			case "context":
				value =
					analytics.ContextModel(
						stats,
						player.Position,
					)

			case "overall":
				normalized :=
					analytics.NormalizedModel(stats)

				context :=
					analytics.ContextModel(
						stats,
						player.Position,
					)

				value =
					normalized*0.3 +
						context*0.4

			default:

				c.JSON(
					http.StatusBadRequest,
					gin.H{
						"error": "invalid stat",
					},
				)

				return
			}

			response = append(
				response,
				LeaderboardResponse{
					Player:   player.Name,
					Team:     stats.Team,
					Position: player.Position,
					Season:   stats.Season,
					Value:    value,
				},
			)
		}
	}

	sort.Slice(
		response,
		func(i, j int) bool {

			if order == "asc" {
				return response[i].Value <
					response[j].Value
			}

			return response[i].Value >
				response[j].Value
		},
	)

	if limit < len(response) {
		response = response[:limit]
	}

	c.JSON(http.StatusOK, response)
}
