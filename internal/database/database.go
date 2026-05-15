package database

import (
	"fmt"
	"log"

	"hockeyAnalytics/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=123123 dbname=hockey_analytics port=5432 sslmode=disable"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	fmt.Println("Database connected")

	database.AutoMigrate(
		&models.Team{},
		&models.Player{},
		&models.PlayerStats{},
	)

	DB = database
}
