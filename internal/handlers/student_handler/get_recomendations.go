package student_handler

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Recommendation представляет структуру рекомендации
type Recommendation struct {
	Title string `json:"title"`
	Link  string `json:"url"`
}

// GetRecommendedArticles
//
//	@Summary		Получение рекомендованных статей
//	@Description	Получение списка рекомендованных статей (пока случайный пример)
//	@Tags			Supervisor.RecommendedArticles
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]RecommendedArticle	"Данные"
//	@Param			token	path		string					true	"Токен пользователя"
//	@Failure		400		{string}	string					"Неверный формат данных"
//	@Failure		401		{string}	string					"Токен протух"
//	@Failure		500		{string}	string					"Ошибка на стороне сервера"
//	@Router			/supervisors/recommended-articles/{token} [get]
func (h *StudentHandler) GetRecommendedArticles(ctx *gin.Context) {
	// Путь к Python-скрипту, который возвращает JSON
	pythonScript := "recomendations.py"

	// Запускаем скрипт, передаем text как аргумент
	// cmd := exec.Command("python3", pythonScript, text)

	cmd := exec.Command("python3", pythonScript)
	log, err := cmd.CombinedOutput()
	if err != nil {
		// log содержит и stdout, и stderr
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Script error: %v\nOutput:\n%s", err, string(log)),
		})
		return
	}

	// output, err := cmd.Output()
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": fmt.Sprintf("Failed to execute script: %v", err),
	// 	})
	// 	return
	// }

	// Парсим JSON
	var recommendations []Recommendation
	if err := json.Unmarshal(log, &recommendations); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to parse JSON: %v. Script output: %s", err, string(log)),
		})
		return
	}

	// Отдаём результат
	ctx.JSON(http.StatusOK, recommendations)
}
