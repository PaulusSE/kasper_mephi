package supervisor_handler

import (
	"net/http"

	"uir_draft/internal/handlers/supervisor_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// UpsertFeedback
//
//	@Summary		Отправка комментария и оценки от научника
//	@Tags			Supervisor.Dissertation
//	@Description	Отправка комментария и оценки от научника
//	@Success		200
//	@Param			token	path	string	true	"Токен пользователя"
//	@Accept			json
//
//	@Produce		json
//
//	@Param			input	body		request_models.UpsertFeedbackRequest	true	"Запрос"
//
//	@Failure		401		{string}	string									"Токен протух"
//	@Failure		204		{string}	string									"Нет записей в БД"
//	@Failure		500		{string}	string									"Ошибка на стороне сервера"
//	@Router			/supervisors/student/feedback/{token} [post]
func (h *SupervisorHandler) UpsertFeedback(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.UpsertFeedbackRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	err = h.supervisor.UpsertFeedback(ctx, reqBody.StudentID, user.KasperID, reqBody.Feedback)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
