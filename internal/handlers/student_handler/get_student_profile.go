package student_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetStudentProfile
//
//	@Summary		Получение профиля студента
//
//	@Description	Получение профиля студента
//
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	models.StudentProfile	"Данные"
//	@Param			token	path		string					true	"Токен пользователя"
//	@Failure		400		{string}	string					"Неверный формат данных"
//	@Failure		401		{string}	string					"Токен протух"
//	@Failure		204		{string}	string					"Нет записей в БД"
//	@Failure		500		{string}	string					"Ошибка на стороне сервера"
//	@Router			/student/profile/{token} [get]
func (h *StudentHandler) GetStudentProfile(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	student, err := h.student.GetStudentsProfile(ctx, user.KasperID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}
	ctx.JSON(http.StatusOK, student)
}
