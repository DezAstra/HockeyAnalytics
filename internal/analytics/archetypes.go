package analytics

import (
	"hockeyAnalytics/internal/models"
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

		return "Sniper"
	}

	// playmaker

	if stats.Assists >= 50 &&
		stats.Assists > stats.Goals {

		return "Playmaker"
	}

	// two-way forward

	if (position == "C" ||
		position == "LW" ||
		position == "RW") &&
		stats.PlusMinus >= 15 &&
		stats.FaceoffPercent != nil &&
		*stats.FaceoffPercent >= 52 {

		return "Two-Way Forward"
	}

	// offensive defenseman

	if position == "D" &&
		stats.Assists >= 40 {

		return "Offensive Defenseman"
	}

	// defensive defenseman

	if position == "D" &&
		stats.BlockedShots >= 120 {

		return "Defensive Defenseman"
	}

	// enforcer

	if stats.Hits >= 180 &&
		stats.PenaltyMinutes >= 70 {

		return "Enforcer"
	}

	// grinder

	if stats.Hits >= 120 &&
		stats.Goals < 20 {

		return "Grinder"
	}

	// faceoff specialist

	if position == "C" &&
		stats.FaceoffPercent != nil &&
		*stats.FaceoffPercent >= 57 {

		return "Faceoff Specialist"
	}

	return "Balanced"
}
