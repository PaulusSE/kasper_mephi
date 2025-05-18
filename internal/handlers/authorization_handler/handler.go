package authorization_handler

import (
	"context"
	"fmt"
	"net/http"
	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/authorization_handler/request_models"
	"uir_draft/internal/pkg/helpers"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	Authenticator interface {
		CreateAnonymousToken(ctx context.Context, token string) error
		TokenExists(ctx context.Context, token string) (bool, error)
		Authorize(ctx context.Context, request models.AuthorizeRequest) (*models.AuthorizeResponse, bool, error)
		AuthenticateWithUserType(ctx context.Context, token, userType string) (*model.Users, error)
		Authenticate(ctx context.Context, token string) (*model.Users, error)
		TokenCheck(ctx context.Context, token string) (*model.Users, error)
		ChangePassword(ctx context.Context, userID uuid.UUID, request request_models.ChangePasswordRequest) error
		GetUserProfile(ctx context.Context, userID uuid.UUID) (model.Users, error)
		CreateUser(ctx context.Context, user *model.Users) error
		AttachTokenToUser(ctx context.Context, token string, userID uuid.UUID) error
	}

	StudentService interface {
		InitStudent(ctx context.Context, user model.Users, req request_models.FirstStudentRegistry) error
	}

	SupervisorService interface {
		InitSupervisor(ctx context.Context, user model.Users, registry request_models.FirstSupervisorRegistry) error
	}

	RegistrationService interface {
		CreateStudentRequest(ctx context.Context, hashedPassword []byte, data request_models.FirstStudentRegistry) error
		CreateSupervisorRequest(ctx context.Context, hashedPassword []byte, data request_models.FirstSupervisorRegistry) error
	}
)

type AuthorizationHandler struct {
	authenticator Authenticator
	student       StudentService
	supervisor    SupervisorService
	registration  RegistrationService
}

func NewHandler(
	authenticator Authenticator,
	student StudentService,
	supervisor SupervisorService,
	registration RegistrationService,
) *AuthorizationHandler {
	return &AuthorizationHandler{
		authenticator: authenticator,
		student:       student,
		supervisor:    supervisor,
		registration:  registration,
	}
}

func (h *AuthorizationHandler) authenticateStudent(ctx *gin.Context) (*model.Users, error) { // nolint
	token := helpers.GetToken(ctx)

	user, err := h.authenticator.AuthenticateWithUserType(ctx, token, model.UserType_Student.String())
	if err != nil {
		return user, err
	}

	return user, nil
}

func (h *AuthorizationHandler) authenticateSupervisor(ctx *gin.Context) (*model.Users, error) { // nolint
	token := helpers.GetToken(ctx)

	user, err := h.authenticator.AuthenticateWithUserType(ctx, token, model.UserType_Supervisor.String())
	if err != nil {
		return user, err
	}

	return user, nil
}

func (h *AuthorizationHandler) authenticate(ctx *gin.Context) (*model.Users, error) {
	token := helpers.GetToken(ctx)

	user, err := h.authenticator.Authenticate(ctx, token)
	if err != nil {
		return user, err
	}

	return user, nil
}

// ValidateToken проверяет, что переданный в запросе токен есть в таблице authorization_token
// Если токен невалиден — прервёт контекст с соответствующим статусом
func (h *AuthorizationHandler) ValidateToken(ctx *gin.Context) error {
	token := helpers.GetToken(ctx)
	ok, err := h.authenticator.TokenExists(ctx.Request.Context(), token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return err
	}
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return fmt.Errorf("token not found")
	}
	return nil
}
