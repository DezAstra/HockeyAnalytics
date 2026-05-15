package main

import (
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	database.ConnectDB()

	router := gin.Default()

	router.POST("/import/csv", handlers.ImportCSV)
	router.GET("/players", handlers.GetPlayers)
	router.GET("/analytics", handlers.GetPlayersAnalytics)

	router.Run(":8080")
}
