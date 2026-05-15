package analytics

import "hockeyAnalytics/internal/models"

type AnalyticsResult struct {
	BaseScore       float64 `json:"base_score"`
	NormalizedScore float64 `json:"normalized_score"`
	ContextScore    float64 `json:"context_score"`
}

func CalculateAllModels(
	stats models.PlayerStats,
	position string,
) AnalyticsResult {

	return AnalyticsResult{
		BaseScore: BaseStatModel(stats),

		NormalizedScore: NormalizedModel(stats),

		ContextScore: ContextModel(stats, position),
	}
}
