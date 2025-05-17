package student_handler

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os/exec"
	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
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
		log.Printf("Authentication error: %v", err)
		ctx.AbortWithStatusJSON(models.MapErrorToCode(err), gin.H{"error": err.Error()})
		return
	}

	// Получаем параметр semester из тела запроса
	var requestData struct {
		Semester int       `json:"semester"`
		UserID   uuid.UUID `json:"student_id"`
	}

	if err := ctx.BindJSON(&requestData); err != nil {
		log.Printf("Invalid request data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID := user.KasperID
	if user.UserType != model.UserType_Student {
		userID = requestData.UserID
	}

	log.Printf("User authenticated: %s, Semester: %d", userID, requestData.Semester)

	presentationData, err := h.student.GetPresentation(ctx, userID, requestData.Semester)
	if err != nil {
		log.Printf("Error fetching presentation data: %v", err)
		ctx.AbortWithStatusJSON(models.MapErrorToCode(err), gin.H{"error": err.Error()})
		return
	}

	//load, err := h.student.GetStudentLoad(ctx, userID, int32(requestData.Semester))
	//if err != nil {
	//	log.Printf("Error fetching student load: %v", err)
	//	ctx.AbortWithStatusJSON(models.MapErrorToCode(err), gin.H{"error": err.Error()})
	//	return
	//}

	// presentationData.PedagogicalData = load

	reportDataJSON, err := json.Marshal(presentationData)
	if err != nil {
		log.Printf("Error marshaling presentation data: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare presentation data"})
		return
	}

	log.Printf("Presentation data JSON prepared successfully")

	cmd := exec.Command("python3", "generate_presentation.py")
	cmd.Stdin = bytes.NewReader(reportDataJSON)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	log.Printf("Executing Python script for presentation generation")
	if err := cmd.Run(); err != nil {
		log.Printf("Python script error: %s", stderr.String())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error running presentation script"})
		return
	}

	if out.Len() == 0 {
		log.Printf("Python script returned empty output")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Generated presentation is empty"})
		return
	}

	log.Printf("Presentation generated successfully, sending response")

	ctx.Header("Content-Disposition", "attachment; filename=report.pptx")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.presentationml.presentation", out.Bytes())
}
