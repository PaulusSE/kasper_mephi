package student_handler

import (
	"net/http"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// ScientificWorksToReview
//
//	@Summary		Отправление на проверку научных работ
//
//	@Description	Отправление на проверку научных работ
//
//	@Tags			Student.ScientificWorks
//	@Accept			json
//	@Param			input body request_models.ToReviewRequest true "Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/works/review/{token} [post]
func (h *StudentHandler) ScientificWorksToReview(ctx *gin.Context) {
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

	if err = h.scientific.ScientificWorksToStatus(ctx, user.KasperID, model.ApprovalStatus_OnReview, reqBody.Semester); err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	//err = h.email.SendMailToSupervisor(ctx, user.KasperID, "path", "Научные работы")
	//if err != nil {
	//	ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
	//	return
	//}

	ctx.Status(http.StatusOK)
}
