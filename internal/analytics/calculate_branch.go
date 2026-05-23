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
	Mean   float64 // Среднее значение по лиге
	StdDev float64 // Стандартное отклонение по лиге
}

// CalculateZScore рассчитывает Z-оценку конкретного показателя
func CalculateZScore(value float64, dist ModelDistribution) float64 {
	if dist.StdDev == 0 {
		return 0
	}
	return (value - dist.Mean) / dist.StdDev
}

// CalculateOverallScore объединяет Z-оценки в 100-балльную шкалу
func CalculateOverallScore(
	normScore float64, normDist ModelDistribution,
	contextScore float64, contextDist ModelDistribution,
) float64 {
	// 1. Переводим сырые скоры в Z-оценки относительно лиги
	zNorm := CalculateZScore(normScore, normDist)
	zContext := CalculateZScore(contextScore, contextDist)

	// 2. Объединяем с весами 50/50 (можно настроить под себя)
	combinedZ := (zNorm * 0.5) + (zContext * 0.5)

	// 3. Масштабируем Z-score в красивую 100-балльную шкалу.
	// Средний игрок (Z=0) получает 50 баллов.
	// Опережение на 1 стандартное отклонение (Z=1) дает +15 баллов (итого 65).
	overall := 50.0 + (combinedZ * 15.0)

	// Ограничиваем жесткими рамками
	if overall > 100 {
		overall = 100
	}
	if overall < 0 {
		overall = 0
	}

	// Округляем до 1 знака после запятой для красивого JSON/фронтенда
	return math.Round(overall*10) / 10
}

// CalculateDistribution вычисляет среднее и стандартное отклонение генеральной совокупности
func CalculateDistribution(values []float64) ModelDistribution {
	if len(values) == 0 {
		return ModelDistribution{}
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	varianceSum := 0.0
	for _, v := range values {
		varianceSum += math.Pow(v-mean, 2)
	}
	variance := varianceSum / float64(len(values))
	stdDev := math.Sqrt(variance)

	return ModelDistribution{Mean: mean, StdDev: stdDev}
}

// CalculateBatch выполняет пакетный расчет метрик для всей выборки игроков за сезон
func CalculateBatch(playersStats []models.PlayerSeasonStats) map[uint]AnalyticsResult {
	results := make(map[uint]AnalyticsResult)

	var normScores []float64
	var contextScores []float64

	// Внутренняя структура для сохранения промежуточных сырых расчетов
	type rawPlayerScores struct {
		playerID uint
		base     float64
		norm     float64
		ctx      float64
	}

	rawCalculations := make([]rawPlayerScores, 0, len(playersStats))

	// ШАГ 1: Считаем сырые баллы моделей для каждого игрока
	for _, stats := range playersStats {
		// Защита: пропускаем игроков с аномально маленьким ToI или GP,
		// чтобы они не искажали Mean и StdDev всей лиги.
		if stats.GamesPlayed < 5 {
			continue
		}

		// Безопасно извлекаем позицию игрока, если объект Player был подгружен (Preload)
		position := ""
		playerID := stats.PlayerID
		if stats.Player.ID != 0 {
			position = stats.Player.Position
			playerID = stats.Player.ID
		}

		base := BaseStatModel(stats)
		norm := NormalizedModel(stats)
		ctx := ContextModel(stats, position)

		rawCalculations = append(rawCalculations, rawPlayerScores{
			playerID: playerID,
			base:     base,
			norm:     norm,
			ctx:      ctx,
		})

		// Собираем массивы для вычисления параметров распределения лиги
		normScores = append(normScores, norm)
		contextScores = append(contextScores, ctx)
	}

	// ШАГ 2: Вычисляем глобальные параметры распределения лиги за этот сезон
	normDist := CalculateDistribution(normScores)
	contextDist := CalculateDistribution(contextScores)

	// ШАГ 3: Зная среднее по лиге, рассчитываем справедливый Overall через Z-score
	for _, raw := range rawCalculations {
		overall := CalculateOverallScore(raw.norm, normDist, raw.ctx, contextDist)

		results[raw.playerID] = AnalyticsResult{
			BaseScore:       math.Round(raw.base*10) / 10,
			NormalizedScore: math.Round(raw.norm*10) / 10,
			ContextScore:    math.Round(raw.ctx*10) / 10,
			Overall:         overall,
		}
	}

	return results
}
