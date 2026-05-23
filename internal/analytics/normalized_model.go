package analytics

import "hockeyAnalytics/internal/models"

func NormalizedModel(
	stats models.PlayerSeasonStats,
) float64 {

	base := BaseStatModel(stats)
	M_minutes := 200.0
	M_games := 15.0
	AvgBasePerMinute := 0.03
	AvgBasePerGame := 0.7

	if stats.TimeOfIce != nil &&
		*stats.TimeOfIce > 0 {

		// Формула: (Base + (AvgPerMin * M)) / (TOI + M) * 1000
		return (base + (AvgBasePerMinute*M_minutes)/(*stats.TimeOfIce+M_minutes))
	}

	if stats.GamesPlayed > 0 {

		return (base + (AvgBasePerGame*M_games)/(float64(stats.GamesPlayed)+M_games))
	}

	return 0
}
