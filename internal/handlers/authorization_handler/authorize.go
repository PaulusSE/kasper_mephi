package authorization_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Authorize
//
//	@Summary		Авторизация в сервисе
//
//	@Description	Авторизация в сервисе
//
//	@Tags			Authorization
//	@Accept			json
//
//	@Produce		json
//
//	@Success		200		{object}	models.AuthorizeResponse	"Данные"
//	@Param			input	body		models.AuthorizeRequest		true	"Данные"
//	@Failure		400		{string}	string						"Неверный формат данных"
//	@Failure		401		{string}	string						"Токен протух"
//	@Failure		204		{string}	string						"Нет записей в БД"
//	@Failure		500		{string}	string						"Ошибка на стороне сервера"
//
// @Router			/authorize [post]
func (h *AuthorizationHandler) Authorize(ctx *gin.Context) {
	reqBody := models.AuthorizeRequest{}
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	resp, authorized, err := h.authenticator.Authorize(ctx, reqBody)
	if err != nil {
		_ = ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}
	if !authorized {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("wrong password or email")) //nolint
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
