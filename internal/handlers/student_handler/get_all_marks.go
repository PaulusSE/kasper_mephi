package student_handler

import (
	"net/http"

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
//	@Success		200		{object}	models.AllMarks	"Данные"
//	@Param			token	path		string			true	"Токен пользователя"
//	@Failure		400		{string}	string			"Неверный формат данных"
//	@Failure		401		{string}	string			"Токен протух"
//	@Failure		204		{string}	string			"Нет записей в БД"
//	@Failure		500		{string}	string			"Ошибка на стороне сервера"
//	@Router			/students/marks/{token} [get]
func (h *StudentHandler) GetAllMarks(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	marks, err := h.mark.GetAllMarks(ctx, user.KasperID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, marks)
}
