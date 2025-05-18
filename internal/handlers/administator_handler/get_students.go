package administator_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetStudents
//
//	@Summary		Получение списка всех аспирантов
//
//	@Description	Получение списка всех аспирантов
//
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]models.Student	"Данные"
//	@Param			token	path		string				true	"Токен пользователя"
//	@Failure		400		{string}	string				"Неверный формат данных"
//	@Failure		401		{string}	string				"Токен протух"
//	@Failure		204		{string}	string				"Нет записей в БД"
//	@Failure		500		{string}	string				"Ошибка на стороне сервера"
//	@Router			/administrator/students/list/{token} [get]
func (h *AdministratorHandler) GetStudents(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	students, err := h.user.GetStudentsList(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, students)
}
