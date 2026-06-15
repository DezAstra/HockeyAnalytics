package analytics

import "hockeyAnalytics/internal/models"

// BaseStatModel оценивает базовый вклад игрока по сырым событиям сезона.

func BaseStatModel(stats models.PlayerSeasonStats) float64 {
	cfg := DefaultConfig

	score :=
		float64(stats.Goals)*cfg.GoalWeight +
			float64(stats.Assists)*cfg.AssistWeight +
			float64(stats.EvenStrengthGoals)*cfg.EvenStrengthGoalBonus +
			float64(stats.PowerPlayGoals)*cfg.PowerPlayGoalModifier +
			float64(stats.ShortHandedGoals)*cfg.ShortHandedGoalBonus +
			float64(stats.PlusMinus)*cfg.PlusMinusWeight

	if stats.PenaltyMinutes > cfg.PenaltyGraceMinutes {
		excess := stats.PenaltyMinutes - cfg.PenaltyGraceMinutes
		score -= float64(excess) * cfg.PenaltyMinuteWeight
	}

	return score
}
