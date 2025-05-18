package student_handler

import (
	"net/http"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// AllToReview
//
//	@Summary		Отправка на проверку
//	@Tags			Student
//	@Description	Переводит научные работы, пед нагрузку и диссертацию в статус 'in_review', который является блокирующим
//	@Param			input	body	request_models.AllToReviewRequest	true	"Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/review/{token} [post]
func (h *StudentHandler) AllToReview(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		_ = ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.AllToReviewRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	err = h.dissertation.AllToStatus(ctx, user.KasperID, reqBody.Commentary, model.ApprovalStatus_OnReview.String())
	if err != nil {
		_ = ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	err = h.email.SendMailToSupervisor(ctx, user.KasperID, "./internal/templates/student.html", "Работы")
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
