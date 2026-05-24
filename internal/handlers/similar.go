package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"hockeyAnalytics/internal/analytics"
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/utils"

	"github.com/gin-gonic/gin"
)

// SimilarPlayerResponse — структура ответа для фронтенда
type SimilarPlayerResponse struct {
	Player       string  `json:"player"`
	PlayerID     uint    `json:"player_id"`
	NHLID        int     `json:"nhl_id"`
	Team         string  `json:"team"`
	Position     string  `json:"position"`
	Similarity   float64 `json:"similarity"`
	OverallScore float64 `json:"overall_score"`
	Archetype    string  `json:"archetype"`
}

// Потокобезопасный глобальный экземпляр движка аналитики.
// (В идеале передавать его через DI структуру хэндлера, но для быстрой интеграции можно инициализировать так).
var analyticsEngine = analytics.NewAnalyticsEngine()

// GetSimilarPlayers возвращает список топ-10 похожих игроков с использованием пакетного расчета
func GetSimilarPlayers(c *gin.Context) {
	idParam := c.Param("id")
	playerID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}
	season := c.Query("season")
	if season != "" {
		displaySeason, err := utils.ToDisplaySeason(season)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		season = displaySeason
	}

	// 1. Находим целевого игрока
	var target models.Player
	if err := database.DB.Preload("Stats").First(&target, playerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		return
	}
	if len(target.Stats) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "player stats not found"})
		return
	}

	// 2. Определяем рабочий сезон (линейным поиском находим самый свежий, если не передан явный)
	var targetStats models.PlayerSeasonStats
	if season != "" {
		found := false
		for _, s := range target.Stats {
			if s.Season == season {
				targetStats = s
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "stats for specified season not found"})
			return
		}
	} else {
		targetStats = target.Stats[0]
		for _, s := range target.Stats {
			if s.Season > targetStats.Season {
				targetStats = s
			}
		}
	}
	actualSeason := targetStats.Season

	// 3. Вместо ручного расчета моделей ВЫЗЫВАЕМ НАШ ДВИЖОК.
	// Он сам заберет всю лигу из БД, посчитает среднее/отклонение и вернет мапу с готовыми Overall.
	batchAnalytics, err := analyticsEngine.GetSeasonAnalytics(c.Request.Context(), actualSeason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate batch analytics: " + err.Error()})
		return
	}

	// Получаем готовые расчетные баллы для нашего целевого игрока из мапы
	targetResult, exists := batchAnalytics[target.ID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "target player did not pass analytics filter (e.g. min games)"})
		return
	}

	// 4. Запрашиваем из БД статистику всех остальных игроков для расчета дистанции схожести
	var allSeasonStats []models.PlayerSeasonStats
	if err := database.DB.
		Preload("Player").
		Where("season = ? AND player_id != ?", actualSeason, target.ID).
		Find(&allSeasonStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch season stats"})
		return
	}

	// 5. Просчитываем схожесть на основе готовых пакетных данных
	results := make([]SimilarPlayerResponse, 0, len(allSeasonStats))
	for _, stats := range allSeasonStats {
		if stats.Player.ID == 0 {
			continue
		}

		// Достаем результаты пакетного анализа для текущего хоккеиста за O(1)
		playerResult, passedFilter := batchAnalytics[stats.Player.ID]
		if !passedFilter {
			continue // Пропускаем игрока, если он не прошел фильтры пакетного расчета (например, сыграл < 5 матчей)
		}

		distance := analytics.SimilarityDistance(analytics.SimilarityInput{
			Stats:    targetStats,
			Position: target.Position,
			Overall:  targetResult.Overall,
		}, analytics.SimilarityInput{
			Stats:    stats,
			Position: stats.Player.Position,
			Overall:  playerResult.Overall,
		})
		similarity := analytics.SimilarityPercent(distance)

		// Определение динамического архетипа
		archetype := analytics.DetectArchetype(stats, stats.Player.Position)

		results = append(results, SimilarPlayerResponse{
			Player:       stats.Player.Name,
			PlayerID:     stats.Player.ID,
			NHLID:        stats.Player.NHLID,
			Team:         stats.Team,
			Position:     stats.Player.Position,
			Similarity:   similarity,
			OverallScore: playerResult.Overall, // Значение уже округлено внутри calculate_batch.go
			Archetype:    archetype,
		})
	}

	// 6. Сортируем результаты по убыванию процента сходства
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Ограничиваем выборку топ-10 похожими игроками
	if len(results) > 10 {
		results = results[:10]
	}

	c.JSON(http.StatusOK, results)
}
