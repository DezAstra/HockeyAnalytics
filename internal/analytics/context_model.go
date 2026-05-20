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

	// голы

	score += float64(stats.Goals) * 1.0

	// ассисты

	score += float64(stats.Assists) * 0.8

	// броски

	score += float64(stats.Shots) * 0.1

	// коэффициенты усиления голов

	score += float64(stats.EvenStrengthGoals) * 0.3

	score -= float64(stats.PowerPlayGoals) * 0.1

	score += float64(stats.ShortHandedGoals) * 0.7

	// силовая и блоки

	score += float64(stats.BlockedShots) * 0.05

	score += float64(stats.Hits) * 0.03

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

	// эффективность бросков

	if stats.Shots > 50 {

		shootingPercent :=
			(float64(stats.Goals) /
				float64(stats.Shots)) * 100

		if shootingPercent >= 18 {

			score += 3
		}

		if shootingPercent >= 14 &&
			shootingPercent < 18 {

			score += 2
		}

		if shootingPercent <= 7 {

			score -= 2
		}
	}

	// защитники

	if position == "D" {

		// голы менее значимы

		score -= float64(stats.Goals) * 0.2

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

		// если есть raw данные

		if totalFaceoffs > 100 {

			faceoffPercent :=
				CalculateFaceoffPercent(
					stats.FaceoffsWon,
					stats.FaceoffsLost,
				)

			if faceoffPercent >= 60 {

				score += 8
			}

			if faceoffPercent >= 50 &&
				faceoffPercent < 60 {

				score += 4
			}

			if faceoffPercent < 45 &&
				faceoffPercent >= 40 {

				score -= 4
			}

			if faceoffPercent < 40 {

				score -= 8
			}

			// бонус за объём

			if totalFaceoffs >= 800 {

				score += 2
			}

			if totalFaceoffs >= 500 &&
				totalFaceoffs < 800 {

				score += 1
			}

		} else {

			// fallback для NHL API

			if stats.FaceoffPercent != nil &&
				*stats.FaceoffPercent >= 60 {

				score += 8
			}

			if stats.FaceoffPercent != nil &&
				*stats.FaceoffPercent >= 50 &&
				*stats.FaceoffPercent < 60 {

				score += 4
			}

			if stats.FaceoffPercent != nil &&
				*stats.FaceoffPercent < 45 &&
				*stats.FaceoffPercent >= 40 {

				score -= 4
			}

			if stats.FaceoffPercent != nil &&
				*stats.FaceoffPercent < 40 {

				score -= 8
			}
		}
	}

	// вингеры

	if position == "LW" ||
		position == "RW" ||
		position == "L" ||
		position == "R" {

		score += float64(stats.Goals) * 1.0

		score += float64(stats.Shots) * 0.08

		score += float64(stats.Hits) * 0.1
	}

	return score
}
