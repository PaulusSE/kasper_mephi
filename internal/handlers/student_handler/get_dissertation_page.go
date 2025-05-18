package student_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetDissertationPage
//
//	@Summary		Получение данных для страницы диссертации
//
//	@Description	Получение данных для страницы диссертации
//
//	@Tags			Student.Dissertation
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	models.DissertationPageResponse	"Данные"
//	@Param			token	path		string							true	"Токен пользователя"
//	@Failure		400		{string}	string							"Неверный формат данных"
//	@Failure		401		{string}	string							"Токен протух"
//	@Failure		204		{string}	string							"Нет записей в БД"
//	@Failure		500		{string}	string							"Ошибка на стороне сервера"
//	@Router			/students/dissertation/{token} [get]
func (h *StudentHandler) GetDissertationPage(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	dissertationPage, err := h.dissertation.GetDissertationPage(ctx, user.KasperID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, dissertationPage)
}
