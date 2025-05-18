package student_handler

import (
	"net/http"

	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// UpsertConferences
//
//	@Summary		Добавление или обновление научных конференций
//
//	@Description	Добавление или обновление научных конференций
//
//	@Tags			Student.ScientificWorks
//	@Accept			json
//	@Param			input body request_models.UpsertConferencesRequest true "Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/works/conferences/{token} [post]
func (h *StudentHandler) UpsertConferences(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.UpsertConferencesRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	if len(reqBody.Conferences) == 0 {
		ctx.AbortWithStatus(http.StatusAccepted)
		return
	}

	err = h.scientific.UpsertConferences(ctx, user.KasperID, reqBody.Semester, reqBody.Conferences)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
