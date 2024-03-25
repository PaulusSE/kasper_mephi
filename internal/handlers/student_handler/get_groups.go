package student_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetGroups
//
//	@Summary		Получение списка всех групп
//
//	@Description	Получение списка всех групп
//
//	@Tags			Student.Enum
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]models.Group	"Данные"
//	@Param			token	path		string			true	"Токен пользователя"
//	@Failure		400		{string}	string			"Неверный формат данных"
//	@Failure		401		{string}	string			"Токен протух"
//	@Failure		204		{string}	string			"Нет записей в БД"
//	@Failure		500		{string}	string			"Ошибка на стороне сервера"
//	@Router			/student/enum/groups/{token} [get]
func (h *StudentHandler) GetGroups(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	groups, err := h.enum.GetGroups(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, groups)
}
