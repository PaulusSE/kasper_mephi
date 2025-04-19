package student_handler

import (
	"log"
	"net/http"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetSupervisors
//
//	@Summary		Получения списка научных руководителей
//
//	@Description	Получения списка научных руководителей
//
//	@Tags			Student
//	@Success		200		{object}	[]models.Supervisor
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/student/supervisors/list/{token} [get]
func (h *StudentHandler) GetSupervisors(ctx *gin.Context) {
	//_, err := h.authenticate(ctx)
	//if err != nil {
	//	ctx.AbortWithError(models.MapErrorToCode(err), err)
	//	return
	//}
	if err := h.ValidateToken(ctx); err != nil {
		return
	}
	sups, err := h.admin.GetSupervisors(ctx)
	if err != nil {
		log.Printf("Error getting supervisors: %v\n", err)
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	ctx.JSON(http.StatusOK, sups)
}
