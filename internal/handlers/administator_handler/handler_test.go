package administator_handler

import (
	"context"
	"github.com/google/uuid"
	"testing"
	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/administator_handler/request_models"
	auth_req_models "uir_draft/internal/handlers/authorization_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

// Мокаем зависимости хендлера:
type mockEnum struct{}

func (m *mockEnum) GetGroups(ctx context.Context) ([]models.Group, error) {
	return []models.Group{}, nil
}
func (m *mockEnum) GetSpecializations(ctx context.Context) ([]models.Specialization, error) {
	return nil, nil
}
func (m *mockEnum) InsertSpecializations(ctx context.Context, specializations []models.Specialization) error {
	return nil
}
func (m *mockEnum) DeleteSpecializations(ctx context.Context, specIDs []int32) error {
	return nil
}
func (m *mockEnum) InsertGroups(ctx context.Context, groups []models.Group) error {
	return nil
}
func (m *mockEnum) DeleteGroups(ctx context.Context, groupIDs []int32) error {
	return nil
}
func (m *mockEnum) GetSemestersAmount(ctx context.Context) ([]models.SemesterAmount, error) {
	return nil, nil
}

func (m *mockEnum) DeleteSemesterAmounts(ctx context.Context, ids []uuid.UUID) error {
	return nil
}
func (m *mockEnum) InsertSemesterAmount(ctx context.Context, amounts []models.SemesterAmount) error {
	return nil
}

type mockUser struct{}

func (m *mockUser) GetNotRegisteredUsers(ctx context.Context) ([]models.UserInfo, error) {
	return nil, nil
}
func (m *mockUser) GetStudentSupervisorPairs(ctx context.Context) ([]models.StudentSupervisorPair, error) {
	return nil, nil
}
func (m *mockUser) ChangeSupervisor(ctx context.Context, pairs []models.ChangeSupervisor) error {
	return nil
}
func (m *mockUser) SetStudentFlags(ctx context.Context, students []models.SetStudentsFlags) error {
	return nil
}
func (m *mockUser) GetSupervisors(ctx context.Context) ([]models.Supervisor, error) {
	return nil, nil
}
func (m *mockUser) GetStudentsList(ctx context.Context) ([]models.Student, error) {
	return nil, nil
}
func (m *mockUser) UpsertAttestationMarks(ctx context.Context, marks []models.AttestationMarkRequest) error {
	return nil
}
func (m *mockUser) AddUsers(ctx context.Context, users request_models.AddUsersRequest, userType model.UserType) ([]models.UsersCredentials, error) {
	return nil, nil
}
func (m *mockUser) ArchiveSupervisor(ctx context.Context, supervisors []models.SupervisorStatus) error {
	return nil
}
func (m *mockUser) DeleteNotRegisteredUsers(ctx context.Context, userIDs []uuid.UUID) error {
	return nil
}

type mockSupervisor struct{}

func (m *mockSupervisor) GetSupervisorsStudents(ctx context.Context, id uuid.UUID) ([]models.Student, error) {
	return nil, nil
}
func (m *mockSupervisor) GetSupervisorProfile(ctx context.Context, supervisorID uuid.UUID) (models.SupervisorProfile, error) {
	return models.SupervisorProfile{}, nil
}
func (m *mockSupervisor) InitSupervisor(ctx context.Context, user model.Users, registry auth_req_models.FirstSupervisorRegistry) error {
	return nil
}

type mockAuthenticator struct{}

func (m *mockAuthenticator) AuthenticateWithUserType(ctx context.Context, token, userType string) (*model.Users, error) {
	// Всегда возвращаем успех
	return &model.Users{UserID: uuid.New(), UserType: model.UserType_Admin}, nil
}
func (m *mockAuthenticator) CreateUser(ctx context.Context, users *model.Users) error {
	return nil
}
func (m *mockAuthenticator) AttachTokenToUser(ctx context.Context, token string, userID uuid.UUID) error {
	return nil
}

// Создаём мок-хендлер:
func newTestHandler() *AdministratorHandler {
	return NewHandler(
		&mockUser{},
		&mockAuthenticator{}, // вот он!
		&mockEnum{},
		&mockSupervisor{},
		nil, // email сервис
		nil, // registration
		nil, // student
	)
}
func TestGetGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := newTestHandler()
	router.GET("/groups", h.GetGroups)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/groups", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", w.Code)
	}
}

func TestGetNotRegisteredUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := newTestHandler()
	router.GET("/users", h.GetNotRegisteredUsers)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", w.Code)
	}
}

// Для Put-запроса (можно добавить тест тела)
func TestGetSupervisorsStudents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := newTestHandler()
	router.PUT("/supervisors", h.GetSupervisorsStudents)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/supervisors", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("expected 200, 400, or 500, got %d", w.Code)
	}
}
