package analytics

import (
	"hockeyAnalytics/internal/models"
	"math"
	"strings"
)

func DetectArchetype(
	stats models.PlayerSeasonStats,
	position string,
) string {

	position =
		strings.ToUpper(position)

	// sniper

	if stats.Goals >= 35 &&
		stats.Shots >= 220 &&
		stats.Goals > stats.Assists {

		return "Снайпер"
	}

	// playmaker

	if stats.Assists+stats.Goals >= 50 &&
		math.Abs(float64(stats.Assists-stats.Goals)) <= 20 {

		return "Бомбардир"
	}

	// playmaker

	if stats.Assists >= 50 &&
		stats.Assists-stats.Goals >= 30 {

		return "Ассистент"
	}

	// offensive defenseman

	if position == "D" &&
		stats.Assists+stats.Goals >= 40 {

		return "Атакующий защитник"
	}

	// iron defenseman

	if position == "D" &&
		stats.BlockedShots >= 120 {

		return "Защитник-стена"
	}

	// enforcer

	if stats.PenaltyMinutes >= 60 {

		return "Нарушитель"
	}

	// grinder

	if stats.Hits >= 120 &&
		stats.Goals < 20 {

		return "Силовик"
	}

	// faceoff specialist

	if position == "C" &&
		stats.FaceoffPercent != nil &&
		*stats.FaceoffPercent >= 57 {

		return "Специалист по вбрасываниям"
	}

	return "Баланс"
}
