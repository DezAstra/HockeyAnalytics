package analytics

import (
	"context"
	"errors"
	"sync"
	"time"

	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
)

// CalculateAllModels вычисляет все аналитические оценки для конкретного игрока.
// В отличие от старой версии, теперь функция принимает предопределенные распределения всей лиги
// (normDist и contextDist), что полностью решает проблему неверного расчета распределения на основе одного игрока.
func CalculateAllModels(
	stats models.PlayerSeasonStats,
	position string,
	normDist ModelDistribution,
	contextDist ModelDistribution,
) AnalyticsResult {
	normScore := NormalizedModel(stats)
	ctxScore := ContextModel(stats, position)

	return AnalyticsResult{
		BaseScore:       Round1(BaseStatModel(stats)),
		NormalizedScore: Round1(normScore),
		ContextScore:    Round1(ctxScore),
		// Передаем сырые баллы игрока, глобальные параметры распределения сезона и доверие к выборке
		Overall: CalculateOverallScoreWithConfidence(
			normScore,
			normDist,
			ctxScore,
			contextDist,
			ConfidenceFactor(stats.GamesPlayed),
		),
	}
}

// ============================================================================
// Слой архитектурной оркестрации (Сервисный Движок Аналитики)
// ============================================================================

// Engine описывает интерфейс аналитического движка для удобного тестирования и DI
type Engine interface {
	GetSeasonAnalytics(ctx context.Context, season string) (map[uint]AnalyticsResult, error)
	ClearCache()
}

// AnalyticsEngine инкапсулирует логику обращения к БД, кэширования и запуска пакетной математики
type AnalyticsEngine struct {
	cache      map[string]map[uint]AnalyticsResult
	cacheMutex sync.RWMutex
	cacheTTL   time.Duration
	lastUpdate map[string]time.Time
}

// NewAnalyticsEngine — конструктор движка аналитики
func NewAnalyticsEngine() *AnalyticsEngine {
	return &AnalyticsEngine{
		cache:      make(map[string]map[uint]AnalyticsResult),
		lastUpdate: make(map[string]time.Time),
		cacheTTL:   5 * time.Minute, // Кэшируем расчеты на 5 минут для оптимизации нагрузки на БД
	}
}

// GetSeasonAnalytics координирует сбор сырых данных из БД за весь сезон,
// передает их в пакетный процессор и кэширует результаты расчета.
func (e *AnalyticsEngine) GetSeasonAnalytics(ctx context.Context, season string) (map[uint]AnalyticsResult, error) {
	if season == "" {
		return nil, errors.New("season parameter cannot be empty")
	}

	// 1. Проверяем кэш (используем RLock для безопасного параллельного чтения)
	e.cacheMutex.RLock()
	cachedData, exists := e.cache[season]
	lastUpd, hasTime := e.lastUpdate[season]
	e.cacheMutex.RUnlock()

	if exists && hasTime && time.Since(lastUpd) < e.cacheTTL {
		return cachedData, nil
	}

	// 2. Если кэша нет, запрашиваем статистику ВСЕХ игроков за этот сезон из ORM
	var allSeasonStats []models.PlayerSeasonStats
	err := database.DB.WithContext(ctx).
		Preload("Player").
		Where("season = ?", season).
		Find(&allSeasonStats).Error

	if err != nil {
		return nil, err
	}

	if len(allSeasonStats) == 0 {
		return nil, errors.New("no statistics found for the specified season: " + season)
	}

	// 3. Вызываем функцию пакетного расчета из твоего пакета calculate_batch.go
	// Она один раз соберет срезы, посчитает правильные средние/отклонения по лиге и вернет мапу результатов
	analysisMap := CalculateBatch(allSeasonStats)

	// 4. Записываем результаты в кэш (используем Lock для безопасной записи)
	e.cacheMutex.Lock()
	e.cache[season] = analysisMap
	e.lastUpdate[season] = time.Now()
	e.cacheMutex.Unlock()

	return analysisMap, nil
}

// ClearCache принудительно сбрасывает кэш аналитики (необходимо вызывать при импорте новых CSV)
func (e *AnalyticsEngine) ClearCache() {
	e.cacheMutex.Lock()
	defer e.cacheMutex.Unlock()
	e.cache = make(map[string]map[uint]AnalyticsResult)
	e.lastUpdate = make(map[string]time.Time)
}
