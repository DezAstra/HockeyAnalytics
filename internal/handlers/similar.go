package handlers

import (
	"math"
	"net/http"
	"sort"
	"strconv"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"

	"github.com/gin-gonic/gin"
)

type SimilarPlayer struct {
	Player string `json:"player"`

	PlayerID uint `json:"player_id"`

	Team string `json:"team"`

	Position string `json:"position"`

	Similarity float64 `json:"similarity"`

	OverallScore float64 `json:"overall_score"`

	Archetype string `json:"archetype"`
}

func GetSimilarPlayers(
	c *gin.Context,
) {

	idParam :=
		c.Param("id")

	playerID, err :=
		strconv.Atoi(idParam)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid id",
			},
		)

		return
	}

	var target models.Player

	database.DB.
		Preload("Stats").
		First(&target, playerID)

	if len(target.Stats) == 0 {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "player stats not found",
			},
		)

		return
	}

	targetStats :=
		target.Stats[len(target.Stats)-1]

	var players []models.Player

	database.DB.
		Preload("Stats").
		Find(&players)

	var results []SimilarPlayer

	targetOverall :=
		analytics.NormalizedModel(
			targetStats,
		)*0.3 +
			analytics.ContextModel(
				targetStats,
				target.Position,
			)*0.4

	for _, player := range players {

		if player.ID == target.ID {

			continue
		}

		if len(player.Stats) == 0 {

			continue
		}

		stats := player.Stats[len(player.Stats)-1]

		overall :=
			analytics.NormalizedModel(
				stats,
			)*0.3 +
				analytics.ContextModel(
					stats,
					player.Position,
				)*0.4

		distance :=
			calculateDistance(
				targetStats,
				stats,
				targetOverall,
				overall,
			)

		similarity :=
			math.Max(
				0,
				100-distance,
			)

		results = append(
			results,
			SimilarPlayer{

				Player: player.Name,

				PlayerID: player.ID,

				Team: stats.Team,

				Position: player.Position,

				Similarity: similarity,

				OverallScore: overall,

				Archetype: analytics.DetectArchetype(
					stats,
					player.Position,
				),
			},
		)
	}

	sort.Slice(
		results,
		func(i, j int) bool {

			return results[i].Similarity >
				results[j].Similarity
		},
	)

	if len(results) > 10 {

		results = results[:10]
	}

	c.JSON(
		http.StatusOK,
		results,
	)
}

func calculateDistance(
	a models.PlayerSeasonStats,
	b models.PlayerSeasonStats,
	overallA float64,
	overallB float64,
) float64 {

	diffGoals :=
		math.Abs(
			float64(
				a.Goals - b.Goals,
			),
		)

	diffAssists :=
		math.Abs(
			float64(
				a.Assists - b.Assists,
			),
		)

	diffShots :=
		math.Abs(
			float64(
				a.Shots - b.Shots,
			),
		)

	diffHits :=
		math.Abs(
			float64(
				a.Hits - b.Hits,
			),
		)

	diffBlocks :=
		math.Abs(
			float64(
				a.BlockedShots -
					b.BlockedShots,
			),
		)

	diffOverall :=
		math.Abs(
			overallA - overallB,
		)

	return (diffGoals*0.8 + diffAssists*0.7 + diffShots*0.2 + diffHits*0.15 + diffBlocks*0.15 + diffOverall*1.5) / 10
}
