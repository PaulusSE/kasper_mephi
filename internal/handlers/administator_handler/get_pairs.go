package administator_handler

import (
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetPairs
//
//	@Summary		Получение пар аспирант/руководитель
//
//	@Description	Ручка для получения пар аспирант/руководитель
//
//	@Tags			Admin
//	@Success		200		{object}	[]models.StudentSupervisorPair
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/administrator/pairs/{token} [get]
func (h *AdministratorHandler) GetPairs(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	pairs, err := h.user.GetStudentSupervisorPairs(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.JSON(http.StatusOK, pairs)
}
