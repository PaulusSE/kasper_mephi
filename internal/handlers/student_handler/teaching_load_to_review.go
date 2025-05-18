package student_handler

import (
	"net/http"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// TeachingLoadToReview
//
//	@Summary		Отправление на проверку пед нагрузки
//
//	@Description	Отправление на проверку пед нагрузки
//
//	@Tags			Student.TeachingLoad
//	@Accept			json
//	@Param			input body request_models.ToReviewRequest true "Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/load/review/{token} [post]
func (h *StudentHandler) TeachingLoadToReview(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.ToReviewRequest{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	if err = h.load.TeachingLoadToStatus(ctx, user.KasperID, model.ApprovalStatus_OnReview, reqBody.Semester); err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	//err = h.email.SendMailToSupervisor(ctx, user.KasperID, "path", "Педагогическая нагрузка")
	//if err != nil {
	//	ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
	//	return
	//}

	ctx.Status(http.StatusOK)
}
