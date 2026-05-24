package analytics

import (
	"math"

	"hockeyAnalytics/internal/models"
)

type SimilarityInput struct {
	Stats    models.PlayerSeasonStats
	Position string
	Overall  float64
}

func scaledAbsDiff(a float64, b float64, scale float64) float64 {
	if scale <= 0 {
		return 0
	}

	return math.Abs(a-b) / scale
}

// SimilarityDistance считает расстояние между игроками по нормализованным признакам.
// Используются per-82 значения, чтобы не сравнивать напрямую totals игроков с разным числом матчей.
func SimilarityDistance(a SimilarityInput, b SimilarityInput) float64 {
	aStats := a.Stats
	bStats := b.Stats

	distance := 0.0
	distance += scaledAbsDiff(Per82(aStats.Goals, aStats.GamesPlayed), Per82(bStats.Goals, bStats.GamesPlayed), 35) * 1.00
	distance += scaledAbsDiff(Per82(aStats.Assists, aStats.GamesPlayed), Per82(bStats.Assists, bStats.GamesPlayed), 45) * 0.90
	distance += scaledAbsDiff(Per82(aStats.Shots, aStats.GamesPlayed), Per82(bStats.Shots, bStats.GamesPlayed), 220) * 0.45
	distance += scaledAbsDiff(Per82(aStats.Hits, aStats.GamesPlayed), Per82(bStats.Hits, bStats.GamesPlayed), 150) * 0.35
	distance += scaledAbsDiff(Per82(aStats.BlockedShots, aStats.GamesPlayed), Per82(bStats.BlockedShots, bStats.GamesPlayed), 130) * 0.35
	distance += scaledAbsDiff(a.Overall, b.Overall, 20) * 1.25

	aGroup := PositionGroup(a.Position)
	bGroup := PositionGroup(b.Position)

	if aGroup != bGroup {
		distance += 0.55
	} else if NormalizePosition(a.Position) != NormalizePosition(b.Position) {
		distance += 0.15
	}

	return distance
}

func SimilarityPercent(distance float64) float64 {
	similarity := 100 - (distance * 18)

	if similarity < 0 {
		return 0
	}

	if similarity > 100 {
		return 100
	}

	return Round2(similarity)
}
