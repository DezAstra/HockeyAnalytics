package analytics

import "hockeyAnalytics/internal/models"

// DetectArchetype определяет игровой профиль по темповым метрикам per-82, а не по сырым totals.
// Так архетипы меньше зависят от количества матчей и корректнее работают для неполных сезонов.
func DetectArchetype(
	stats models.PlayerSeasonStats,
	position string,
) string {
	pos := NormalizePosition(position)
	games := stats.GamesPlayed

	if games <= 0 {
		return "Баланс"
	}

	goals82 := Per82(stats.Goals, games)
	assists82 := Per82(stats.Assists, games)
	points82 := Per82(stats.Points, games)
	shots82 := Per82(stats.Shots, games)
	hits82 := Per82(stats.Hits, games)
	blocks82 := Per82(stats.BlockedShots, games)
	pim82 := Per82(stats.PenaltyMinutes, games)

	if goals82 >= 35 && shots82 >= 220 && goals82 > assists82 {
		return "Снайпер"
	}

	if points82 >= 65 && assists82 >= goals82*0.85 && goals82 >= assists82*0.65 {
		return "Бомбардир"
	}

	if assists82 >= 55 && assists82-goals82 >= 30 {
		return "Ассистент"
	}

	if pos == "D" && points82 >= 45 {
		return "Атакующий защитник"
	}

	if pos == "D" && blocks82 >= 130 {
		return "Защитник-стена"
	}

	if pim82 >= 70 {
		return "Нарушитель"
	}

	if hits82 >= 150 && goals82 < 20 {
		return "Силовик"
	}

	if pos == "C" && stats.FaceoffPercent != nil && *stats.FaceoffPercent >= 57 {
		return "Специалист по вбрасываниям"
	}

	return "Баланс"
}
