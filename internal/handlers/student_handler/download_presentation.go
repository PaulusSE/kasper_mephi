package student_handler

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/exec"

// 	"uir_draft/internal/handlers/student_handler/request_models"
// 	"uir_draft/internal/pkg/models"

// 	"github.com/gin-gonic/gin"
// 	"github.com/samber/lo"
// )

// // DownloadDissertation
// //
// //	@Summary		Скачивание файла презентации
// //
// //	@Description	Скачивание файла презентации
// //
// //	@Tags			Student
// //	@Accept			json
// //
// //	@Produce		json
// //
// //	@Success		200		"Файл"
// //	@Param			token	path		string										true	"Токен пользователя"
// //	@Param			input	body		request_models.DownloadPresentationRequest	true	"Данные"
// //	@Failure		400		{string}	string										"Неверный формат данных"
// //	@Failure		401		{string}	string										"Токен протух"
// //	@Failure		204		{string}	string										"Нет записей в БД"
// //	@Failure		500		{string}	string										"Ошибка на стороне сервера"
// //	@Router			/students/dissertation/file/{token} [put]
// func (h *StudentHandler) DownloadPresentation(ctx *gin.Context) {
// 	user, err := h.authenticate(ctx)
// 	if err != nil {
// 		ctx.AbortWithError(models.MapErrorToCode(err), err)
// 		return
// 	}

// 	reqBody := request_models.DownloadPresentationRequest{}
// 	if err = ctx.ShouldBindJSON(&reqBody); err != nil {
// 		ctx.AbortWithError(http.StatusBadRequest, err)
// 		return
// 	}

// 	presentation, err := h.dissertation.GetDissertationData(ctx, user.KasperID, reqBody.Semester)
// 	if err != nil {
// 		ctx.AbortWithError(models.MapErrorToCode(err), err)
// 		return
// 	}

// 	log.Printf("download dissertation info: %v", presentation)
// 	log.Printf("download dissertation file_name: %v", lo.FromPtr(presentation.FileName))

// 	if dissertation.FileName == nil {
// 		ctx.Status(http.StatusNoContent)
// 		return
// 	}

// 	data := map[string]interface{}{
// 		"full_name":       "Евченко Игорь Владимирович",
// 		"supervisor_name": "Тихомирова Анна Николаевна",
// 		"slides_data": []map[string]interface{}{
// 			{
// 				"type":      "text",
// 				"content":   "Пример текста",
// 				"left":      50,       // Pt(50) в Python будет умножено на Pt
// 				"top":       50,       // Pt(50)
// 				"width":     620,      // Pt(620)
// 				"height":    60,       // Pt(60)
// 				"font_size": 12,       // Размер шрифта
// 				"bold":      false,    // Жирный шрифт
// 				"alignment": "CENTER", // Выравнивание
// 			},
// 		},
// 	}

// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		log.Fatalf("Ошибка сериализации данных: %s", err)
// 	}

// 	cmd := exec.Command("python", "script.py")
// 	cmd.Stdin = bytes.NewReader(jsonData)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	if err := cmd.Run(); err != nil {
// 		log.Fatalf("Ошибка выполнения скрипта: %s", err)
// 	}

// 	log.Println("Python скрипт успешно выполнен.")

// 	dst := fmt.Sprintf("./dissertations/%s/semester%d/%s",
// 		dissertation.StudentID.String(), dissertation.Semester, lo.FromPtr(dissertation.FileName))

// 	_, err = os.Stat(dst)
// 	if err != nil {
// 		ctx.AbortWithError(models.MapErrorToCode(err), err)
// 		return
// 	}

// 	ctx.Header("Content-Disposition", lo.FromPtr(dissertation.FileName))
// 	ctx.File(dst)
// }
