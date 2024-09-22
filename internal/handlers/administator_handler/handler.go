package administator_handler

import (
	"context"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/administator_handler/request_models"
	"uir_draft/internal/pkg/helpers"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetStudentSupervisorPairs(ctx context.Context) ([]models.StudentSupervisorPair, error)

type (
	UserService interface {
		GetStudentSupervisorPairs(ctx context.Context) ([]models.StudentSupervisorPair, error)
		ChangeSupervisor(ctx context.Context, pairs []models.ChangeSupervisor) error
		SetStudentFlags(ctx context.Context, students []models.SetStudentsFlags) error
		GetSupervisors(ctx context.Context) ([]models.Supervisor, error)
		GetStudentsList(ctx context.Context) ([]models.Student, error)
		UpsertAttestationMarks(ctx context.Context, marks []models.AttestationMarkRequest) error
		AddUsers(ctx context.Context, users request_models.AddUsersRequest, userType model.UserType) ([]models.UsersCredentials, error)
		ArchiveSupervisor(ctx context.Context, supervisors []models.SupervisorStatus) error
		GetNotRegisteredUsers(ctx context.Context) ([]models.UserInfo, error)
		DeleteNotRegisteredUsers(ctx context.Context, userIDs []uuid.UUID) error
	}

	Authenticator interface {
		// AuthenticateWithUserType - проводит аутентификацию пользователя
		AuthenticateWithUserType(ctx context.Context, token, userType string) (*model.Users, error)
	}

	EnumService interface {
		GetSpecializations(ctx context.Context) ([]models.Specialization, error)
		InsertSpecializations(ctx context.Context, specializations []models.Specialization) error
		DeleteSpecializations(ctx context.Context, specIDs []int32) error

		GetGroups(ctx context.Context) ([]models.Group, error)
		InsertGroups(ctx context.Context, groups []models.Group) error
		DeleteGroups(ctx context.Context, groupIDs []int32) error

		GetSemestersAmount(ctx context.Context) ([]models.SemesterAmount, error)
		DeleteSemesterAmounts(ctx context.Context, ids []uuid.UUID) error

		InsertSemesterAmount(ctx context.Context, amounts []models.SemesterAmount) error
	}

	SupervisorService interface {
		GetSupervisorsStudents(ctx context.Context, supervisorID uuid.UUID) ([]models.Student, error)
		GetSupervisorProfile(ctx context.Context, supervisorID uuid.UUID) (models.SupervisorProfile, error)
	}

	EmailService interface {
		SendInviteEmails(ctx context.Context, credentials []models.UsersCredentials, templatePath string) error
	}

	PresentationService interface {
		GetPresentation(ctx context.Context, studentID uuid.UUID) (models.ReportData, error)
		GetStudentLoad(ctx context.Context, studentID uuid.UUID, actSem int32) ([]models.PedagogicalWork, error)
	}
)

type AdministratorHandler struct {
	user          UserService
	authenticator Authenticator
	enum          EnumService
	supervisor    SupervisorService
	email         EmailService
	presentation  PresentationService
}

func NewHandler(
	user UserService,
	authenticator Authenticator,
	enum EnumService,
	supervisor SupervisorService,
	email EmailService,
	presentation PresentationService,
) *AdministratorHandler {
	return &AdministratorHandler{
		user:          user,
		authenticator: authenticator,
		enum:          enum,
		supervisor:    supervisor,
		email:         email,
		presentation:  presentation,
	}
}

func (h *AdministratorHandler) authenticate(ctx *gin.Context) (*model.Users, error) {
	token := helpers.GetToken(ctx)

	user, err := h.authenticator.AuthenticateWithUserType(ctx, token, model.UserType_Admin.String())
	if err != nil {
		return user, err
	}

	return user, nil
}
