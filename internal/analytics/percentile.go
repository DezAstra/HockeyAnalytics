package analytics

import "sort"

func CalculatePercentile(
	value float64,
	allValues []float64,
) float64 {

	sort.Float64s(allValues)

	count := 0

	for _, v := range allValues {

		if v <= value {

			count++
		}
	}

	return float64(count) /
		float64(len(allValues)) * 100
}
