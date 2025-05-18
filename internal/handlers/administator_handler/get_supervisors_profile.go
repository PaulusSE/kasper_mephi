package administator_handler

import (
	"net/http"

	"uir_draft/internal/handlers/administator_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetSupervisorProfile
//
//	@Summary		Получение профиля научного руководителя
//	@Tags			Admin
//	@Description	Получение профиля научного руководителя
//	@Success		200		{object}	models.SupervisorProfile			"Данные"
//
//	@Param			input	body		request_models.GetBySupervisorID	true	"Данные"
//	@Produce		json
//	@Param			token	path		string	true	"Токен пользователя"
//
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/administrator/supervisors/profile/{token} [put]
func (h *AdministratorHandler) GetSupervisorProfile(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.GetBySupervisorID{}
	if err = ctx.ShouldBind(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	resp, err := h.supervisor.GetSupervisorProfile(ctx, reqBody.SupervisorID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
