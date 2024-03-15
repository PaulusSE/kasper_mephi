package supervisor_handler

import (
	"net/http"

	"uir_draft/internal/handlers/supervisor_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// AllToStatus
//
//	@Summary		Проставление статуса для всего для аспиранта
//	@Tags			Supervisor
//	@Description	Проставление статуса для всего для аспиранта
//	@Success		200
//	@Param			token	path		string							true	"Токен пользователя"
//
//	@Param			input	body		request_models.ToStatusRequest	true	"Запрос"
//
//	@Failure		401		{string}	string							"Токен протух"
//	@Failure		204		{string}	string							"Нет записей в БД"
//	@Failure		500		{string}	string							"Ошибка на стороне сервера"
//	@Router			/supervisors/student/review/{token} [post]
func (h *SupervisorHandler) AllToStatus(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	reqBody := request_models.ToStatusRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = h.dissertation.AllToStatus(ctx, reqBody.StudentID, reqBody.Status)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	ctx.Status(http.StatusOK)
}
