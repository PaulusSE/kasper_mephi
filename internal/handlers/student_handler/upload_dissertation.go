package student_handler

import (
	"fmt"
	"net/http"
	"os"

	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// UploadDissertation
//
//	@Summary		Загрузка файла диссертации
//
//	@Description	Загрузка файла диссертации
//
//	@Tags			Student
//	@Accept			mpfd
//
//	@Param			file	formData	request_models.UploadDissertationRequest	true	"Данные"
//
//	@Param			file	formData	file										true	"Файл"
//
//	@Param			token	path		string										true	"Токен пользователя"
//
//	@Success		200
//
//	@Failure		400	{string}	string	"Неверный формат данных"
//	@Failure		401	{string}	string	"Токен протух"
//	@Failure		204	{string}	string	"Нет записей в БД"
//	@Failure		500	{string}	string	"Ошибка на стороне сервера"
//	@Router			/students/dissertation/file/{token} [post]
func (h *StudentHandler) UploadDissertation(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.UploadDissertationRequest{}
	if err = ctx.ShouldBind(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	dirPath := fmt.Sprintf("./dissertations/%s/semester%d", user.KasperID, reqBody.Semester)
	err = os.RemoveAll(dirPath)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	dst := fmt.Sprintf("./dissertations/%s/semester%d/%s",
		user.KasperID, reqBody.Semester, reqBody.File.Filename)

	err = ctx.SaveUploadedFile(reqBody.File, dst)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err) //nolint
		return
	}

	err = h.dissertation.UpsertDissertationInfo(ctx, user.KasperID, reqBody.Semester, reqBody.File.Filename)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Status(http.StatusOK)
}
