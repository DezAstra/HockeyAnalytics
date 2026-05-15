package analytics

import "hockeyAnalytics/internal/models"

func BaseStatModel(stats models.PlayerStats) float64 {

	score :=
		float64(stats.Goals*2) +
			float64(stats.Assists) +
			float64(stats.PlusMinus)

	score -= float64(stats.PenaltyMinutes) * 0.2

	return score
}
