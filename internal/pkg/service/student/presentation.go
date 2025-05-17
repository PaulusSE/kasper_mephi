package student

import (
	"context"
	"fmt"

	"uir_draft/internal/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (s *Service) GetPresentation(ctx context.Context, studentID uuid.UUID, semester int) (models.ReportData, error) {
	reportData := models.ReportData{}
	universityFullName := "МИНИСТЕРСТВО НАУКИ И ВЫСШЕГО ОБРАЗОВАНИЯ РОССИЙСКОЙ ФЕДЕРАЦИИ\nФЕДЕРАЛЬНОЕ ... НИЯУ МИФИ"
	instituteAndDepartment := "ИНСТИТУТ ... КАФЕДРА 22 (Кибернетика)"
	educationDirection := "Информатика и вычислительная техника"
	educationDirectionCode := "09.06.01"
	educationProfile := "Системный анализ, управление и обработка информации, статистика"
	educationProfileCode := "2.3.5"
	reportTitle := fmt.Sprintf("Отчет аспиранта за %d семестр", semester)
	city := "Москва"
	year := "2024" // можно генерировать
	logoMephi := "MEPhI_Logo2014_en.png"
	logoKafedra := "kaf22.png"

	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		fmt.Printf("STUDEN_ID - %s \n", studentID)
		student, err := s.studRepo.GetStudentStatusTx(ctx, tx, studentID)
		if err != nil {
			return errors.Wrap(err, "GetStudentStatusTx()")
		}

		// Получаем данные о текущем научном руководителе
		supervisor, err := s.studRepo.GetStudentsActualSupervisorTx(ctx, tx, studentID)
		if err != nil {
			// Можно возвращать "Не назначен" или пусто, если не найдено
			supervisor.FullName = "Не назначен"
		}

		semesterProgress, err := s.dissertationRepo.GetStudentsProgressiveness(ctx, tx, studentID)
		if err != nil {
			return err
		}
		var progressPerCents []int
		for _, progItem := range semesterProgress {
			progressPerCents = append(progressPerCents, int(progItem.Progressiveness))
		}

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

		//load, err := s.GetStudentLoad(ctx, student.StudentID, int32(semester))
		//if err != nil {
		//	return err
		//}

		reportData = models.ReportData{
			UniversityFullName:     universityFullName,
			InstituteAndDepartment: instituteAndDepartment,
			EducationDirection:     educationDirection,
			EducationDirectionCode: educationDirectionCode,
			EducationProfile:       educationProfile,
			EducationProfileCode:   educationProfileCode,
			ReportTitle:            reportTitle,
			City:                   city,
			Year:                   year,
			FullName:               student.FullName,
			SupervisorName:         supervisor.FullName,
			CurrentSemester:        semester,
			EnrollmentDate:         student.StartDate.Format("2006-01-02"),
			Specialty:              student.Specialization,
			//TrainingYearFGOS:       student.TrainingYearFGOS,
			//CandidateExams:         candidateExams,
			Category: student.Category,
			//Topic:                  student.Topic,
			//ReportPeriodWork:       student.ReportPeriodWork,
			//ScientificObject:       student.ScientificObject,
			//ScientificSubject:      student.ScientificSubject,
			//MentorRate:             student.MentorRate,
			ProgressPercents: progressPerCents,
			//ProgressDescriptions:   progressDescriptions,
			//Publications:           publications,
			//AllPublications:        allPublications,
			//PedagogicalData:        pedagogicalData,
			//ReportOtherAchievments: student.ReportOtherAchievments,
			//PedagogicalDataAll:     load,
			//NextSemesterPlan:       student.NextSemesterPlan,
			LogoMephi:   logoMephi,
			LogoKafedra: logoKafedra,
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
