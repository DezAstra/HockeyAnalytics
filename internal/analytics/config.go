package analytics

// AnalyticsConfig хранит веса и пороги аналитических моделей.
// Значения вынесены в один файл, чтобы модель было проще объяснять,
// настраивать и защищать в документации.
type AnalyticsConfig struct {
	GoalWeight            float64
	AssistWeight          float64
	EvenStrengthGoalBonus float64
	PowerPlayGoalModifier float64
	ShortHandedGoalBonus  float64
	PlusMinusWeight       float64
	PenaltyMinuteWeight   float64
	PenaltyGraceMinutes   int

	NormalizedSeasonGames    float64
	NormalizedSampleGames    float64
	NormalizedAvgBasePerGame float64

	OverallNormWeight       float64
	OverallContextWeight    float64
	OverallMean             float64
	OverallStdScale         float64
	ConfidenceFullGames     float64
	DistributionMinGames    int
	PositionDistributionMin int
}

var DefaultConfig = AnalyticsConfig{
	GoalWeight:            1.00,
	AssistWeight:          0.80,
	EvenStrengthGoalBonus: 0.10,
	PowerPlayGoalModifier: -0.05,
	ShortHandedGoalBonus:  0.80,
	PlusMinusWeight:       0.25,
	PenaltyMinuteWeight:   0.12,
	PenaltyGraceMinutes:   20,

	NormalizedSeasonGames:    82,
	NormalizedSampleGames:    15,
	NormalizedAvgBasePerGame: 0.70,

	OverallNormWeight:       0.50,
	OverallContextWeight:    0.50,
	OverallMean:             50,
	OverallStdScale:         15,
	ConfidenceFullGames:     40,
	DistributionMinGames:    5,
	PositionDistributionMin: 10,
}

type OverallWeights struct {
	Normalized float64
	Context    float64
}

var PositionOverallWeights = map[string]OverallWeights{
	"C": {
		Normalized: 0.35,
		Context:    0.65,
	},
	"W": {
		Normalized: 0.55,
		Context:    0.45,
	},
	"D": {
		Normalized: 0.25,
		Context:    0.75,
	},
}

var DefaultOverallWeights = OverallWeights{
	Normalized: 0.50,
	Context:    0.50,
}
