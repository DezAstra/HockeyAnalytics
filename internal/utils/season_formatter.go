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

func ToAPISeason(season string) (string, error) {

	season = strings.TrimSpace(season)

	if len(season) == 8 && !strings.Contains(season, "/") {
		_, err := strconv.Atoi(season)
		if err != nil {
			return "", fmt.Errorf("invalid NHL season format: %s", season)
		}

		return season, nil
	}

	parts := strings.Split(season, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid season format: %s", season)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid season start year: %s", season)
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid season end year: %s", season)
	}

	if start < 0 || start > 99 || end < 0 || end > 99 {
		return "", fmt.Errorf("invalid season range: %s", season)
	}

	startCentury := 2000
	if start >= 90 {
		startCentury = 1900
	}

	startYear := startCentury + start
	endYear := startCentury + end

	if end < start {
		endYear += 100
	}

	return fmt.Sprintf(
		"%04d%04d",
		startYear,
		endYear,
	), nil
}

func ToDisplaySeason(season string) (string, error) {

	season = strings.TrimSpace(season)

	if strings.Contains(season, "/") {
		parts := strings.Split(season, "/")
		if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
			return "", fmt.Errorf("invalid season format: %s", season)
		}

		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return "", fmt.Errorf("invalid season start year: %s", season)
		}

		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid season end year: %s", season)
		}

		return fmt.Sprintf(
			"%02d/%02d",
			start%100,
			end%100,
		), nil
	}

	if len(season) != 8 {
		return "", fmt.Errorf("invalid NHL season format: %s", season)
	}

	seasonID, err := strconv.Atoi(season)
	if err != nil {
		return "", fmt.Errorf("invalid NHL season format: %s", season)
	}

	return FormatNHLSeason(seasonID), nil
}

// NormalizeSeason is kept for backward compatibility.
// New code should use ToAPISeason or ToDisplaySeason explicitly.
func NormalizeSeason(season string) (string, error) {
	return ToAPISeason(season)
}
