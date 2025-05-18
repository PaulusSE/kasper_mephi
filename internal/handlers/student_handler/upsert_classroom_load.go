package student_handler

import (
	"net/http"

	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// UpsertClassroomLoads
//
//	@Summary		Добавление или обновление аудиторной нагрузки
//
//	@Description	Добавление или обновление аудиторной нагрузки
//
//	@Tags			Student.TeachingLoad
//	@Accept			json
//	@Param			input body request_models.UpsertClassroomLoadRequest true "Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/load/classroom/{token} [post]
func (h *StudentHandler) UpsertClassroomLoads(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.UpsertClassroomLoadRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	if len(reqBody.Loads) == 0 {
		ctx.AbortWithStatus(http.StatusAccepted)
		return
	}

	err = h.load.UpsertClassroomLoad(ctx, user.KasperID, reqBody.Semester, reqBody.Loads)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
