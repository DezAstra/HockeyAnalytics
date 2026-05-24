package handlers

import (
	"net/http"
	"sort"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/services"
	"hockeyAnalytics/internal/utils"

	"github.com/gin-gonic/gin"
)

type TeamListItemResponse struct {
	Team              string  `json:"team"`
	Season            string  `json:"season"`
	PlayersCount      int     `json:"players_count"`
	AverageOverall    float64 `json:"average_overall"`
	AveragePercentile float64 `json:"average_percentile"`
	BestPlayerID      uint    `json:"best_player_id"`
	BestPlayerNHLID   int     `json:"best_player_nhl_id"`
	BestPlayer        string  `json:"best_player"`
	BestOverall       float64 `json:"best_overall"`
	Goals             int     `json:"goals"`
	Assists           int     `json:"assists"`
	Points            int     `json:"points"`
	TopArchetype      string  `json:"top_archetype"`
}

func GetTeams(c *gin.Context) {
	season := c.Query("season")
	if season == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "season is required"})
		return
	}

	displaySeason, err := utils.ToDisplaySeason(season)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.SyncSeasonIfMissing(displaySeason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	seasonAnalytics, err := analyticsEngine.GetSeasonAnalytics(c.Request.Context(), displaySeason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	allOverallScores := make([]float64, 0, len(seasonAnalytics))
	for _, result := range seasonAnalytics {
		allOverallScores = append(allOverallScores, result.Overall)
	}

	var stats []models.PlayerSeasonStats
	err = database.DB.
		Preload("Player").
		Where("season = ?", displaySeason).
		Find(&stats).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	teams := make(map[string]*TeamListItemResponse)
	archetypes := make(map[string]map[string]int)

	for _, item := range stats {
		if item.Team == "" {
			continue
		}

		result, ok := seasonAnalytics[item.Player.ID]
		if !ok {
			continue
		}

		team, ok := teams[item.Team]
		if !ok {
			team = &TeamListItemResponse{Team: item.Team, Season: displaySeason}
			teams[item.Team] = team
			archetypes[item.Team] = make(map[string]int)
		}

		percentile := analytics.CalculatePercentile(result.Overall, allOverallScores)
		archetype := analytics.DetectArchetype(item, item.Player.Position)

		team.PlayersCount++
		team.AverageOverall += result.Overall
		team.AveragePercentile += percentile
		team.Goals += item.Goals
		team.Assists += item.Assists
		team.Points += item.Points
		archetypes[item.Team][archetype]++

		if result.Overall > team.BestOverall || team.BestPlayerID == 0 {
			team.BestOverall = result.Overall
			team.BestPlayerID = item.Player.ID
			team.BestPlayerNHLID = item.Player.NHLID
			team.BestPlayer = item.Player.Name
		}
	}

	response := make([]TeamListItemResponse, 0, len(teams))
	for code, team := range teams {
		if team.PlayersCount > 0 {
			team.AverageOverall = team.AverageOverall / float64(team.PlayersCount)
			team.AveragePercentile = team.AveragePercentile / float64(team.PlayersCount)
		}

		maxCount := 0
		for archetype, count := range archetypes[code] {
			if count > maxCount {
				team.TopArchetype = archetype
				maxCount = count
			}
		}

		response = append(response, *team)
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].AverageOverall > response[j].AverageOverall
	})

	c.JSON(http.StatusOK, response)
}
