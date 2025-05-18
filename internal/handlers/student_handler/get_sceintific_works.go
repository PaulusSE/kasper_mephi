package student_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetScientificWorks
//
//	@Summary		Получение данных для страницы научных работ
//
//	@Description	Получение данных для страницы научных работ
//
//	@Tags			Student.ScientificWorks
//	@Accept			json
//
//	@Produce		json
//
//	@Success		200		{object}	[]models.ScientificWork	"Данные"
//	@Param			token	path		string					true	"Токен пользователя"
//	@Failure		400		{string}	string					"Неверный формат данных"
//	@Failure		401		{string}	string					"Токен протух"
//	@Failure		204		{string}	string					"Нет записей в БД"
//	@Failure		500		{string}	string					"Ошибка на стороне сервера"
//	@Router			/students/works/{token} [get]
func (h *StudentHandler) GetScientificWorks(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	works, err := h.scientific.GetScientificWorks(ctx, user.KasperID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, works)
}
