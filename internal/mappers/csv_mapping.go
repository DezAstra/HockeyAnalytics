package mappers

import (
	"strconv"
	"strings"

	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/utils"
)

func ExtractTeamName(
	row []string,
) string {

	return row[4]
}

func ExtractPlayerName(
	row []string,
) string {

	return strings.Split(
		row[1],
		"\\",
	)[0]
}

func ExtractPlayerPosition(
	row []string,
) string {

	return row[3]
}

func MapCSVRowToStats(
	row []string,
	playerID uint,
) models.PlayerSeasonStats {

	seasonYear, _ :=
		strconv.Atoi(row[30])

	season :=
		utils.FormatSeason(
			seasonYear,
		)
	gamesPlayed, _ := strconv.Atoi(row[5])

	goals, _ := strconv.Atoi(row[6])
	assists, _ := strconv.Atoi(row[7])
	plusMinus, _ := strconv.Atoi(row[9])

	penaltyMinutes, _ := strconv.Atoi(row[10])

	evenStrengthGoals, _ := strconv.Atoi(row[12])
	powerPlayGoals, _ := strconv.Atoi(row[13])
	shortHandedGoals, _ := strconv.Atoi(row[14])

	shots, _ := strconv.Atoi(row[19])
	blockedShots, _ := strconv.Atoi(row[23])

	hits, _ := strconv.Atoi(row[24])

	faceoffsWon, _ := strconv.Atoi(row[25])
	faceoffsLost, _ := strconv.Atoi(row[26])

	var toi *float64

	if row[21] != "" {

		value, _ := strconv.ParseFloat(
			row[21],
			64,
		)

		toi = &value
	}

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

		Goals:          goals,
		Assists:        assists,
		PlusMinus:      plusMinus,
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

		TimeOfIce: toi,
	}
}
