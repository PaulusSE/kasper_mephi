package supervisor_handler

import (
	"net/http"

	"uir_draft/internal/handlers/supervisor_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetStudentsReportComments
//
//	@Summary		Получение комментариев к диссертации аспиранта (от лица научника)
//	@Tags			NEW
//	@Description	Получение комментариев к диссертации аспиранта (от лица научника)
//	@Success		200	{object}	models.ReportComments	"Данные"
//	@Accept			json
//
//	@Produce		json
//
//	@Param			token	path		string							true	"Токе пользователя"
//
//	@Param			input	body		request_models.GetByStudentID	true	"Запрос"
//
//	@Failure		401		{string}	string							"Токен протух"
//	@Failure		204		{string}	string							"Нет записей в БД"
//	@Failure		500		{string}	string							"Ошибка на стороне сервера"
//	@Router			/supervisors/report/comments/{token} [put]
func (h *SupervisorHandler) GetStudentsReportComments(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.GetByStudentID{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	list, err := h.report.GetReportComments(ctx, reqBody.StudentID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, list)
}
