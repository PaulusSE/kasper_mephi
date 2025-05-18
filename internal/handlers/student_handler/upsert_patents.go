package student_handler

import (
	"net/http"

	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// UpsertPatents
//
//	@Summary		Добавление или обновление патентов
//
//	@Description	Добавление или обновление патентов
//
//	@Tags			NEW
//	@Accept			json
//	@Param			input body request_models.UpsertPatentsRequest true "Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/works/patents/{token} [post]
func (h *StudentHandler) UpsertPatents(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.UpsertPatentsRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	if len(reqBody.Patents) == 0 {
		ctx.AbortWithStatus(http.StatusAccepted)
		return
	}

	err = h.scientific.UpsertPatents(ctx, user.KasperID, reqBody.Semester, reqBody.Patents)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
