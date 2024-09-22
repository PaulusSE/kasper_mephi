package student

import (
	"context"

	"uir_draft/internal/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (s *Service) GetPresentation(ctx context.Context, studentID uuid.UUID) (models.ReportData, error) {
	reportData := models.ReportData{}

	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		student, err := s.studRepo.GetStudentStatusTx(ctx, tx, studentID)
		if err != nil {
			return err
		}

		// semesterProgress, err := s.dissertationRepo.GetSemesterProgressTx(ctx, tx, studentID)
		// if err != nil {
		// 	return err
		// }

		// disTitles, err := s.dissertationRepo.GetDissertationTitlesTx(ctx, tx, studentID)
		// if err != nil {
		// 	return err
		// }

		// dissertationsStatuses, err := s.dissertationRepo.GetDissertationsTx(ctx, tx, studentID)
		// if err != nil {
		// 	return err
		// }

		// feedback, err := s.dissertationRepo.GetFeedbackTx(ctx, tx, studentID)
		// if err != nil {
		// 	return err
		// }

		// supervisors, err := s.studRepo.GetAllStudentsSupervisors(ctx, tx, studentID)
		// if err != nil {
		// 	return err
		// }

		// comments, err := s.commentRepo.GetStudentsCommentaries(ctx, tx, studentID)
		// if err != nil {
		// 	return err
		// }

		// progresses, err := s.dissertationRepo.GetStudentsProgressiveness(ctx, tx, studentID)
		// if err != nil {
		// 	return err
		// }

		// load, err := s.studRepo.GetLoad(ctx, student.StudentID, student.ActualSemester)
		// if err != nil {
		// 	return err
		// }

		reportData = models.ReportData{
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
			// PedagogicalData:      load,                   // Заполнить данными
			// ReportOtherAchievements: "Other Achievements", // Здесь нужно заменить на реальные данные
			// PedagogicalDataAll:    []PedagogicalWorkSummary{}, // Заполнить данными
			NextSemesterPlan: []string{"Plan 1", "Plan 2"}, // Заполнить данными
		}

		return nil
	})
	if err != nil {
		return models.ReportData{}, errors.Wrap(err, "GetPresentaion()")
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

	// Создаем структуру для передачи в Python скрипт

	// Преобразуем объект в JSON

	return reportData, nil
}

// Метод для получения данных по педагогической нагрузке студента
func (s *Service) GetStudentLoad(ctx context.Context, studentID uuid.UUID, actSem int32) ([]models.PedagogicalWork, error) {
	var pedagogicalWorks []models.PedagogicalWork

	// Получаем всю педагогическую нагрузку студента
	teachingLoads, err := s.GetTeachingLoad(ctx, studentID)
	if err != nil {
		return nil, errors.Wrap(err, "GetTeachingLoad()")
	}

	// Проходим по каждому элементу педагогической нагрузки
	for _, load := range teachingLoads {
		if load.Semester != int(actSem) {
			continue // Пропускаем нагрузки, не соответствующие заданному семестру
		}

		// Обрабатываем аудиторную нагрузку
		for _, classroomLoad := range load.ClassroomLoads {
			pedagogicalWorks = append(pedagogicalWorks, models.PedagogicalWork{
				Semester:    load.Semester,
				WorkType:    *classroomLoad.LoadType,
				Hours:       int(*classroomLoad.Hours),
				MainTeacher: *classroomLoad.MainTeacher,
				GroupName:   *classroomLoad.GroupName,
			})
		}

		// Обрабатываем индивидуальную работу со студентами
		for _, individualLoad := range load.IndividualStudentsLoads {
			pedagogicalWorks = append(pedagogicalWorks, models.PedagogicalWork{
				Semester:    load.Semester,
				WorkType:    *individualLoad.LoadType,
				Hours:       0,
				MainTeacher: "", // Укажите, если необходимо
				GroupName:   "", // Укажите, если необходимо
			})
		}

		// Обрабатываем дополнительную нагрузку
		for _, additionalLoad := range load.AdditionalLoads {
			pedagogicalWorks = append(pedagogicalWorks, models.PedagogicalWork{
				Semester:    load.Semester,
				WorkType:    *additionalLoad.Name,
				Hours:       0,  // Предполагается, что Volume это часы
				MainTeacher: "", // Укажите, если необходимо
				GroupName:   "", // Укажите, если необходимо
			})
		}
	}

	return pedagogicalWorks, nil
}
