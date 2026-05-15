package analytics

import "hockeyAnalytics/internal/models"

func ContextModel(
	stats models.PlayerStats,
	position string,
) float64 {

	if position == "D" {

		return float64(stats.Assists) +
			float64(stats.PlusMinus*2)
	}

	return float64(stats.Goals*2) +
		float64(stats.Assists)
}
