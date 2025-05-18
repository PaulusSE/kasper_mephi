package authorization_handler

import (
	"net/http"

	"uir_draft/internal/pkg/helpers"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

type TokenCheckResponse struct {
	UserType   string `json:"user_type" enums:"admin,student,supervisor"`
	Registered bool   `json:"registered"`
}

// TokenCheck
//
//	@Summary		Проверка токена пользователя
//
//	@Description	Проверка токена пользователя
//
//	@Tags			Authorization
//
//	@Produce		json
//
//	@Success		200		{object}	TokenCheckResponse	"Данные"
//	@Param			token	path		string				true	"Токен пользователя"
//	@Failure		400		{string}	string				"Неверный формат данных"
//	@Failure		401		{string}	string				"Токен протух"
//	@Failure		204		{string}	string				"Нет записей в БД"
//	@Failure		500		{string}	string				"Ошибка на стороне сервера"
//	@Router			/authorize/token/check/{token} [get]
func (h *AuthorizationHandler) TokenCheck(ctx *gin.Context) {
	token := helpers.GetToken(ctx)

	user, err := h.authenticator.TokenCheck(ctx, token)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	resp := TokenCheckResponse{
		UserType:   user.UserType.String(),
		Registered: user.Registered,
	}

	ctx.JSON(http.StatusOK, resp)
}
