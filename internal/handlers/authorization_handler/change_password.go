package authorization_handler

import (
	"net/http"

	"uir_draft/internal/handlers/authorization_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// TODO

// ChangePassword
//
//	@Summary		Изменение пароля пользователя
//
//	@Description	Изменение пароля пользователя
//
//	@Tags			Authorization
//	@Accept			json
//	@Param			input	body	request_models.ChangePasswordRequest	true	"Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/authorize/password/change/{token} [post]
func (h *AuthorizationHandler) ChangePassword(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		_ = ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.ChangePasswordRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	if err = h.authenticator.ChangePassword(ctx, user.UserID, reqBody); err != nil {
		_ = ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
