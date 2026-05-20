package services

import (
	"fmt"

	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/mappers"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/utils"
)

func ImportSeason(
	season string,
) error {

	season, err := utils.NormalizeSeason(season)
	if err != nil {
		return err
	}

	fmt.Println("IMPORTING:", season)

	summary, err :=
		FetchSeasonSummary(
			season,
		)

	if err != nil {
		return err
	}

	realtime, err :=
		FetchRealtimeStats(
			season,
		)

	if err != nil {
		return err
	}

	data :=
		MergePlayerStats(
			summary,
			realtime,
		)

	for _, item := range data {

		playerModel :=
			mappers.MapNHLPlayerToModel(
				item,
			)

		player, err :=
			UpsertPlayer(
				playerModel,
			)

		if err != nil {
			continue
		}

		stats :=
			mappers.MapNHLStatsToModel(
				item,
				player.ID,
			)

		stats.Team =
			mappers.ExtractLastTeam(
				item.TeamAbbrevs,
			)

		err =
			SaveSeasonStats(
				stats,
			)

		if err != nil {
			continue
		}
	}

	return nil
}

func SyncSeasonIfMissing(
	season string,
) error {

	var count int64

	err := database.DB.
		Model(&models.PlayerSeasonStats{}).
		Where(
			"season = ?",
			season,
		).
		Count(&count).Error

	if err != nil {
		return err
	}

	season, err = utils.NormalizeSeason(season)
	if err != nil {
		return err
	}

	fmt.Println("SEASON:", season)
	fmt.Println("COUNT:", count)

	if count > 0 {
		return nil
	}

	return ImportSeason(season)
}
