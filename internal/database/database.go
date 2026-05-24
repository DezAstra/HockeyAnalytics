package database

import (
	"fmt"
	"log"
	"os"

	"hockeyAnalytics/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func envOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func buildDSN() string {
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		return dsn
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		envOrDefault("DB_HOST", "localhost"),
		envOrDefault("DB_USER", "postgres"),
		envOrDefault("DB_PASSWORD", "123123"),
		envOrDefault("DB_NAME", "hockey_analytics"),
		envOrDefault("DB_PORT", "5432"),
		envOrDefault("DB_SSLMODE", "disable"),
	)
}

func ConnectDB() {

	database, err := gorm.Open(
		postgres.Open(buildDSN()),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal(err)
	}

	DB = database

	if err := DB.AutoMigrate(
		&models.Player{},
		&models.PlayerSeasonStats{},
		&models.ImportLog{},
	); err != nil {
		log.Fatal(err)
	}

	if err := ensurePlayerIndexes(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected")
}

func ensurePlayerIndexes() error {
	// Старый uniqueIndex от GORM запрещал несколько CSV-игроков с nhl_id = 0.
	// Для CSV 0 означает "NHL ID пока неизвестен", поэтому уникальность нужна
	// только для реальных NHL ID, то есть для значений больше нуля.
	if err := DB.Exec(`DROP INDEX IF EXISTS idx_players_nhl_id`).Error; err != nil {
		return err
	}

	return DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_players_nhl_id_positive
		ON players (nhl_id)
		WHERE nhl_id > 0 AND deleted_at IS NULL
	`).Error
}
