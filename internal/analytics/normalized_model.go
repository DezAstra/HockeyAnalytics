package analytics

import "hockeyAnalytics/internal/models"

func NormalizedModel(stats models.PlayerStats) float64 {

	base :=
		float64(stats.Goals*2) +
			float64(stats.Assists)

	if stats.TimeOfIce == nil || *stats.TimeOfIce == 0 {
		return base
	}

	return (base / *stats.TimeOfIce) * 1000
}
