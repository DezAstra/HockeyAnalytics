package analytics

import "hockeyAnalytics/internal/models"

func BaseStatModel(stats models.PlayerSeasonStats) float64 {

	score :=
		float64(stats.Goals)*1.0 +
			float64(stats.Assists)*0.8 +
			float64(stats.EvenStrengthGoals)*1.5 +
			float64(stats.PowerPlayGoals)*1.0 +
			float64(stats.ShortHandedGoals)*2.0 +
			float64(stats.PlusMinus)*0.3

	score -= float64(stats.PenaltyMinutes) * 0.5

	return score
}
