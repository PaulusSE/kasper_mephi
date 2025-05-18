package supervisor_handler

import (
	"net/http"

	"uir_draft/internal/handlers/supervisor_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetTeachingLoad
//
//	@Summary		Получение страницы информации для страницы научных работ аспиранта
//	@Tags			Supervisor.TeachingLoad
//	@Description	Получение страницы информации для страницы научных работ аспиранта
//	@Success		200		{object}	[]models.TeachingLoad	"Данные"
//	@Param			token	path		string					true	"Токен пользователя"
//	@Accept			json
//
//	@Produce		json
//
//	@Param			input	body		request_models.GetByStudentID	true	"Запрос"
//
//	@Failure		401		{string}	string							"Токен протух"
//	@Failure		204		{string}	string							"Нет записей в БД"
//	@Failure		500		{string}	string							"Ошибка на стороне сервера"
//	@Router			/supervisors/student/load/{token} [put]
func (h *SupervisorHandler) GetTeachingLoad(ctx *gin.Context) {
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

	loads, err := h.load.GetTeachingLoad(ctx, reqBody.StudentID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, loads)
}
