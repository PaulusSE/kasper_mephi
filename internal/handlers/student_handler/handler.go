package student_handler

import (
	"context"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/student_handler/request_models"
	"uir_draft/internal/pkg/helpers"
	"uir_draft/internal/pkg/models"
	"uir_draft/internal/pkg/service/student"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	StudentService interface {
		// SetStudentStatus - проставляет аспиранту статус подтверждения
		SetStudentStatus(ctx context.Context, studentID uuid.UUID, status model.ApprovalStatus) error
		// GetStudentStatus - возвращает статус студента
		GetStudentStatus(ctx context.Context, studentID uuid.UUID) (models.Student, error)
		UpdateStudentsProgressiveness(ctx context.Context, studentID uuid.UUID, progress int32) error
		GetStudentsProfile(ctx context.Context, studentID uuid.UUID) (models.StudentProfile, error)

		UpdateStudentsProfile(ctx context.Context, userID, studentID uuid.UUID, studentInfo models.UpdateProfile) error

		GetLoad(ctx context.Context, studentID uuid.UUID, actSem int32) ([]models.PedagogicalWork, error)
	}

	DissertationService interface {
		// DissertationToStatus - проставляет статус для диссертации
		DissertationToStatus(ctx context.Context, studentID uuid.UUID, status model.ApprovalStatus, semester int32) error
		// AllToStatus - проставляет статус для всего
		AllToStatus(ctx context.Context, studentID uuid.UUID, comment *string, status string) error
		// GetDissertationPage - возвращает всю информацию для отрисовки страницы диссертации
		GetDissertationPage(ctx context.Context, studentID uuid.UUID) (models.DissertationPageResponse, error)
		// UpsertSemesterProgress - обновляет план подготовки диссертации
		UpsertSemesterProgress(ctx context.Context, studentID uuid.UUID, progress []models.SemesterProgressRequest) error
		// UpsertDissertationInfo - обновляет информацию о диссертации (файле)
		UpsertDissertationInfo(ctx context.Context, studentID uuid.UUID, semester int32, fileName string) error
		UpsertDissertationTitle(ctx context.Context, studentID uuid.UUID, title, object, order string) error
		// GetDissertationData - возвращает информацию о диссертации (файле)
		GetDissertationData(ctx context.Context, studentID uuid.UUID, semester int32) (model.Dissertations, error)
	}

	ScientificWorksService interface {
		// ScientificWorksToStatus - ставит статус для научных работ
		ScientificWorksToStatus(ctx context.Context, studentID uuid.UUID, status model.ApprovalStatus, semester int32) error
		// GetScientificWorks - возвращает все научные работы студента
		GetScientificWorks(ctx context.Context, studentID uuid.UUID) ([]models.ScientificWork, error)
		// UpsertPublications - добавляет или обновляет научные публикации
		UpsertPublications(ctx context.Context, studentID uuid.UUID, semester int32, publications []models.Publication) error
		// UpsertConferences - добавляет или обновляет научные конференции
		UpsertConferences(ctx context.Context, studentID uuid.UUID, semester int32, conferences []models.Conference) error
		// UpsertResearchProjects - добавляет или обновляет научные исследования
		UpsertResearchProjects(ctx context.Context, studentID uuid.UUID, semester int32, projects []models.ResearchProject) error
		// UpsertPatents - добавляет или обновляет патенты
		UpsertPatents(ctx context.Context, studentID uuid.UUID, semester int32, patents []models.Patent) error
		// DeletePublications - удаляет научные публикации
		DeletePublications(ctx context.Context, studentID uuid.UUID, semester int32, ids []uuid.UUID) error
		// DeleteConferences - удаляет научные конференции
		DeleteConferences(ctx context.Context, studentID uuid.UUID, semester int32, ids []uuid.UUID) error
		// DeleteResearchProjects - удаляет научные исследования
		DeleteResearchProjects(ctx context.Context, studentID uuid.UUID, semester int32, ids []uuid.UUID) error
		// DeletePatents - удаляет патенты
		DeletePatents(ctx context.Context, studentID uuid.UUID, semester int32, ids []uuid.UUID) error
	}

	TeachingLoadService interface {
		// TeachingLoadToStatus - ставит статус для пед нагрузки
		TeachingLoadToStatus(ctx context.Context, studentID uuid.UUID, status model.ApprovalStatus, semester int32) error
		// GetTeachingLoad - возвращает всю педагогическую нагрузку студента
		GetTeachingLoad(ctx context.Context, studentID uuid.UUID) ([]models.TeachingLoad, error)
		// UpsertClassroomLoad - добавляет или обновляет аудиторную педагогическую нагрузку
		UpsertClassroomLoad(ctx context.Context, studentID uuid.UUID, semester int32, loads []models.ClassroomLoad) error
		// UpsertIndividualLoad - добавляет или обновляет индивидуальную педагогическую нагрузку
		UpsertIndividualLoad(ctx context.Context, studentID uuid.UUID, semester int32, loads []models.IndividualStudentsLoad) error
		// UpsertAdditionalLoad - добавляет или обновляет дополнительную педагогическую нагрузку
		UpsertAdditionalLoad(ctx context.Context, studentID uuid.UUID, semester int32, loads []models.AdditionalLoad) error
		// DeleteClassroomLoad - удаляет аудиторную педагогическую нагрузку студента
		DeleteClassroomLoad(ctx context.Context, studentID uuid.UUID, semester int32, loads []uuid.UUID) error
		// DeleteIndividualLoad - удаляет индивидуальную педагогическую нагрузку студента
		DeleteIndividualLoad(ctx context.Context, studentID uuid.UUID, semester int32, loads []uuid.UUID) error
		// DeleteAdditionalLoad - удаляет дополнительную аудиторную педагогическую нагрузку студента
		DeleteAdditionalLoad(ctx context.Context, studentID uuid.UUID, semester int32, loads []uuid.UUID) error
	}

	ReportService interface {
		GetReportComments(ctx context.Context, studentID uuid.UUID) (models.ReportComments, error)
		UpsertReportComments(ctx context.Context, studentID uuid.UUID, req request_models.UpsertReportCommentsRequest) error
	}

	Authenticator interface {
		// Authenticate - проводит аутентификацию пользователя
		AuthenticateWithUserType(ctx context.Context, token, userType string) (*model.Users, error)
	}

	EmailService interface {
		SendMailToSupervisor(ctx context.Context, studentID uuid.UUID, templatePath, tt string) error
	}

	EnumService interface {
		GetSpecializations(ctx context.Context) ([]models.Specialization, error)
		GetGroups(ctx context.Context) ([]models.Group, error)
		GetSemestersAmount(ctx context.Context) ([]models.SemesterAmount, error)
	}

	AdminService interface {
		GetSupervisors(ctx context.Context) ([]models.Supervisor, error)
	}

	MarksService interface {
		GetAllMarks(ctx context.Context, studentID uuid.UUID) (models.AllMarks, error)
		UpsertExamResults(ctx context.Context, studentID uuid.UUID, exams []models.ExamRequest) error
		DeleteExamMarks(ctx context.Context, semester int32, ids []uuid.UUID) error
	}

	PresentationService interface {
		GetPresentation(ctx context.Context, studentID uuid.UUID) (models.ReportData, error)
	}
)

type StudentHandler struct {
	student      StudentService
	dissertation DissertationService
	scientific   ScientificWorksService
	load         TeachingLoadService
	mark         MarksService

	authenticator Authenticator
	email         EmailService
	enum          EnumService
	admin         AdminService
	report        ReportService
}

func NewHandler(
	student *student.Service,
	authenticator Authenticator,
	email EmailService,
	enum EnumService,
	admin AdminService,
) *StudentHandler {
	return &StudentHandler{
		student:       student,
		dissertation:  student,
		scientific:    student,
		load:          student,
		authenticator: authenticator,
		mark:          student,
		email:         email,
		enum:          enum,
		admin:         admin,
		report:        student,
	}
}

func (h *StudentHandler) authenticate(ctx *gin.Context) (*model.Users, error) {
	token := helpers.GetToken(ctx)

	user, err := h.authenticator.AuthenticateWithUserType(ctx, token, model.UserType_Student.String())
	if err != nil {
		return user, err
	}

	return user, nil
}
