package analytics

import (
	"math"
	"strings"
)

func CalculateShootingPercent(
	goals int,
	shots int,
) float64 {

	if shots == 0 {
		return 0
	}

	return (float64(goals) / float64(shots)) * 100
}

func CalculateFaceoffPercent(
	won int,
	lost int,
) float64 {

	total := won + lost

	if total == 0 {
		return 0
	}

	return (float64(won) / float64(total)) * 100
}

func Round1(value float64) float64 {
	return math.Round(value*10) / 10
}

func Round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func SafePerGame(value int, games int) float64 {
	if games <= 0 {
		return 0
	}

	return float64(value) / float64(games)
}

func Per82(value int, games int) float64 {
	return SafePerGame(value, games) * DefaultConfig.NormalizedSeasonGames
}

func NormalizePosition(position string) string {
	return strings.ToUpper(strings.TrimSpace(position))
}

func PositionGroup(position string) string {
	switch NormalizePosition(position) {
	case "C":
		return "C"
	case "D":
		return "D"
	case "L", "LW", "R", "RW":
		return "W"
	default:
		return "SKATER"
	}
}

func ConfidenceFactor(gamesPlayed int) float64 {
	if gamesPlayed <= 0 {
		return 0
	}

	confidence := float64(gamesPlayed) / DefaultConfig.ConfidenceFullGames

	if confidence > 1 {
		return 1
	}

	if confidence < 0 {
		return 0
	}

	return confidence
}
