package student_handler

import (
	"net/http"

	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// UpsertPublications
//
//	@Summary		Добавление или обновление научных публикаций
//
//	@Description	Добавление или обновление научных публикаций
//
//	@Tags			Student.ScientificWorks
//	@Accept			json
//	@Param			input body request_models.UpsertPublicationsRequest true "Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/works/publications/{token} [post]
func (h *StudentHandler) UpsertPublications(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.UpsertPublicationsRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	if len(reqBody.Publications) == 0 {
		ctx.AbortWithStatus(http.StatusAccepted)
		return
	}

	err = h.scientific.UpsertPublications(ctx, user.KasperID, reqBody.Semester, reqBody.Publications)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
