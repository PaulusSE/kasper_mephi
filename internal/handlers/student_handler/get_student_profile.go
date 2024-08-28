package student_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"

	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
)

// GetStudentProfile
//
//	@Summary		Получение профиля студента
//
//	@Description	Получение профиля студента
//
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	models.StudentProfile	"Данные"
//	@Param			token	path		string					true	"Токен пользователя"
//	@Failure		400		{string}	string					"Неверный формат данных"
//	@Failure		401		{string}	string					"Токен протух"
//	@Failure		204		{string}	string					"Нет записей в БД"
//	@Failure		500		{string}	string					"Ошибка на стороне сервера"
//	@Router			/student/profile/{token} [get]
func (h *StudentHandler) GetStudentProfile(ctx *gin.Context) {
	user, err := h.authenticate(ctx)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	student, err := h.student.GetStudentsProfile(ctx, user.KasperID)
	if err != nil {
		ctx.AbortWithError(models.MapErrorToCode(err), err)
		return
	}

	// Получаем дополнительные данные

	// publications, err := h.student.GetPublications(ctx, student.StudentID)
	// if err != nil {
	// 	ctx.AbortWithError(http.StatusInternalServerError, err)
	// 	return
	// }

	// exams, err := h.student.GetExams(ctx, student.StudentID)
	// if err != nil {
	// 	ctx.AbortWithError(http.StatusInternalServerError, err)
	// 	return
	// }

	load, err := h.student.GetLoad(ctx, student.StudentID, student.ActualSemester)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Создаем структуру для передачи в Python скрипт
	reportData := models.ReportData{
		CurrentSemester:      student.ActualSemester,
		FullName:             student.FullName,
		SupervisorName:       "Supervisor Name",     // Здесь нужно заменить на реальные данные
		EducationDirection:   "Education Direction", // Здесь нужно заменить на реальные данные
		EducationProfile:     student.Specialization,
		EnrollmentDate:       student.StartDate,
		Specialty:            student.Specialization,
		TrainingYearFGOS:     "Training Year FGOS", // Здесь нужно заменить на реальные данные
		CandidateExams:       []models.Exam{},
		Category:             student.Category,
		Topic:                "Dissertation Topic", // Здесь нужно заменить на реальные данные
		ReportPeriodWork:     "Report Period Work", // Здесь нужно заменить на реальные данные
		ScientificObject:     "Scientific Object",  // Здесь нужно заменить на реальные данные
		ScientificSubject:    "Scientific Subject", // Здесь нужно заменить на реальные данные
		MentorRate:           "Mentor Rate",        // Здесь нужно заменить на реальные данные
		ProgressPercents:     []int{},              // Заполнить данными
		ProgressDescriptions: []string{},           // Заполнить данными
		Publications:         []models.Publication{},
		AllPublications:      []models.Publication{}, // Заполнить данными
		PedagogicalData:      load,                   // Заполнить данными
		// ReportOtherAchievements: "Other Achievements", // Здесь нужно заменить на реальные данные
		// PedagogicalDataAll:    []PedagogicalWorkSummary{}, // Заполнить данными
		NextSemesterPlan: []string{"Plan 1", "Plan 2"}, // Заполнить данными
	}

	// Преобразуем объект в JSON
	reportDataJSON, err := json.Marshal(reportData)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Вызов Python скрипта для генерации презентации
	cmd := exec.Command("python3", "generate_presentation.py")
	cmd.Stdin = bytes.NewReader(reportDataJSON)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Отправляем файл пользователю
	ctx.Header("Content-Disposition", "attachment; filename=report.pptx")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.presentationml.presentation", out.Bytes())

	ctx.JSON(http.StatusOK, student)
}
