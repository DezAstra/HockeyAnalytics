package handlers

import (
	"net/http"

	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"

	"github.com/gin-gonic/gin"
)

func GetPlayers(c *gin.Context) {

	var players []models.Player

	database.DB.
		Preload("Team").
		Preload("Stats").
		Find(&players)

	c.JSON(http.StatusOK, players)
}
