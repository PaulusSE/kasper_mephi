package administator_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetSupervisors
//
//	@Summary		Получения списка научных руководителей
//
//	@Description	Получения списка научных руководителей
//
//	@Tags			Admin
//	@Success		200		{object}	[]models.Supervisor
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/administrator/supervisors/list/{token} [get]
func (h *AdministratorHandler) GetSupervisors(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	sups, err := h.user.GetSupervisors(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, sups)
}
