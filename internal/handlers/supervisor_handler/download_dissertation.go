package supervisor_handler

import (
	"fmt"
	"net/http"
	"os"

	"uir_draft/internal/handlers/supervisor_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

// DownloadDissertation
//
//	@Summary		Скачивание файла диссертации
//
//	@Description	Скачивание файла диссертации
//
//	@Tags			Supervisor.Dissertation
//	@Accept			json
//
//	@Produce		json
//
//	@Success		200		"Файл"
//	@Param			token	path		string											true	"Токен пользователя"
//	@Param			input	body		request_models.DownloadDissertationRequestSup	true	"Данные"
//	@Failure		400		{string}	string											"Неверный формат данных"
//	@Failure		401		{string}	string											"Токен протух"
//	@Failure		204		{string}	string											"Нет записей в БД"
//	@Failure		500		{string}	string											"Ошибка на стороне сервера"
//	@Router			/supervisors/student/dissertation/file/{token} [put]
func (h *SupervisorHandler) DownloadDissertation(ctx *gin.Context) {
	_, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	reqBody := request_models.DownloadDissertationRequestSup{}
	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err) //nolint
		return
	}

	dissertation, err := h.dissertation.GetDissertationData(ctx, reqBody.StudentID, reqBody.Semester)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	if dissertation.FileName == nil {
		ctx.Status(http.StatusNoContent)
		return
	}

	dst := fmt.Sprintf("./dissertations/%s/semester%d/%s",
		dissertation.StudentID.String(), dissertation.Semester, lo.FromPtr(dissertation.FileName))

	_, err = os.Stat(dst)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err) //nolint
		return
	}

	ctx.Header("Content-Disposition", lo.FromPtr(dissertation.FileName))
	ctx.File(dst)
}
