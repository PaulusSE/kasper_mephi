package student_handler

import (
	"net/http"

	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// DeleteClassroomLoads
//
//	@Summary		Удаление аудиторной педнагрузки
//	@Description	Удаление аудиторной нагрузки по ID нагрузки
//	@Tags			Student.TeachingLoad
//	@Accept			json
//
//	@Param			input	body	request_models.DeleteIDs	true	"ID нагрузок и семестр"
//
//	@Success		200
//	@Param			token	path		string	true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/load/classroom/{token} [put]
func (h *StudentHandler) DeleteClassroomLoads(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	reqBody := request_models.DeleteIDs{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if len(reqBody.IDs) == 0 {
		ctx.AbortWithStatus(http.StatusCreated)
		return
	}

	err = h.load.DeleteClassroomLoad(ctx, user.KasperID, reqBody.Semester, reqBody.IDs)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	ctx.Status(http.StatusOK)
}
