package student_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetStudentStatus
//
//	@Summary		Получение данных о студенте (статус студента)
//
//	@Description	Получение данных о студенте (статус студента)
//
//	@Tags			Student
//	@Accept			json
//
//	@Produce		json
//
//	@Success		200		{object}	models.Student	"Данные"
//	@Param			token	path		string			true	"Токен пользователя"
//	@Failure		400		{string}	string			"Неверный формат данных"
//	@Failure		401		{string}	string			"Токен протух"
//	@Failure		204		{string}	string			"Нет записей в БД"
//	@Failure		500		{string}	string			"Ошибка на стороне сервера"
//	@Router			/students/info/{token} [get]
func (h *StudentHandler) GetStudentStatus(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	student, err := h.student.GetStudentStatus(ctx, user.KasperID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, student)
}
