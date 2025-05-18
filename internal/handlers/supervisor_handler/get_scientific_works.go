package supervisor_handler

import (
	"net/http"

	"uir_draft/internal/handlers/supervisor_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetScientificWorks
//
//	@Summary		Получение страницы информации для страницы научных работ аспиранта
//	@Tags			Supervisor.ScientificWorks
//	@Description	Получение страницы информации для страницы научных работ аспиранта
//	@Success		200	{object}	[]models.ScientificWork	"Данные"
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
//	@Router			/supervisors/student/works/{token} [put]
func (h *SupervisorHandler) GetScientificWorks(ctx *gin.Context) {
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

	works, err := h.scientific.GetScientificWorks(ctx, reqBody.StudentID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, works)
}
