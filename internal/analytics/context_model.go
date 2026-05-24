package analytics

import "hockeyAnalytics/internal/models"

// ContextModel оценивает игрока с учетом роли и позиции.
// Модель оставлена эвристической:
// сначала считается универсальный вклад, затем позиция добавляет именно роль, а не полностью дублирует scoring.
func ContextModel(
	stats models.PlayerSeasonStats,
	position string,
) float64 {
	pos := NormalizePosition(position)
	score := 0.0

	// Универсальный атакующий вклад.
	score += float64(stats.Goals) * 1.00
	score += float64(stats.Assists) * 0.80
	score += float64(stats.Shots) * 0.08

	// Тип гола — модификатор.
	score += float64(stats.EvenStrengthGoals) * 0.10
	score -= float64(stats.PowerPlayGoals) * 0.05
	score += float64(stats.ShortHandedGoals) * 0.70

	// Двусторонние и силовые действия.
	score += float64(stats.BlockedShots) * 0.07
	score += float64(stats.Hits) * 0.04

	if stats.PenaltyMinutes > DefaultConfig.PenaltyGraceMinutes {
		excessPIM := stats.PenaltyMinutes - DefaultConfig.PenaltyGraceMinutes
		score -= float64(excessPIM) * 0.12
	}

	if stats.GamesPlayed > 20 {
		score += float64(stats.PlusMinus) * 0.20
	}

	if stats.Shots > 50 {
		shootingPercent := CalculateShootingPercent(stats.Goals, stats.Shots)
		score += (shootingPercent - 10) * 0.20
	}

	switch pos {
	case "D":
		// Для защитников усиливаем defensive involvement и playmaking,
		// но не заставляем их соревноваться с форвардами только голами.
		score -= float64(stats.Goals) * 0.20
		score += float64(stats.Assists) * 0.20
		score += float64(stats.BlockedShots) * 0.45
		score += float64(stats.Hits) * 0.12

	case "C":
		// Для центров важны розыгрыш, двухсторонний вклад и вбрасывания.
		score += float64(stats.Assists) * 0.20
		score += float64(stats.BlockedShots) * 0.08

		totalFaceoffs := stats.FaceoffsWon + stats.FaceoffsLost
		if totalFaceoffs > 100 {
			faceoffPercent := CalculateFaceoffPercent(stats.FaceoffsWon, stats.FaceoffsLost)
			score += (faceoffPercent - 50) * 0.35

			if totalFaceoffs >= 800 && faceoffPercent >= 55 {
				score += 2
			} else if totalFaceoffs >= 500 && faceoffPercent >= 55 {
				score += 1
			}
		} else if stats.FaceoffPercent != nil {
			score += (*stats.FaceoffPercent - 50) * 0.35
		}

	case "LW", "RW", "L", "R":
		// Для крайних нападающих усиливаем бросковый и финишерский профиль,
		// но только умеренно, чтобы не дублировать голы второй раз полностью.
		score += float64(stats.Goals) * 0.25
		score += float64(stats.Shots) * 0.05
		score += float64(stats.Hits) * 0.08
	}

	return score
}
