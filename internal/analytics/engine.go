package analytics

import "hockeyAnalytics/internal/models"

type AnalyticsResult struct {
	BaseScore       float64 `json:"base_score"`
	NormalizedScore float64 `json:"normalized_score"`
	ContextScore    float64 `json:"context_score"`
	Overall         float64 `json:"overall"`
}

func OverallStat(stats models.PlayerSeasonStats, position string) float64 {

	score :=
		NormalizedModel(stats)*0.3 +
			ContextModel(stats, position)*0.4

	return score
}

func CalculateAllModels(stats models.PlayerSeasonStats, position string) AnalyticsResult {

	return AnalyticsResult{
		BaseScore: BaseStatModel(stats),

		NormalizedScore: NormalizedModel(stats),

		ContextScore: ContextModel(stats, position),

		Overall: OverallStat(stats, position),
	}
}
