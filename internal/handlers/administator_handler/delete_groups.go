package administator_handler

import (
	"net/http"

	"uir_draft/internal/handlers/administator_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// DeleteGroups
//
//	@Summary		Удаление групп
//
//	@Description	Удаление групп
//
//	@Tags			Admin
//	@Accept			json
//	@Param			input body request_models.DeleteEnumRequest true "Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/administrator/enum/groups/{token} [put]
func (h *AdministratorHandler) DeleteGroups(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.DeleteEnumRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	if len(reqBody.IDs) == 0 {
		ctx.AbortWithStatus(http.StatusAccepted)
		return
	}

	if err = h.enum.DeleteGroups(ctx, reqBody.IDs); err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
