package services

import (
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
)

func UpsertPlayer(
	player models.Player,
) (models.Player, error) {

	var existing models.Player

	// 1. Ищем по NHLID

	if player.NHLID != 0 {

		err :=
			database.DB.
				Where(
					"nhl_id = ?",
					player.NHLID,
				).
				First(&existing).Error

		if err == nil {

			existing.Name =
				player.Name

			existing.Position =
				player.Position

			database.DB.Save(
				&existing,
			)

			return existing, nil
		}
	}

	// 2. Ищем CSV игрока

	err :=
		database.DB.
			Where(
				"name = ?",
				player.Name,
			).
			First(&existing).Error

	if err == nil {

		// Привязываем NHLID
		// к существующему игроку

		if existing.NHLID == 0 {

			existing.NHLID =
				player.NHLID
		}

		existing.Position =
			player.Position

		database.DB.Save(
			&existing,
		)

		return existing, nil
	}

	// 3. Создаём нового

	err =
		database.DB.
			Create(&player).Error

	return player, err
}

func SaveSeasonStats(
	stats models.PlayerSeasonStats,
) error {

	var existing models.PlayerSeasonStats

	err :=
		database.DB.
			Where(
				"player_id = ? AND season = ?",
				stats.PlayerID,
				stats.Season,
			).
			First(&existing).Error

	if err == nil {

		existing.GamesPlayed = stats.GamesPlayed

		existing.Goals = stats.Goals

		existing.Assists = stats.Assists

		existing.Points = stats.Points

		existing.Team = stats.Team

		existing.PlusMinus = stats.PlusMinus

		existing.PenaltyMinutes = stats.PenaltyMinutes

		existing.EvenStrengthGoals = stats.EvenStrengthGoals

		existing.PowerPlayGoals = stats.PowerPlayGoals

		existing.ShortHandedGoals = stats.ShortHandedGoals

		existing.Shots = stats.Shots

		existing.FaceoffPercent = stats.FaceoffPercent

		existing.TimeOfIce = stats.TimeOfIce

		return database.DB.
			Save(&existing).Error
	}

	return database.DB.
		Create(&stats).Error
}
