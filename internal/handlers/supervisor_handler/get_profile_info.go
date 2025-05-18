package supervisor_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetSupervisorProfile
//
//	@Summary		Получение профиля научного руководителя
//	@Tags			Supervisor
//	@Description	Получение профиля научного руководителя
//	@Success		200	{object}	models.SupervisorProfile	"Данные"
//
//	@Produce		json
//	@Param			token	path		string	true	"Токен пользователя"
//
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/supervisors/profile/{token} [get]
func (h *SupervisorHandler) GetSupervisorProfile(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	super, err := h.supervisor.GetSupervisorProfile(ctx, user.KasperID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, super)
}
