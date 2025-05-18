package administator_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetSemesterAmounts
//
//	@Summary		Добавление количеств семестров
//
//	@Description	Добавление количеств семестров
//
//	@Tags			NEW
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]models.SemesterAmount	"Данные"
//	@Param			token	path		string					true	"Токен пользователя"
//	@Failure		400		{string}	string					"Неверный формат данных"
//	@Failure		401		{string}	string					"Токен протух"
//	@Failure		204		{string}	string					"Нет записей в БД"
//	@Failure		500		{string}	string					"Ошибка на стороне сервера"
//	@Router			/administrator/enum/amounts/{token} [get]
func (h *AdministratorHandler) GetSemesterAmounts(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	amounts, err := h.enum.GetSemestersAmount(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, amounts)
}
