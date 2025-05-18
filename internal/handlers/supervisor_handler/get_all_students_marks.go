package supervisor_handler

import (
	"net/http"

	"uir_draft/internal/handlers/supervisor_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetAllMarks
//
//	@Summary		Получение всех оценок
//
//	@Description	Получение всех оценок
//
//	@Tags			NEW
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	request_models.GetByStudentID	"Данные"
//	@Param			token	path		string							true	"Токен пользователя"
//	@Failure		400		{string}	string							"Неверный формат данных"
//	@Failure		401		{string}	string							"Токен протух"
//	@Failure		204		{string}	string							"Нет записей в БД"
//	@Failure		500		{string}	string							"Ошибка на стороне сервера"
//	@Router			/supervisors/student/marks/{token} [put]
func (h *SupervisorHandler) GetAllMarks(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.GetByStudentID{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	marks, err := h.student.GetAllMarks(ctx, reqBody.StudentID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, marks)
}
