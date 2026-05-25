package analytics

import (
	"hockeyAnalytics/internal/models"
	"math"
)

type AnalyticsResult struct {
	BaseScore       float64 `json:"base_score"`
	NormalizedScore float64 `json:"normalized_score"`
	ContextScore    float64 `json:"context_score"`
	Overall         float64 `json:"overall"`
}

type ModelDistribution struct {
	Mean   float64
	StdDev float64
}

func CalculateZScore(
	value float64,
	dist ModelDistribution,
) float64 {
	if dist.StdDev == 0 {
		return 0
	}

	return (value - dist.Mean) / dist.StdDev
}

func getOverallWeights(
	group string,
) OverallWeights {
	if weights, ok := PositionOverallWeights[group]; ok {
		return weights
	}

	return DefaultOverallWeights
}

func normalizeOverallWeights(
	weights OverallWeights,
) OverallWeights {
	total :=
		weights.Normalized +
			weights.Context

	if total <= 0 {
		return DefaultOverallWeights
	}

	return OverallWeights{
		Normalized: weights.Normalized / total,
		Context:    weights.Context / total,
	}
}

func CalculateOverallScore(
	normScore float64,
	normDist ModelDistribution,
	contextScore float64,
	contextDist ModelDistribution,
) float64 {
	return CalculateOverallScoreWithConfidence(
		normScore,
		normDist,
		contextScore,
		contextDist,
		1,
	)
}

// CalculateOverallScoreWithConfidence оставлен для обратной совместимости.
// Если позиционная группа неизвестна, используются веса из DefaultConfig.
func CalculateOverallScoreWithConfidence(
	normScore float64,
	normDist ModelDistribution,
	contextScore float64,
	contextDist ModelDistribution,
	confidence float64,
) float64 {
	cfg := DefaultConfig

	zNorm :=
		CalculateZScore(
			normScore,
			normDist,
		)

	zContext :=
		CalculateZScore(
			contextScore,
			contextDist,
		)

	combinedZ :=
		(zNorm * cfg.OverallNormWeight) +
			(zContext * cfg.OverallContextWeight)

	return scaleOverallWithConfidence(
		combinedZ,
		confidence,
	)
}

// CalculateOverallScoreWithGroup считает Overall через Normalized и Context.
//
// BaseScore здесь намеренно НЕ используется,
// потому что Normalized уже является развитием базовой продуктивности.
// Если добавить Base отдельно, голы/ассисты/очки будут частично учитываться дважды.
func CalculateOverallScoreWithGroup(
	normScore float64,
	normDist ModelDistribution,
	contextScore float64,
	contextDist ModelDistribution,
	confidence float64,
	group string,
) float64 {
	zNorm :=
		CalculateZScore(
			normScore,
			normDist,
		)

	zContext :=
		CalculateZScore(
			contextScore,
			contextDist,
		)

	weights :=
		normalizeOverallWeights(
			getOverallWeights(group),
		)

	combinedZ :=
		(zNorm * weights.Normalized) +
			(zContext * weights.Context)

	return scaleOverallWithConfidence(
		combinedZ,
		confidence,
	)
}

func scaleOverallWithConfidence(
	combinedZ float64,
	confidence float64,
) float64 {
	cfg := DefaultConfig

	if confidence < 0 {
		confidence = 0
	}

	if confidence > 1 {
		confidence = 1
	}

	overall :=
		cfg.OverallMean +
			(combinedZ * cfg.OverallStdScale * confidence)

	return Round1(overall)
}

func CalculateDistribution(
	values []float64,
) ModelDistribution {
	if len(values) == 0 {
		return ModelDistribution{}
	}

	sum := 0.0

	for _, value := range values {
		sum += value
	}

	mean :=
		sum / float64(len(values))

	varianceSum := 0.0

	for _, value := range values {
		varianceSum += math.Pow(
			value-mean,
			2,
		)
	}

	variance :=
		varianceSum / float64(len(values))

	stdDev :=
		math.Sqrt(variance)

	return ModelDistribution{
		Mean:   mean,
		StdDev: stdDev,
	}
}

