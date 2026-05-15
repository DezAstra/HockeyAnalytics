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

// ImportCSV godoc
// @Summary импорт csv
// @Description Импортировать статистику csv файлом
// @Tags Импорт
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV File"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /import/csv [post]
func ImportCSV(c *gin.Context) {

	file, err := c.FormFile("file")

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ошибка загрузки csv файла",
		})

		return
	}

	tempFilePath := "./csv/" + file.Filename

	err = c.SaveUploadedFile(file, tempFilePath)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка сохранения csv файла",
		})

		return
	}

	f, err := os.Open(tempFilePath)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка открытия csv файла",
		})

		return
	}

	defer f.Close()

	rows, err := services.ReadCSV(tempFilePath)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка чтения csv файла",
		})

		return
	}

	for index, row := range rows {

		if index == 0 {
			continue
		}
		season := row[30]
		gamesPlayed, _ := strconv.Atoi(row[5])
		goals, _ := strconv.Atoi(row[6])
		assists, _ := strconv.Atoi(row[7])
		plusMinus, _ := strconv.Atoi(row[9])
		penaltyMinutes, _ := strconv.Atoi(row[10])
		evenStrengthGoals, _ := strconv.Atoi(row[12])
		powerPlayGoals, _ := strconv.Atoi(row[13])
		shortHandedGoals, _ := strconv.Atoi(row[14])
		shots, _ := strconv.Atoi(row[19])
		blockedShots, _ := strconv.Atoi(row[23])
		hits, _ := strconv.Atoi(row[24])
		faceoffsWon, _ := strconv.Atoi(row[25])
		faceoffsLost, _ := strconv.Atoi(row[26])

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

		playerName := strings.Split(row[1], "\\")[0]
		position := row[3]

		var player models.Player

		database.DB.
			Where(
				"name = ? AND team_id = ?",
				playerName,
				team.ID,
			).
			First(&player)

		if player.ID == 0 {

			player = models.Player{
				Name:     playerName,
				Position: position,
				TeamID:   team.ID,
			}

			database.DB.Create(&player)
		}

		stats := models.PlayerSeasonStats{
			PlayerID:    player.ID,
			Season:      season,
			GamesPlayed: gamesPlayed,

			Goals:          goals,
			Assists:        assists,
			PlusMinus:      plusMinus,
			PenaltyMinutes: penaltyMinutes,

			EvenStrengthGoals: evenStrengthGoals,
			PowerPlayGoals:    powerPlayGoals,
			ShortHandedGoals:  shortHandedGoals,

			Shots: shots,

			BlockedShots: blockedShots,
			Hits:         hits,

			FaceoffsWon:  faceoffsWon,
			FaceoffsLost: faceoffsLost,

			TimeOfIce: toi,
		}

		database.DB.Create(&stats)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "csv файл импортирован",
	})
}
