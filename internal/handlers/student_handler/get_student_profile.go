package student_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"os/exec"
	"time"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetStudentProfile
//
//	@Summary		Получение профиля студента
//	@Description	Получение профиля студента
//	@Tags			Student
//	@Accept		json
//	@Produce		json
//	@Success		200		{object}		models.Student
//	@Param		token	path		string		true	"Токен пользователя"
//	@Failure		400		{string}	string	"Неверный формат данных"
//	@Failure		401		{string}	string	"Токен протух"
//	@Failure		204		{string}	string	"Нет записей в БД"
//	@Failure		500		{string}	string	"Ошибка на стороне сервера"
//	@Router			/student/profile/{token} [get]

func (h *StudentHandler) GetStudentProfile(ctx *gin.Context) {
	log.Println("Начало обработки запроса для получения профиля студента")

	user, err := h.authenticate(ctx)
	if err != nil {
		log.Printf("Ошибка аутентификации: %v\n", err)
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	student, err := h.student.GetStudentsProfile(ctx, user.KasperID)
	if err != nil {
		log.Printf("Ошибка получения профиля студента: %v\n", err)
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	// err = h.handleWithPptx(ctx, student)
	// if err != nil {
	// 	log.Printf("Ошибка при обработке PPTX: %v\n", err)
	// 	// Ошибка обрабатывается, но клиенту всё равно отправляется профиль
	// }

	ctx.JSON(http.StatusOK, student)
}

func (h *StudentHandler) handleWithPptx(ctx *gin.Context, student models.StudentProfile) error {
	load, err := h.student.GetLoad(ctx, student.StudentID, student.ActualSemester)
	if err != nil {
		log.Printf("Ошибка получения учебной нагрузки: %v\n", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		load = []models.PedagogicalWork{}
	}

	log.Println("Учебная нагрузка получена")

	// reportData := prepareReportData(student)
	reportData := models.ReportData{
		CurrentSemester:      student.ActualSemester,
		FullName:             student.FullName,
		SupervisorName:       "Supervisor Name",
		EducationDirection:   "Education Direction",
		EducationProfile:     student.Specialization,
		EnrollmentDate:       student.StartDate,
		Specialty:            student.Specialization,
		TrainingYearFGOS:     "Training Year FGOS",
		CandidateExams:       []models.Exam{},
		Category:             student.Category,
		Topic:                "Dissertation Topic",
		ReportPeriodWork:     "Report Period Work",
		ScientificObject:     "Scientific Object",
		ScientificSubject:    "Scientific Subject",
		MentorRate:           "Mentor Rate",
		ProgressPercents:     []int{},
		ProgressDescriptions: []string{},
		Publications:         []models.Publication{},
		AllPublications:      []models.Publication{},
		PedagogicalData:      load,
		NextSemesterPlan:     []string{"Plan 1", "Plan 2"},
	}

	reportDataJSON, err := json.Marshal(reportData)
	if err != nil {
		return err
	}

	if err := generatePptx(reportDataJSON); err != nil {
		return err
	}

	// Подготовка multipart ответа
	w := multipart.NewWriter(ctx.Writer)
	defer w.Close()
	ctx.Header("Content-Type", w.FormDataContentType())

	// Добавление JSON части
	jsonPart, err := w.CreatePart(textproto.MIMEHeader{"Content-Type": {"application/json"}})
	if err != nil {
		return err
	}
	jsonPart.Write(reportDataJSON)

	// Добавление файла с указанием Content-Disposition
	filePartHeader := textproto.MIMEHeader{
		"Content-Disposition": {`form-data; name="file"; filename="report.pptx"`},
		"Content-Type":        {"application/vnd.openxmlformats-officedocument.presentationml.presentation"},
	}
	filePart, err := w.CreatePart(filePartHeader)
	if err != nil {
		return err
	}
	filepath := "/usr/src/app/internal/pkg/service/pptx_generator/report.pptx"
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	io.Copy(filePart, file)

	return nil
}

func generatePptx(reportDataJSON []byte) error {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctxWithTimeout, "python3", "pptx_generator.py")
	cmd.Stdin = bytes.NewReader(reportDataJSON)
	cmd.Dir = "/usr/src/app/internal/pkg/service/pptx_generator"

	return cmd.Run()
}
