package handlers

import (
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/services"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ImportCSV(c *gin.Context) {

	file, err := c.FormFile("file")

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file upload error",
		})

		return
	}

	tempFilePath := "./csv/" + file.Filename

	err = c.SaveUploadedFile(file, tempFilePath)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "file save error",
		})

		return
	}

	f, err := os.Open(tempFilePath)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "file open error",
		})

		return
	}

	defer f.Close()

	rows, err := services.ReadCSV(tempFilePath)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "csv read error",
		})

		return
	}

	for index, row := range rows {

		if index == 0 {
			continue
		}

		gamesPlayed, _ := strconv.Atoi(row[5])
		goals, _ := strconv.Atoi(row[6])
		assists, _ := strconv.Atoi(row[7])
		plusMinus, _ := strconv.Atoi(row[9])
		penaltyMinutes, _ := strconv.Atoi(row[10])

		var toi *float64

		if row[21] != "" {

			value, _ := strconv.ParseFloat(row[21], 64)
			toi = &value
		}

		teamName := row[4]

		var team models.Team

		database.DB.
			Where("name = ?", teamName).
			First(&team)

		if team.ID == 0 {

			team = models.Team{
				Name: teamName,
			}

			database.DB.Create(&team)
		}

		player := models.Player{
			Name:     strings.Split(row[1], "\\")[0],
			Position: row[3],
			TeamID:   team.ID,
		}

		database.DB.Create(&player)

		stats := models.PlayerStats{
			PlayerID: player.ID,

			GamesPlayed:    gamesPlayed,
			Goals:          goals,
			Assists:        assists,
			PlusMinus:      plusMinus,
			PenaltyMinutes: penaltyMinutes,
			TimeOfIce:      toi,
		}

		database.DB.Create(&stats)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "csv imported",
	})
}
