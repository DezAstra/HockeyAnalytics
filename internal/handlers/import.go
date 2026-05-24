package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"hockeyAnalytics/internal/mappers"
	"hockeyAnalytics/internal/models"
	"hockeyAnalytics/internal/services"

	"github.com/gin-gonic/gin"
)

const (
	csvPositionIdx = 3
	csvTeamIdx     = 4
)

type CSVImportResponse struct {
	Message  string   `json:"message"`
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors"`
}

// ImportCSV godoc
// @Summary импорт csv
// @Description Импортировать статистику csv файлом
// @Tags Импорт
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV File"
// @Success 200 {object} CSVImportResponse
// @Failure 400 {object} map[string]interface{}
// @Router /import/csv [post]
func ImportCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка загрузки csv файла"})
		return
	}

	if err := os.MkdirAll("./csv", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка подготовки папки csv"})
		return
	}

	tempFilePath := filepath.Join("./csv", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения csv файла"})
		return
	}

	rows, err := services.ReadCSV(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка чтения csv файла"})
		return
	}

	imported := 0
	skipped := 0
	errorsList := make([]string, 0)

	for index, row := range rows {
		if index == 0 {
			continue
		}

		if len(row) < mappers.MinCSVColumns {
			skipped++
			errorsList = append(errorsList, "row skipped: not enough columns")
			continue
		}

		teamName := strings.TrimSpace(row[csvTeamIdx])
		playerName := mappers.ExtractCSVPlayerName(row)
		position := strings.TrimSpace(row[csvPositionIdx])

		if playerName == "" {
			skipped++
			continue
		}

		player, err := services.UpsertPlayer(models.Player{
			Name:     playerName,
			Position: position,
		})
		if err != nil {
			skipped++
			errorsList = append(errorsList, err.Error())
			continue
		}

		stats, err := mappers.MapCSVRowToStats(row, player.ID)
		if err != nil {
			skipped++
			errorsList = append(errorsList, err.Error())
			continue
		}

		stats.Team = teamName

		if err := services.SaveSeasonStats(stats); err != nil {
			skipped++
			errorsList = append(errorsList, err.Error())
			continue
		}

		imported++
	}

	analyticsEngine.ClearCache()

	status := "success"
	if len(errorsList) > 0 {
		status = "warning"
	}

	_ = services.CreateImportLog(services.ImportLogInput{
		Source:   "CSV",
		Season:   "",
		Status:   status,
		Message:  "csv файл импортирован",
		Imported: imported,
		Skipped:  skipped,
		Errors:   errorsList,
	})

	c.JSON(http.StatusOK, CSVImportResponse{
		Message:  "csv файл импортирован",
		Imported: imported,
		Skipped:  skipped,
		Errors:   errorsList,
	})
}
