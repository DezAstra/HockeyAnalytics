package mappers

import (
	"hockeyAnalytics/internal/dto/nhl"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/utils"
	"strings"
)

func ExtractLastTeam(
	teams string,
) string {

	splitted :=
		strings.Split(
			teams,
			",",
		)

	return strings.TrimSpace(
		splitted[len(splitted)-1],
	)
}

func MapNHLPlayerToModel(
	data nhl.CombinedPlayerStats,
) models.Player {

	return models.Player{
		NHLID:    data.PlayerID,
		Name:     data.SkaterFullName,
		Position: data.PositionCode,
	}
}

func MapNHLStatsToModel(
	data nhl.CombinedPlayerStats,
	playerID uint,
) models.PlayerSeasonStats {

	FoP := data.FaceoffWinPct * 100

	totalTOI := data.TimeOnIcePerGame

	season :=
		utils.FormatNHLSeason(
			data.SeasonID,
		)

	return models.PlayerSeasonStats{
		PlayerID:    playerID,
		Season:      season,
		GamesPlayed: data.GamesPlayed,

		Goals:     data.Goals,
		Assists:   data.Assists,
		Points:    data.Points,
		PlusMinus: data.PlusMinus,

		PenaltyMinutes: data.PenaltyMinutes,

		EvenStrengthGoals: data.EvenStrengthGoals,
		PowerPlayGoals:    data.PowerPlayGoals,
		ShortHandedGoals:  data.ShortHandedGoals,

		Shots:        data.Shots,
		BlockedShots: data.BlockedShots,
		Hits:         data.Hits,

		FaceoffPercent: &FoP,

		TimeOfIce: &totalTOI,
	}
}
