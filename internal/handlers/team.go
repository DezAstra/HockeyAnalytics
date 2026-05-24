package handlers

import (
	"net/http"
	"sort"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/utils"

	"github.com/gin-gonic/gin"
)

type TeamPlayerResponse struct {
	PlayerID uint `json:"player_id"`

	NHLID int `json:"nhl_id"`

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

	season =
		displaySeason

	var stats []models.PlayerSeasonStats

	err =
		database.DB.
			Preload("Player").
			Where(
				"team = ? AND season = ?",
				team,
				season,
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

	seasonAnalytics, err :=
		analyticsEngine.GetSeasonAnalytics(
			c.Request.Context(),
			season,
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

	response :=
		make(
			[]TeamPlayerResponse,
			0,
			len(stats),
		)

	for _, s := range stats {

		result, ok :=
			seasonAnalytics[s.Player.ID]

		if !ok {
			continue
		}

		archetype :=
			analytics.DetectArchetype(
				s,
				s.Player.Position,
			)

		response =
			append(
				response,
				TeamPlayerResponse{
					PlayerID: s.Player.ID,

					NHLID: s.Player.NHLID,

					Player: s.Player.Name,

					Position: s.Player.Position,

					Team: s.Team,

					Goals: s.Goals,

					Assists: s.Assists,

					Points: s.Points,

					NormalizedScore: result.NormalizedScore,

					ContextScore: result.ContextScore,

					OverallScore: result.Overall,

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

	c.JSON(
		http.StatusOK,
		TeamResponse{
			Team: team,

			Season: season,

			Players: response,
		},
	)
}

type TeamSeasonHistoryItemResponse struct {
	Season            string  `json:"season"`
	PlayersCount      int     `json:"players_count"`
	AverageOverall    float64 `json:"average_overall"`
	AverageNormalized float64 `json:"average_normalized"`
	AverageContext    float64 `json:"average_context"`
	BestPlayerID      uint    `json:"best_player_id"`
	BestPlayerNHLID   int     `json:"best_player_nhl_id"`
	BestPlayer        string  `json:"best_player"`
	BestOverall       float64 `json:"best_overall"`
	Goals             int     `json:"goals"`
	Assists           int     `json:"assists"`
	Points            int     `json:"points"`
}

type TeamHistoryResponse struct {
	Team       string                          `json:"team"`
	History    []TeamSeasonHistoryItemResponse `json:"history"`
	Trend      string                          `json:"trend"`
	TrendDelta float64                         `json:"trend_delta"`
	BestSeason string                          `json:"best_season"`
}

func seasonSortValue(season string) int {
	apiSeason, err := utils.ToAPISeason(season)
	if err != nil || len(apiSeason) < 4 {
		return 0
	}

	value := 0
	for _, ch := range apiSeason[:4] {
		if ch < '0' || ch > '9' {
			return 0
		}
		value = value*10 + int(ch-'0')
	}

	return value
}

func getTeamTrend(history []TeamSeasonHistoryItemResponse) (string, float64) {
	if len(history) < 2 {
		return "not_enough_data", 0
	}

	first := history[len(history)-4]
	last := history[len(history)-1]
	delta := analytics.Round1(last.AverageOverall - first.AverageOverall)

	if delta > 1 {
		return "up", delta
	}

	if delta < -1 {
		return "down", delta
	}

	return "stable", delta
}

func GetTeamHistory(c *gin.Context) {
	team := c.Param("team")

	var seasons []string
	err := database.DB.
		Model(&models.PlayerSeasonStats{}).
		Where("team = ?", team).
		Distinct("season").
		Pluck("season", &seasons).Error

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	sort.Slice(seasons, func(i, j int) bool {
		return seasonSortValue(seasons[i]) < seasonSortValue(seasons[j])
	})

	history := make([]TeamSeasonHistoryItemResponse, 0, len(seasons))

	for _, season := range seasons {
		seasonAnalytics, err := analyticsEngine.GetSeasonAnalytics(
			c.Request.Context(),
			season,
		)

		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
			return
		}

		var stats []models.PlayerSeasonStats
		err = database.DB.
			Preload("Player").
			Where("team = ? AND season = ?", team, season).
			Find(&stats).Error

		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
			return
		}

		item := TeamSeasonHistoryItemResponse{
			Season: season,
		}

		for _, stat := range stats {
			result, ok := seasonAnalytics[stat.Player.ID]
			if !ok {
				continue
			}

			item.PlayersCount++
			item.AverageOverall += result.Overall
			item.AverageNormalized += result.NormalizedScore
			item.AverageContext += result.ContextScore
			item.Goals += stat.Goals
			item.Assists += stat.Assists
			item.Points += stat.Points

			if result.Overall > item.BestOverall || item.BestPlayerID == 0 {
				item.BestOverall = result.Overall
				item.BestPlayerID = stat.Player.ID
				item.BestPlayerNHLID = stat.Player.NHLID
				item.BestPlayer = stat.Player.Name
			}
		}

		if item.PlayersCount == 0 {
			continue
		}

		count := float64(item.PlayersCount)
		item.AverageOverall = analytics.Round1(item.AverageOverall / count)
		item.AverageNormalized = analytics.Round1(item.AverageNormalized / count)
		item.AverageContext = analytics.Round1(item.AverageContext / count)
		item.BestOverall = analytics.Round1(item.BestOverall)

		history = append(history, item)
	}

	bestSeason := ""
	bestAverage := -1.0
	for _, item := range history {
		if item.AverageOverall > bestAverage {
			bestAverage = item.AverageOverall
			bestSeason = item.Season
		}
	}

	trend, trendDelta := getTeamTrend(history)

	c.JSON(
		http.StatusOK,
		TeamHistoryResponse{
			Team:       team,
			History:    history,
			Trend:      trend,
			TrendDelta: trendDelta,
			BestSeason: bestSeason,
		},
	)
}
