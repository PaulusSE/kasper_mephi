package student_handler

import (
	"net/http"

	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// UpsertAdditionalLoads
//
//	@Summary		Добавление или обновление дополнительной нагрузки
//
//	@Description	Добавление или обновление дополнительной нагрузки
//
//	@Tags			Student.TeachingLoad
//	@Accept			json
//	@Param			input body request_models.UpsertAdditionalLoadRequest true "Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/load/additional/{token} [post]
func (h *StudentHandler) UpsertAdditionalLoads(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.UpsertAdditionalLoadRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	if len(reqBody.Loads) == 0 {
		ctx.AbortWithStatus(http.StatusAccepted)
		return
	}

	err = h.load.UpsertAdditionalLoad(ctx, user.KasperID, reqBody.Semester, reqBody.Loads)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
