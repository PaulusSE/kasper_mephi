package supervisor_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetStudentsList
//
//	@Summary		Получение списка аспирантов научного руководителя
//	@Tags			Supervisor
//	@Description	Получение списка аспирантов научного руководителя
//	@Success		200	{object}	[]models.Student	"Данные"
//
//	@Produce		json
//	@Param			token	path		string	true	"Токен пользователя"
//
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/supervisors/student/list/{token} [put]
func (h *SupervisorHandler) GetStudentsList(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	list, err := h.supervisor.GetStudentList(ctx, user.KasperID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, list)
}
