package student_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetReportComments
//
//	@Summary		Получение комментариев к диссертации аспиранта
//	@Tags			NEW
//	@Description	Получение комментариев к диссертации аспиранта
//	@Success		200	{object}	models.ReportComments	"Данные"
//	@Accept			json
//
//	@Produce		json
//
//	@Param			token	path		string	true	"Токе пользователя"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/report/comments/{token} [get]
func (h *StudentHandler) GetReportComments(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	comments, err := h.report.GetReportComments(ctx, user.KasperID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, comments)
}
