// @title HockeyAnalytics API
// @version 1.5
// @description Аналитика статистических показателей хоккейных игроков
// @host localhost:8080
// @BasePath /
package main

import (
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/handlers"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "hockeyAnalytics/docs"
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
	router.GET("/analytics/player/:id/history", handlers.GetPlayerHistory)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/analytics/compare", handlers.ComparePlayers)

	router.Run(":8080")
}
