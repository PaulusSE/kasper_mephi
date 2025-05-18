package supervisor_handler

import (
	"net/http"

	"uir_draft/internal/handlers/supervisor_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetStudentStatus
//
//	@Summary		Получение данных о студенте (статус студента)
//
//	@Description	Получение данных о студенте (статус студента)
//
//	@Tags			Supervisor
//	@Accept			json
//
//	@Produce		json
//	@Param			input	body		request_models.GetByStudentID	true	"Запрос"
//	@Success		200		{object}	models.Student					"Данные"
//	@Param			token	path		string							true	"Токен пользователя"
//	@Failure		400		{string}	string							"Неверный формат данных"
//	@Failure		401		{string}	string							"Токен протух"
//	@Failure		204		{string}	string							"Нет записей в БД"
//	@Failure		500		{string}	string							"Ошибка на стороне сервера"
//	@Router			/supervisors/student/info/{token} [put]
func (h *SupervisorHandler) GetStudentStatus(ctx *gin.Context) {
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

	student, err := h.student.GetStudentStatus(ctx, reqBody.StudentID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, student)
}
