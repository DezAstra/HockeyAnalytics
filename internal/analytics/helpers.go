package analytics

func CalculateShootingPercent(
	goals int,
	shots int,
) float64 {

	if shots == 0 {
		return 0
	}

	return (float64(goals) / float64(shots)) * 100
}

func CalculateFaceoffPercent(
	won int,
	lost int,
) float64 {

	total := won + lost

	if total == 0 {
		return 0
	}

	return (float64(won) / float64(total)) * 100
}
