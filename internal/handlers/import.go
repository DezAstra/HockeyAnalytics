package handlers

import (
	"hockeyAnalytics/internal/database"
	"hockeyAnalytics/internal/mappers"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/services"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	csvPlayerIdx   = 1
	csvPositionIdx = 3
	csvTeamIdx     = 4
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

		if len(row) <= csvTeamIdx {
			continue
		}

		teamName := row[csvTeamIdx]
		playerName := strings.Split(row[csvPlayerIdx], "\\")[0]
		position := row[csvPositionIdx]

		var player models.Player

		database.DB.
			Where(
				"LOWER(name) = LOWER(?)",
				playerName,
			).
			First(&player)

		if player.ID == 0 {

			player = models.Player{
				Name:     playerName,
				Position: position,
			}

			database.DB.Create(&player)
		}

		stats :=
			mappers.MapCSVRowToStats(
				row,
				player.ID,
			)

		stats.Team = teamName

		services.SaveSeasonStats(stats)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "csv файл импортирован",
	})
}
