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

	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	router.GET("/player", func(c *gin.Context) {
		c.File("./static/player.html")
	})

	router.GET("/comparison", func(c *gin.Context) {
		c.File("./static/compare.html")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/import/csv", handlers.ImportCSV)
	router.POST("/nhl/sync", handlers.SyncNHLSeason)

	router.GET("/nhl/stats", handlers.GetNHLSeasonStats)

	router.GET("/analytics", handlers.GetPlayersAnalytics)
	router.GET("/analytics/leaders/:stat", handlers.GetLeaderboard)
	router.GET("/analytics/player/:id/history", handlers.GetPlayerHistory)
	router.GET("/analytics/compare", handlers.ComparePlayers)

	router.GET("/players", handlers.GetPlayers)
	router.GET("/players/:id/career", handlers.GetPlayerCareer)
	router.GET("/players/:id/similar", handlers.GetSimilarPlayers)

	router.GET(
		"/seasons",
		handlers.GetSeasons,
	)

	router.GET(
		"/compare",
		func(c *gin.Context) {
			c.File("./static/comparison.html")
		},
	)

	router.GET(
		"/team/:team",
		func(c *gin.Context) {
			c.File("./static/team.html")
		},
	)

	router.GET(
		"/api/team/:team",
		handlers.GetTeam,
	)

	router.Run(":8080")
}
