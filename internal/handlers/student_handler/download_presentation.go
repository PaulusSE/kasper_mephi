package student_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetPresentation
//
//	@Summary		Загрузка презентации
//	@Description	Генерация и загрузка презентации для студента
//	@Tags			Student.Presentation
//	@Accept			json
//	@Produce		application/vnd.openxmlformats-officedocument.presentationml.presentation
//	@Success		200		{file}		file	"Презентация"
//	@Param			token		path	string	true	"Токен пользователя"
//	@Param			semester	query	int		true	"Семестр"
//	@Failure		400			{string}	string	"Неверный формат данных"
//	@Failure		401			{string}	string	"Токен протух"
//	@Failure		500			{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/report/download/{token} [post]
func (h *StudentHandler) GetPresentation(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(models.MapErrorToCode(err), gin.H{"error": err.Error()})
		return
	}

	// Получаем параметр semester из тела запроса
	var requestData struct {
		Semester int `json:"semester"`
	}

	if err := ctx.BindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	presentationData, err := h.student.GetPresentation(ctx, user.KasperID)
	if err != nil {
		ctx.AbortWithStatusJSON(models.MapErrorToCode(err), gin.H{"error": err.Error()})
		return
	}

	load, err := h.student.GetStudentLoad(ctx, user.KasperID, int32(requestData.Semester))
	if err != nil {
		ctx.AbortWithStatusJSON(models.MapErrorToCode(err), gin.H{"error": err.Error()})
		return
	}

	presentationData.PedagogicalData = load

	reportDataJSON, err := json.Marshal(presentationData)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Вызов Python скрипта для генерации презентации
	cmd := exec.Command("python3", "generate_presentation.py")
	cmd.Stdin = bytes.NewReader(reportDataJSON)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Отправляем файл пользователю
	ctx.Header("Content-Disposition", "attachment; filename=report.pptx")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.presentationml.presentation", out.Bytes())
}
