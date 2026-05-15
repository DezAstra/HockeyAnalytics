package analytics

import (
	"strings"

	"hockeyAnalytics/internal/models"
)

func ContextModel(
	stats models.PlayerSeasonStats,
	position string,
) float64 {

	score := 0.0

	position = strings.ToUpper(position)

	// универсальные метрики
	// ассисты (г/п)

	score += float64(stats.Assists) * 0.8

	// броски в створ ворот

	score += float64(stats.Shots) * 0.10

	// коэффициенты усиления голов

	score += float64(stats.EvenStrengthGoals) * 0.3

	score -= float64(stats.PowerPlayGoals) * 0.1

	score += float64(stats.ShortHandedGoals) * 0.7

	// нарушения

	if stats.PenaltyMinutes > 20 {

		excessPIM :=
			stats.PenaltyMinutes - 20

		score -= float64(excessPIM) * 0.15
	}

	// +/-

	if stats.GamesPlayed > 20 {

		score += float64(stats.PlusMinus) * 0.25
	}

	// защитники

	if position == "D" {

		score += float64(stats.Goals) * 0.6

		score += float64(stats.BlockedShots) * 0.5
		score += float64(stats.Hits) * 0.15

		score += float64(stats.Assists) * 0.3
	}

	// центры

	if position == "C" {

		score += float64(stats.Goals) * 0.9

		score += float64(stats.Assists) * 0.2

		totalFaceoffs :=
			stats.FaceoffsWon +
				stats.FaceoffsLost

		if totalFaceoffs > 100 {

			faceoffPercent :=
				CalculateFaceoffPercent(
					stats.FaceoffsWon,
					stats.FaceoffsLost,
				)

			// высокий

			if faceoffPercent >= 60 {

				score += 8
			}

			// хороший

			if faceoffPercent >= 50 &&
				faceoffPercent < 60 {

				score += 4
			}

			// умеренный

			if faceoffPercent < 45 &&
				faceoffPercent >= 40 {

				score -= 4
			}

			// низкий

			if faceoffPercent < 40 {

				score -= 8
			}
		}
	}

	// вингеры (крайние напы)

	if position == "LW" ||
		position == "RW" {

		score += float64(stats.Goals) * 1.0

		score += float64(stats.Shots) * 0.08

		score += float64(stats.Hits) * 0.1
	}

	return score
}
