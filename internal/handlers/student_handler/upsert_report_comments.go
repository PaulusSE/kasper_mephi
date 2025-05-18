package student_handler

import (
	"net/http"

	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// UpsertReportComments
//
//	@Summary		Внесение комментариев для диссертации аспиранта
//	@Tags			NEW
//	@Description	Внесение комментариев для диссертации аспиранта
//	@Accept			json
//
//	@Produce		json
//
//	@Param			token	path		string										true	"Токе пользователя"
//
//	@Param			input	body		request_models.UpsertReportCommentsRequest	true	"Запрос"
//
//	@Failure		401		{string}	string										"Токен протух"
//	@Failure		204		{string}	string										"Нет записей в БД"
//	@Failure		500		{string}	string										"Ошибка на стороне сервера"
//	@Router			/students/report/comments/{token} [post]
func (h *StudentHandler) UpsertReportComments(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.UpsertReportCommentsRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	err = h.report.UpsertReportComments(ctx, user.KasperID, reqBody)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
