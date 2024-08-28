package student

import (
	"context"
	"time"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	auth_req "uir_draft/internal/handlers/authorization_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func (s *Service) AllToStatus(ctx context.Context, studentID uuid.UUID, comment *string, status string) error {
	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		student, err := s.studRepo.GetStudentTx(ctx, tx, studentID)
		if err != nil {
			return err
		}

		dStatus, err := models.MapApprovalStatusToDomain(status)
		if err != nil {
			return errors.Wrap(err, "AllToStatus()")
		}

		err = s.dissertationRepo.SetDissertationStatusTx(ctx, tx, student.StudentID, dStatus, student.ActualSemester)
		if err != nil {
			return err
		}

		//err = s.dissertationRepo.SetDissertationTitleStatusTx(ctx, tx, student.StudentID, dStatus, student.ActualSemester, nil)
		//if err != nil {
		//	return err
		//}

		if comment != nil {
			err = s.commentRepo.UpsertStudentsComment(ctx, tx, model.StudentsCommentary{
				CommentaryID: uuid.New(),
				StudentID:    studentID,
				Semester:     student.ActualSemester,
				Commentary:   comment,
				CommentedAt:  time.Now(),
			})

			if err != nil {
				return err
			}
		}

		err = s.dissertationRepo.SetSemesterProgressStatusTx(ctx, tx, student.StudentID, dStatus, nil)
		if err != nil {
			return err
		}

		err = s.scienceRepo.SetScientificWorkStatusTx(ctx, tx, student.StudentID, dStatus, student.ActualSemester, nil)
		if err != nil {
			return err
		}

		err = s.loadRepo.SetTeachingLoadStatusTx(ctx, tx, student.StudentID, dStatus, student.ActualSemester, nil)
		if err != nil {
			return err
		}

		err = s.studRepo.SetStudentStatusTx(ctx, tx, dStatus, student.StudyingStatus, student.StudentID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "AllToStatus()")
	}

	return nil
}

func (s *Service) GetStudentStatus(ctx context.Context, studentID uuid.UUID) (models.Student, error) {
	var student models.Student

	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		var err error
		student, err = s.studRepo.GetStudentStatusTx(ctx, tx, studentID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.Student{}, errors.Wrap(err, "GetStudentStatus()")
	}

	return student, nil
}

func (s *Service) SetStudentStatus(ctx context.Context, studentID uuid.UUID, status model.ApprovalStatus) error {
	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		student, err := s.studRepo.GetStudentTx(ctx, tx, studentID)
		if err != nil {
			return err
		}

		if err = s.studRepo.SetStudentStatusTx(ctx, tx, status, student.StudyingStatus, studentID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "SetStudentStatus()")
	}

	return nil
}

func (s *Service) InitStudent(ctx context.Context, user model.Users, req auth_req.FirstStudentRegistry) error {
	startDate, err := time.Parse(time.DateOnly, req.StartDate)
	if err != nil {
		return errors.Wrap(err, "InitStudent()")
	}

	student := model.Students{
		StudentID: user.KasperID,
		UserID:    user.UserID,
		FullName:  req.FullName,
		//Department:     req.Department,
		SpecID:         req.SpecializationID,
		ActualSemester: req.ActualSemester,
		Years:          req.NumberOfYears,
		StartDate:      startDate,
		GroupID:        req.GroupID,
		Category:       req.Category,
		Phone:          req.Phone,
		EndDate:        startDate.AddDate(4, 0, 0),
	}

	var progresses []model.SemesterProgress
	progressTypes := []model.ProgressType{
		model.ProgressType_Intro,
		model.ProgressType_Ch1,
		model.ProgressType_Ch2,
		model.ProgressType_Ch3,
		model.ProgressType_Ch4,
		model.ProgressType_Ch5,
		model.ProgressType_Ch6,
		model.ProgressType_End,
		model.ProgressType_Literature,
		model.ProgressType_Abstract,
	}

	updatedAt := time.Now()

	for _, progressType := range progressTypes {
		progress := model.SemesterProgress{
			ProgressID:   uuid.New(),
			StudentID:    user.KasperID,
			ProgressType: progressType,
			First:        false,
			Second:       false,
			Third:        false,
			Forth:        false,
			Fifth:        false,
			Sixth:        false,
			Seventh:      false,
			Eighth:       false,
			UpdatedAt:    updatedAt,
			Status:       model.ApprovalStatus_Empty,
			AcceptedAt:   nil,
		}

		progresses = append(progresses, progress)
	}

	err = s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		if err = s.studRepo.InsertStudentTx(ctx, tx, student); err != nil {
			return err
		}

		if err = s.scienceRepo.InitScientificWorkStatusTx(ctx, tx, student.StudentID, student.Years); err != nil {
			return err
		}

		if err = s.loadRepo.InitTeachingLoadsStatusTx(ctx, tx, student.StudentID, student.Years); err != nil {
			return err
		}

		if err = s.dissertationRepo.UpsertSemesterProgressTx(ctx, tx, progresses); err != nil {
			return err
		}

		if err = s.studRepo.SetNewSupervisorTx(ctx, tx, user.KasperID, lo.FromPtr(req.SupervisorID)); err != nil {
			return err
		}

		if err = s.userRepo.SetUserRegisteredTx(ctx, tx, user.UserID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "InitStudent()")
	}

	return nil
}

