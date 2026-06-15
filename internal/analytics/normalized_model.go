package analytics

import "hockeyAnalytics/internal/models"

// NormalizedModel приводит базовую оценку к условному 82-матчевому сезону.

func NormalizedModel(stats models.PlayerSeasonStats) float64 {
	if stats.GamesPlayed <= 0 {
		return 0
	}

	cfg := DefaultConfig
	base := BaseStatModel(stats)
	games := float64(stats.GamesPlayed)
	basePerGame := base / games

	smoothedPerGame :=
		((basePerGame * games) +
			(cfg.NormalizedAvgBasePerGame * cfg.NormalizedSampleGames)) /
			(games + cfg.NormalizedSampleGames)

	return smoothedPerGame * cfg.NormalizedSeasonGames
}
