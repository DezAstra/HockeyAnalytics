package services

import (
	"errors"

	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PlayerUpsertResult struct {
	Player        models.Player
	Created       bool
	Updated       bool
	MatchedBy     string
	AssignedNHLID bool
}

type SeasonStatsSaveResult struct {
	Created bool
	Updated bool
}

func UpsertPlayer(
	player models.Player,
) (models.Player, error) {
	result, err := UpsertPlayerWithResult(player)
	return result.Player, err
}

func UpsertPlayerWithResult(
	player models.Player,
) (PlayerUpsertResult, error) {
	var existing models.Player

	// 1. NHL ID is the strongest identity key. If the player already exists by
	// NHL ID, update stable fields and return that exact row.
	if player.NHLID > 0 {
		err := database.DB.
			Where("nhl_id = ?", player.NHLID).
			First(&existing).Error

		if err == nil {
			changed := false

			if player.Name != "" && existing.Name != player.Name {
				existing.Name = player.Name
				changed = true
			}

			if player.Position != "" && existing.Position != player.Position {
				existing.Position = player.Position
				changed = true
			}

			if changed {
				if err := database.DB.Save(&existing).Error; err != nil {
					return PlayerUpsertResult{}, err
				}
			}

			return PlayerUpsertResult{
				Player:    existing,
				Updated:   changed,
				MatchedBy: "nhl_id",
			}, nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return PlayerUpsertResult{}, err
		}
	}

	// 2. NHL import with a new NHL ID may need to attach to an older CSV row.
	// We only attach when there is exactly one same-name player without NHL ID.
	// If a same-name player already has another real NHL ID, this is a namesake
	// and a new row must be created. Example: Elias Pettersson.
	if player.NHLID > 0 {
		var csvCandidates []models.Player

		err := database.DB.
			Where("LOWER(name) = LOWER(?) AND nhl_id = 0", player.Name).
			Find(&csvCandidates).Error

		if err != nil {
			return PlayerUpsertResult{}, err
		}

		if len(csvCandidates) == 1 {
			existing = csvCandidates[0]
			existing.NHLID = player.NHLID

			if player.Position != "" {
				existing.Position = player.Position
			}

			if err := database.DB.Save(&existing).Error; err != nil {
				return PlayerUpsertResult{}, err
			}

			return PlayerUpsertResult{
				Player:        existing,
				Updated:       true,
				MatchedBy:     "csv_name",
				AssignedNHLID: true,
			}, nil
		}

		if len(csvCandidates) > 1 {
			// Ambiguous CSV duplicates. Do not guess. Create a clean NHL-identified
			// row so stats import is not blocked.
			err = database.DB.Create(&player).Error
			if err != nil {
				return PlayerUpsertResult{}, err
			}

			return PlayerUpsertResult{
				Player:    player,
				Created:   true,
				MatchedBy: "created_ambiguous_csv_name",
			}, nil
		}

		err = database.DB.Create(&player).Error
		if err != nil {
			return PlayerUpsertResult{}, err
		}

		return PlayerUpsertResult{
			Player:    player,
			Created:   true,
			MatchedBy: "created_nhl_id",
		}, nil
	}

	// 3. CSV import has no NHL ID. Use case-insensitive name matching so repeated
	// CSV imports update the same player instead of creating duplicates.
	err := database.DB.
		Where("LOWER(name) = LOWER(?)", player.Name).
		Order("CASE WHEN nhl_id = 0 THEN 0 ELSE 1 END").
		First(&existing).Error

	if err == nil {
		changed := false

		if player.Position != "" && existing.Position != player.Position {
			existing.Position = player.Position
			changed = true
		}

		if changed {
			if err := database.DB.Save(&existing).Error; err != nil {
				return PlayerUpsertResult{}, err
			}
		}

		return PlayerUpsertResult{
			Player:    existing,
			Updated:   changed,
			MatchedBy: "name",
		}, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return PlayerUpsertResult{}, err
	}

	err = database.DB.Create(&player).Error
	if err != nil {
		return PlayerUpsertResult{}, err
	}

	return PlayerUpsertResult{
		Player:    player,
		Created:   true,
		MatchedBy: "created_csv",
	}, nil
}

func SaveSeasonStats(
	stats models.PlayerSeasonStats,
) error {
	_, err := SaveSeasonStatsWithResult(stats)
	return err
}

func SaveSeasonStatsWithResult(
	stats models.PlayerSeasonStats,
) (SeasonStatsSaveResult, error) {
	var existing models.PlayerSeasonStats

	err := database.DB.
		Where(
			"player_id = ? AND season = ?",
			stats.PlayerID,
			stats.Season,
		).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return SeasonStatsSaveResult{}, err
	}

	created := errors.Is(err, gorm.ErrRecordNotFound)

	err = database.DB.
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "player_id"},
					{Name: "season"},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					"team",
					"games_played",
					"goals",
					"assists",
					"points",
					"plus_minus",
					"penalty_minutes",
					"even_strength_goals",
					"power_play_goals",
					"short_handed_goals",
					"shots",
					"blocked_shots",
					"hits",
					"faceoffs_won",
					"faceoffs_lost",
					"faceoff_percent",
					"time_of_ice",
					"updated_at",
				}),
			},
		).
		Create(&stats).Error

	if err != nil {
		return SeasonStatsSaveResult{}, err
	}

	return SeasonStatsSaveResult{
		Created: created,
		Updated: !created,
	}, nil
}
