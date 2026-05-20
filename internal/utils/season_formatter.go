package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func FormatSeason(
	endYear int,
) string {

	start :=
		(endYear - 1) % 100

	end :=
		endYear % 100

	return fmt.Sprintf(
		"%02d/%02d",
		start,
		end,
	)
}

func FormatNHLSeason(
	seasonID int,
) string {

	start :=
		seasonID / 10000

	end :=
		seasonID % 10000

	return fmt.Sprintf(
		"%02d/%02d",
		start%100,
		end%100,
	)
}

func NormalizeSeason(season string) (string, error) {
	// уже в нужном формате
	if len(season) == 8 && !strings.Contains(season, "/") {
		return season, nil
	}

	parts := strings.Split(season, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid season format: %s", season)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", err
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}

	startYear := 2000 + start
	endYear := 2000 + end

	return fmt.Sprintf("%d%d", startYear, endYear), nil
}
