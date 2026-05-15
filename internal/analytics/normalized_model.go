package analytics

import "hockeyAnalytics/internal/models"

func NormalizedModel(
	stats models.PlayerSeasonStats,
) float64 {

	base := BaseStatModel(stats)

	if stats.TimeOfIce != nil &&
		*stats.TimeOfIce > 0 {

		return (base / *stats.TimeOfIce) * 1000
	}

	if stats.GamesPlayed > 0 {

		return (base / float64(stats.GamesPlayed)) * 55
	}

	return 0
}
