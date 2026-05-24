package services

import (
	"fmt"
	"log"

	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/mappers"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/utils"
)

type SeasonImportResult struct {
	RequestedSeason string   `json:"requested_season"`
	APISeason       string   `json:"api_season"`
	DisplaySeason   string   `json:"display_season"`
	FetchedSummary  int      `json:"fetched_summary"`
	FetchedRealtime int      `json:"fetched_realtime"`
	MergedPlayers   int      `json:"merged_players"`
	PlayersCreated  int      `json:"players_created"`
	PlayersUpdated  int      `json:"players_updated"`
	NHLIDsAssigned  int      `json:"nhl_ids_assigned"`
	StatsCreated    int      `json:"stats_created"`
	StatsUpdated    int      `json:"stats_updated"`
	Skipped         int      `json:"skipped"`
	Errors          []string `json:"errors"`
	AlreadyExisted  bool     `json:"already_existed"`
}

func ImportSeason(
	season string,
) (SeasonImportResult, error) {
	apiSeason, err := utils.ToAPISeason(season)
	if err != nil {
		return SeasonImportResult{}, err
	}

	displaySeason, err := utils.ToDisplaySeason(season)
	if err != nil {
		return SeasonImportResult{}, err
	}

	result := SeasonImportResult{
		RequestedSeason: season,
		APISeason:       apiSeason,
		DisplaySeason:   displaySeason,
		Errors:          make([]string, 0),
	}

	fmt.Println("IMPORTING:", apiSeason)

	summary, err :=
		FetchSeasonSummary(
			apiSeason,
		)

	if err != nil {
		return result, err
	}

	result.FetchedSummary = len(summary)

	realtime, err :=
		FetchRealtimeStats(
			apiSeason,
		)

	if err != nil {
		return result, err
	}

	result.FetchedRealtime = len(realtime)

	data :=
		MergePlayerStats(
			summary,
			realtime,
		)

	result.MergedPlayers = len(data)

	for _, item := range data {
		playerModel :=
			mappers.MapNHLPlayerToModel(
				item,
			)

		playerResult, err :=
			UpsertPlayerWithResult(
				playerModel,
			)

		if err != nil {
			message := fmt.Sprintf(
				"failed to upsert player %s: %v",
				item.SkaterFullName,
				err,
			)

			log.Println(message)
			result.Errors = append(result.Errors, message)
			result.Skipped++
			continue
		}

		if playerResult.Created {
			result.PlayersCreated++
		}

		if playerResult.Updated {
			result.PlayersUpdated++
		}

		if playerResult.AssignedNHLID {
			result.NHLIDsAssigned++
		}

		stats :=
			mappers.MapNHLStatsToModel(
				item,
				playerResult.Player.ID,
			)

		stats.Team =
			mappers.ExtractLastTeam(
				item.TeamAbbrevs,
			)

		statsResult, err :=
			SaveSeasonStatsWithResult(
				stats,
			)

		if err != nil {
			message := fmt.Sprintf(
				"failed to save stats for %s: %v",
				item.SkaterFullName,
				err,
			)

			log.Println(message)
			result.Errors = append(result.Errors, message)
			result.Skipped++
			continue
		}

		if statsResult.Created {
			result.StatsCreated++
		}

		if statsResult.Updated {
			result.StatsUpdated++
		}
	}

	return result, nil
}

func SyncSeasonIfMissing(
	season string,
) error {
	_, err := SyncSeasonIfMissingWithResult(season)
	return err
}

func SyncSeasonIfMissingWithResult(
	season string,
) (SeasonImportResult, error) {
	displaySeason, err := utils.ToDisplaySeason(season)
	if err != nil {
		return SeasonImportResult{}, err
	}

	var count int64

	err = database.DB.
		Model(&models.PlayerSeasonStats{}).
		Where(
			"season = ?",
			displaySeason,
		).
		Count(&count).Error

	if err != nil {
		return SeasonImportResult{}, err
	}

	fmt.Println("SEASON:", displaySeason)
	fmt.Println("COUNT:", count)

	if count > 0 {
		apiSeason, err := utils.ToAPISeason(displaySeason)
		if err != nil {
			return SeasonImportResult{}, err
		}

		return SeasonImportResult{
			RequestedSeason: season,
			APISeason:       apiSeason,
			DisplaySeason:   displaySeason,
			AlreadyExisted:  true,
			StatsUpdated:    int(count),
			Errors:          make([]string, 0),
		}, nil
	}

	return ImportSeason(displaySeason)
}
