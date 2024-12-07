package handlers

import (
	"github.com/gin-gonic/gin"
	"music-library/databases"
	"music-library/pkg/logging"
	"music-library/pkg/models"
	"music-library/pkg/services"
	"net/http"
	"strconv"
	"strings"
)

// GetSongs возвращает список песен
// @Summary Получить песни
// @Description Возвращает список песен с фильтрацией по полям и пагинацией
// @Tags Songs
// @Accept json
// @Produce json
// @Param group query string false "Фильтр по группе"
// @Param song query string false "Фильтр по названию песни"
// @Param releaseDate query string false "Фильтр по дате релиза"
// @Param page query int false "Номер страницы"
// @Param size query int false "Размер страницы"
// @Success 200 {array} models.Song
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	var songs []models.Song
	group := c.Query("group")
	title := c.Query("title")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	query := databases.DB
	if group != "" {
		query = query.Where("group LIKE ?", "%"+group+"%")
	}

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&songs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список песен"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": songs})
}

// GetSongText возвращает текст песни с пагинацией по куплетам
// @Summary Получить текст песни
// @Description Возвращает текст песни, разбитый на куплеты, с поддержкой пагинации
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы"
// @Param size query int false "Размер страницы"
// @Success 200 {object} models.SongLyricsResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /songs/{id}/lyrics [get]
func GetSongText(c *gin.Context) {
	var song models.Song
	id := c.Param("id")
	if err := databases.DB.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		return
	}

	couplets := strings.Split(song.Text, "\n\n")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1"))

	start := (page - 1) * limit
	if start >= len(couplets) {
		c.JSON(http.StatusOK, gin.H{"data": []string{}})
		return
	}

	end := start + limit
	if end > len(couplets) {
		end = len(couplets)
	}

	c.JSON(http.StatusOK, gin.H{"data": couplets[start:end]})
}

// AddSong добавляет новую песню в библиотеку
// @Summary Добавить песню
// @Description Добавляет новую песню в библиотеку
// @Tags Songs
// @Accept json
// @Produce json
// @Param song body models.AddSongRequest true "Группа и название песни"
// @Success 201 {object} models.Song
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /songs [post]
func AddSong(c *gin.Context) {
	var newSong struct {
		Title string `json:"title" binding:"required"`
		Group string `json:"group" binding:"required"`
	}
	if err := c.ShouldBindJSON(&newSong); err != nil {
		logging.Log.Warn("Неверный формат данных в запросе на добавление песни")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	logging.Log.Debugf("Попытка добавить песню: группа=%s, песня=%s", newSong.Group, newSong.Title)

	song := models.Song{
		Group: newSong.Group,
		Title: newSong.Title,
	}

	details, err := services.FetchSongDetails(newSong.Group, newSong.Title)
	if err != nil {
		logging.Log.Info("Успешно получили данные из внешнего API")
		song.Text = details.Text
		song.FilePath = details.FilePath
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось получить данные о песне с API",
		})
		return
	}

	if err := databases.DB.Create(&song).Error; err != nil {
		logging.Log.Error("Ошибка сохранения песни в базу данных")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить песню"})
		return
	}
	logging.Log.Infof("Песня успешно добавлена: ID=%d", song.ID)
	c.JSON(http.StatusCreated, gin.H{"data": song})
}

// DeleteSong удаляет песню по ID
// @Summary Удалить песню
// @Description Удаляет песню из библиотеки по указанному ID
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} models.MessageResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	var song models.Song
	id := c.Param("id")
	if err := databases.DB.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		return
	}

	databases.DB.Delete(&song)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

// UpdateSong обновляет данные песни
// @Summary Обновить песню
// @Description Обновляет данные песни по ID
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.UpdateSongRequest true "Данные для обновления"
// @Success 200 {object} models.Song
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	var song models.Song
	id := c.Param("id")
	logging.Log.Debugf("Попытка обновить песню: ID=%s", id)
	if err := databases.DB.First(&song, id).Error; err != nil {
		logging.Log.Warnf("Песня с ID=%s не найдена", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		return
	}

	var updatedSong struct {
		Group    *string `json:"group,omitempty"`
		Title    *string `json:"title,omitempty"`
		Text     *string `json:"text,omitempty"`
		FilePath *string `json:"filePath,omitempty"`
	}
	if err := c.ShouldBindJSON(&updatedSong); err != nil {
		logging.Log.Warnf("Неверный формат данных при обновлении песни: ID=%s", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	logging.Log.Debugf("Данные для обновления: %+v", updatedSong)

	if updatedSong.Group != nil {
		song.Group = *updatedSong.Group
	}
	if updatedSong.Title != nil {
		song.Title = *updatedSong.Title
	}
	if updatedSong.Text != nil {
		song.Text = *updatedSong.Text
	}
	if updatedSong.FilePath != nil {
		song.FilePath = *updatedSong.FilePath
	}

	if err := databases.DB.Save(&song).Error; err != nil {
		logging.Log.Error("Ошибка обновления данных песни")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить песню"})
		return
	}
	logging.Log.Infof("Песня успешно обновлена: ID=%s", id)
	c.JSON(http.StatusOK, gin.H{"data": song})
}