func (s *Service) GetStudentsProfile(ctx context.Context, studentID uuid.UUID) (models.StudentProfile, error) {
	var student models.StudentProfile

	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		var err error
		student, err = s.studRepo.GetStudentProfile(ctx, tx, studentID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return models.StudentProfile{}, errors.Wrap(err, "GetStudentsProfile()")
	}

	return student, nil
}

// Метод для получения данных по публикациям студента
// func (s *Service) GetPublications(ctx context.Context, studentID uuid.UUID) ([]Publication, error) {
// 	var publications []Publication

// 	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
// 		ids, err := s.scienceRepo.GetScientificWorksStatusIDs(ctx, tx, studentID)
// 		if err != nil {
// 			return err
// 		}

// 		// Преобразуйте данные, полученные из базы данных, в нужную структуру
// 		dbPublications, err := s.scienceRepo.GetPublicationsTx(ctx, tx, ids)
// 		if err != nil {
// 			return err
// 		}

// 		for _, dbPub := range dbPublications {
// 			publications = append(publications, Publication{
// 				Name:       dbPub.Name,
// 				Status:     string(dbPub.Status),
// 				Impact:     dbPub.Impact,
// 				OutputData: dbPub.OutputData,
// 				CoAuthors:  dbPub.CoAuthors,
// 				Volume:     dbPub.Volume,
// 			})
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "GetPublications()")
// 	}

// 	return publications, nil
// }

// Метод для получения данных по экзаменам студента
// func (s *Service) GetExams(ctx context.Context, studentID uuid.UUID) ([]Exam, error) {
// 	var exams []Exam

// 	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
// 		dbExams, err := s.marksRepo.GetStudentsExamResults(ctx, tx, studentID)
// 		if err != nil {
// 			return err
// 		}

// 		for _, dbExam := range dbExams {
// 			exams = append(exams, Exam{
// 				ExamName: dbExam.ExamName,
// 				Mark:     dbExam.Mark,
// 			})
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "GetExams()")
// 	}

// 	return exams, nil
// }

// Метод для получения данных по педагогической нагрузке студента
func (s *Service) GetLoad(ctx context.Context, studentID uuid.UUID, actSem int32) ([]models.PedagogicalWork, error) {
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

// func (s *Service) GetPresentationData(ctx context.Context, studentID uuid.UUID) ([]model.Dissertations, error) {
// 	var dissertations []model.Dissertations

// 	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
// 		var err error
// 		dissertations, err = s.dissertationRepo.GetDissertationsTx(ctx, tx, studentID)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "GetDissertationData()")
// 	}

// 	return dissertations, nil
// }

func (s *Service) UpdateStudentsProfile(ctx context.Context, userID, studentID uuid.UUID, studentInfo models.UpdateProfile) error {
	//groupID, err := strconv.ParseInt(studentInfo.GroupID, 10, 32)
	//if err != nil {
	//	return errors.Wrap(err, "UpdateStudentsProfile()")
	//}

	student := model.Students{
		StudentID: studentID,
		FullName:  studentInfo.FullName,
		Years:     studentInfo.Years,
		StartDate: studentInfo.StartDate,
		GroupID:   studentInfo.GroupID,
		Phone:     studentInfo.Phone,
		Category:  studentInfo.Category,
	}

	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		if err := s.studRepo.UpdateStudent(ctx, tx, student); err != nil {
			return err
		}

		if err := s.userRepo.ChangeUsersEmail(ctx, tx, userID, studentInfo.Email); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "UpdateStudentsProfile()")
	}

	return nil
}
