package administator_handler

import (
	"net/http"

	"uir_draft/internal/handlers/administator_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// SetStudentStudyingStatus
//
//	@Summary		Изменение статуса обучения у аспиранта
//
//	@Description	Изменение статуса обучения у аспиранта
//
//	@Tags			Admin
//	@Accept			json
//	@Param			input body request_models.SetStudentFlagsRequest true "Данные"
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/administrator/student/status/{token} [post]
func (h *AdministratorHandler) SetStudentStudyingStatus(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.SetStudentFlagsRequest{}
	if err = ctx.ShouldBind(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	if err = h.user.SetStudentFlags(ctx, reqBody.Students); err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
