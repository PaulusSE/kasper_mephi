package student_handler

import (
	"net/http"

	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// UpsertExamResults
//
//	@Summary		Внесение оценок за кандидатские экзамены
//
//	@Description	Внесение оценок за кандидатские экзамены
//
//	@Tags			NEW
//	@Accept			json
//	@Param			input body request_models.UpsertExamResultsRequest true "Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/exams/{token} [post]
func (h *StudentHandler) UpsertExamResults(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	reqBody := request_models.UpsertExamResultsRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err = h.mark.UpsertExamResults(
		ctx,
		user.KasperID,
		reqBody.Marks,
	); err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	ctx.Status(http.StatusOK)
}
