package administator_handler

import (
	"net/http"

	"uir_draft/internal/handlers/administator_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// UpsertAttestationMarks
//
//	@Summary		Апсерт оценок за аттестацию (от комиссии)
//
//	@Description	Апсерт оценок за аттестацию (от комиссии)
//
//	@Tags			NEW
//	@Accept			json
//	@Param			input body request_models.UpsertAttestationMarksRequest true "Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/administrator/student/attestation/marks/{token} [post]
func (h *AdministratorHandler) UpsertAttestationMarks(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.UpsertAttestationMarksRequest{}
	if err = ctx.ShouldBind(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	if err = h.user.UpsertAttestationMarks(ctx, reqBody.AttestationMarks); err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