type rawPlayerScores struct {
	playerID   uint
	position   string
	group      string
	games      int
	base       float64
	norm       float64
	ctx        float64
	confidence float64
}

func pickDistribution(
	group string,
	leagueDist ModelDistribution,
	groupDists map[string]ModelDistribution,
	groupCounts map[string]int,
) ModelDistribution {
	if groupCounts[group] >= DefaultConfig.PositionDistributionMin {
		if dist, ok := groupDists[group]; ok && dist.StdDev > 0 {
			return dist
		}
	}

	return leagueDist
}

// CalculateBatch выполняет пакетный расчет метрик для всей выборки игроков за сезон.
//
// Логика:
//  1. Base считается и возвращается как отдельная метрика.
//  2. Overall считается только из Normalized и Context.
//  3. Распределения строятся по надежной части выборки.
//  4. Если в позиционной группе достаточно игроков, z-score считается внутри этой группы.
//  5. Позиционные веса берутся из config.go.
//  6. Игроки с малым числом матчей не удаляются,
//     но их Overall сжимается к среднему значению через ConfidenceFactor.
func CalculateBatch(
	playersStats []models.PlayerSeasonStats,
) map[uint]AnalyticsResult {
	results :=
		make(map[uint]AnalyticsResult)

	leagueNormScores :=
		make([]float64, 0, len(playersStats))

	leagueContextScores :=
		make([]float64, 0, len(playersStats))

	groupNormScores :=
		make(map[string][]float64)

	groupContextScores :=
		make(map[string][]float64)

	groupCounts :=
		make(map[string]int)

	rawCalculations :=
		make([]rawPlayerScores, 0, len(playersStats))

	for _, stats := range playersStats {
		if stats.GamesPlayed <= 0 {
			continue
		}

		position := ""
		playerID := stats.PlayerID

		if stats.Player.ID != 0 {
			position = stats.Player.Position
			playerID = stats.Player.ID
		}

		base :=
			BaseStatModel(stats)

		norm :=
			NormalizedModel(stats)

		ctx :=
			ContextModel(
				stats,
				position,
			)

		group :=
			PositionGroup(position)

		confidence :=
			ConfidenceFactor(
				stats.GamesPlayed,
			)

		rawCalculations = append(
			rawCalculations,
			rawPlayerScores{
				playerID:   playerID,
				position:   position,
				group:      group,
				games:      stats.GamesPlayed,
				base:       base,
				norm:       norm,
				ctx:        ctx,
				confidence: confidence,
			},
		)

		if stats.GamesPlayed >= DefaultConfig.DistributionMinGames {
			leagueNormScores = append(
				leagueNormScores,
				norm,
			)

			leagueContextScores = append(
				leagueContextScores,
				ctx,
			)

			groupNormScores[group] = append(
				groupNormScores[group],
				norm,
			)

			groupContextScores[group] = append(
				groupContextScores[group],
				ctx,
			)

			groupCounts[group]++
		}
	}

	leagueNormDist :=
		CalculateDistribution(
			leagueNormScores,
		)

	leagueContextDist :=
		CalculateDistribution(
			leagueContextScores,
		)

	groupNormDists :=
		make(map[string]ModelDistribution)

	groupContextDists :=
		make(map[string]ModelDistribution)

	for group, values := range groupNormScores {
		groupNormDists[group] =
			CalculateDistribution(values)
	}

	for group, values := range groupContextScores {
		groupContextDists[group] =
			CalculateDistribution(values)
	}

	for _, raw := range rawCalculations {
		normDist :=
			pickDistribution(
				raw.group,
				leagueNormDist,
				groupNormDists,
				groupCounts,
			)

		contextDist :=
			pickDistribution(
				raw.group,
				leagueContextDist,
				groupContextDists,
				groupCounts,
			)

		overall :=
			CalculateOverallScoreWithGroup(
				raw.norm,
				normDist,
				raw.ctx,
				contextDist,
				raw.confidence,
				raw.group,
			)

		results[raw.playerID] =
			AnalyticsResult{
				BaseScore:       Round1(raw.base),
				NormalizedScore: Round1(raw.norm),
				ContextScore:    Round1(raw.ctx),
				Overall:         overall,
			}
	}

	return results
}
