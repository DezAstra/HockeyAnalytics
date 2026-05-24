package mappers

import (
	"fmt"
	"strconv"
	"strings"

	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/utils"
)

const MinCSVColumns = 31

func ExtractCSVPlayerName(row []string) string {
	return strings.TrimSpace(
		strings.Split(row[1], "\\")[0],
	)
}

func atoi(value string) int {
	parsed, _ := strconv.Atoi(strings.TrimSpace(value))
	return parsed
}

func parseOptionalFloat(value string) *float64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}

	return &parsed
}

func MapCSVRowToStats(
	row []string,
	playerID uint,
) (models.PlayerSeasonStats, error) {

	if len(row) < MinCSVColumns {
		return models.PlayerSeasonStats{}, fmt.Errorf(
			"csv row has %d columns, need at least %d",
			len(row),
			MinCSVColumns,
		)
	}

	seasonYear, err :=
		strconv.Atoi(strings.TrimSpace(row[30]))

	if err != nil {
		return models.PlayerSeasonStats{}, fmt.Errorf(
			"invalid season year %q: %w",
			row[30],
			err,
		)
	}

	season :=
		utils.FormatSeason(
			seasonYear,
		)

	gamesPlayed := atoi(row[5])
	goals := atoi(row[6])
	assists := atoi(row[7])
	plusMinus := atoi(row[9])
	penaltyMinutes := atoi(row[10])
	evenStrengthGoals := atoi(row[12])
	powerPlayGoals := atoi(row[13])
	shortHandedGoals := atoi(row[14])
	shots := atoi(row[19])
	blockedShots := atoi(row[23])
	hits := atoi(row[24])
	faceoffsWon := atoi(row[25])
	faceoffsLost := atoi(row[26])

	var FoP float64
	totalFaceoffs :=
		faceoffsWon + faceoffsLost

	if totalFaceoffs > 0 {

		FoP =
			(float64(faceoffsWon) /
				float64(totalFaceoffs)) * 100
	}

	return models.PlayerSeasonStats{
		PlayerID:    playerID,
		Season:      season,
		GamesPlayed: gamesPlayed,

		Goals:     goals,
		Assists:   assists,
		Points:    goals + assists,
		PlusMinus: plusMinus,

		PenaltyMinutes: penaltyMinutes,

		EvenStrengthGoals: evenStrengthGoals,
		PowerPlayGoals:    powerPlayGoals,
		ShortHandedGoals:  shortHandedGoals,

		Shots:        shots,
		BlockedShots: blockedShots,
		Hits:         hits,

		FaceoffsWon:    faceoffsWon,
		FaceoffsLost:   faceoffsLost,
		FaceoffPercent: &FoP,

		TimeOfIce: parseOptionalFloat(row[21]),
	}, nil
}
