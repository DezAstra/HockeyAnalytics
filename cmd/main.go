package main

import (
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	database.ConnectDB()

	router := gin.Default()

	router.LoadHTMLGlob("web/templates/*")

	router.GET("/", handlers.HomePage)
	router.POST("/import/csv", handlers.ImportCSV)
	router.GET("/analytics", handlers.GetPlayersAnalytics)
	router.GET("/analytics/leaders/:stat", handlers.GetLeaderboard)
	router.GET("/players", handlers.GetPlayers)

	router.Run(":8080")
}
