package supervisor_handler

import (
	"context"
	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	auth_req_models "uir_draft/internal/handlers/authorization_handler/request_models"
	"uir_draft/internal/pkg/helpers"
	"uir_draft/internal/pkg/models"
	"uir_draft/internal/pkg/repository"
	"uir_draft/internal/pkg/service/student"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	DissertationService interface {
		AllToStatus(ctx context.Context, studentID uuid.UUID, comment *string, status string) error
		// GetDissertationPage - возвращает всю информацию для отрисовки страницы диссертации
		GetDissertationPage(ctx context.Context, studentID uuid.UUID) (models.DissertationPageResponse, error)
		// GetDissertationData - возвращает данные для скачивания файла
		GetDissertationData(ctx context.Context, studentID uuid.UUID, semester int32) (model.Dissertations, error)
	}

	ScientificWorksService interface {
		// GetScientificWorks - возвращает все научные работы студента
		GetScientificWorks(ctx context.Context, studentID uuid.UUID) ([]models.ScientificWork, error)
	}

	TeachingLoadService interface {
		// GetTeachingLoad - возвращает всю педагогическую нагрузку студента
		GetTeachingLoad(ctx context.Context, studentID uuid.UUID) ([]models.TeachingLoad, error)
	}

	ReportService interface {
		GetReportComments(ctx context.Context, studentID uuid.UUID) (models.ReportComments, error)
	}

	RegistrationService interface {
		ReviewRequest(ctx context.Context, requestID uuid.UUID, approve bool) error
		ListRequests(ctx context.Context, userType string) ([]repository.RegistrationRequest, error)
		GetRequestByID(ctx context.Context, id uuid.UUID) (repository.RegistrationRequest, error)
	}

	Authenticator interface {
		// Authenticate - проводит аутентификацию пользователя
		AuthenticateWithUserType(ctx context.Context, token, userType string) (*model.Users, error)
		CreateUser(ctx context.Context, user *model.Users) error
	}

	SupervisorService interface {
		// UpsertFeedback - обновляет или добавляет фидбэк от научного руководителя
		UpsertFeedback(ctx context.Context, studentID, supervisorID uuid.UUID, request models.FeedbackRequest) error
		GetSupervisorsStudents(ctx context.Context, supervisorID uuid.UUID) ([]models.Student, error)
		GetSupervisorProfile(ctx context.Context, supervisorID uuid.UUID) (models.SupervisorProfile, error)
		UpsertSupervisorMark(ctx context.Context, studentID, supervisorID uuid.UUID, semester, mark int32) error
		InitSupervisor(ctx context.Context, user model.Users, registry auth_req_models.FirstSupervisorRegistry) error
	}

	EmailService interface {
		SendMailToStudent(ctx context.Context, studentID, supervisorID uuid.UUID, templatePath, tt, status string) error
	}

	StudentService interface {
		GetStudentStatus(ctx context.Context, studentID uuid.UUID) (models.Student, error)
		GetStudentsProfile(ctx context.Context, studentID uuid.UUID) (models.StudentProfile, error)
		GetAllMarks(ctx context.Context, studentID uuid.UUID) (models.AllMarks, error)
		InitStudent(ctx context.Context, user model.Users, req auth_req_models.FirstStudentRegistry) error
	}
)

type SupervisorHandler struct {
	dissertation  DissertationService
	scientific    ScientificWorksService
	load          TeachingLoadService
	authenticator Authenticator
	supervisor    SupervisorService
	student       StudentService
	email         EmailService
	report        ReportService
	registration  RegistrationService
}

func NewHandler(
	dissertation *student.Service,
	authenticator Authenticator,
	supervisor SupervisorService,
	email EmailService,
	registration RegistrationService,
) *SupervisorHandler {
	return &SupervisorHandler{
		dissertation:  dissertation,
		scientific:    dissertation,
		load:          dissertation,
		authenticator: authenticator,
		supervisor:    supervisor,
		student:       dissertation,
		email:         email,
		report:        dissertation,
		registration:  registration,
	}
}

func (h *SupervisorHandler) authenticate(ctx *gin.Context) (*model.Users, error) {
	token := helpers.GetToken(ctx)

	user, err := h.authenticator.AuthenticateWithUserType(ctx, token, model.UserType_Supervisor.String())
	if err != nil {
		return user, err
	}

	return user, nil
}
