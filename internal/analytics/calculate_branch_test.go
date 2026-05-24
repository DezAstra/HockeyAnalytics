package analytics

import "testing"

func TestDefenseOverallWeightsPrioritizeContext(t *testing.T) {
	weights := normalizeOverallWeights(getOverallWeights("D"))

	if weights.Context <= weights.Normalized {
		t.Fatalf("expected D context weight to be greater than normalized weight, got context=%v normalized=%v", weights.Context, weights.Normalized)
	}

	if weights.Context < 0.70 {
		t.Fatalf("expected D context weight to be at least 0.70, got %v", weights.Context)
	}
}

func TestOverallScoreUsesPositionWeights(t *testing.T) {
	normDist := ModelDistribution{Mean: 50, StdDev: 10}
	contextDist := ModelDistribution{Mean: 50, StdDev: 10}

	// Same scores, different position groups.
	// Defensemen should benefit more from a high context score and lower normalized score.
	defenseOverall := CalculateOverallScoreWithGroup(
		40,
		normDist,
		80,
		contextDist,
		1,
		"D",
	)

	wingOverall := CalculateOverallScoreWithGroup(
		40,
		normDist,
		80,
		contextDist,
		1,
		"W",
	)

	if defenseOverall <= wingOverall {
		t.Fatalf("expected defense overall to be higher when context is high, got D=%v W=%v", defenseOverall, wingOverall)
	}
}
