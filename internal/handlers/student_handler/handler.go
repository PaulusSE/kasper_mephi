package student_handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	"time"

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
		// GetStudentsProgressiveness(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.Progressiveness, error)
		UpdateStudentsProfile(ctx context.Context, userID, studentID uuid.UUID, studentInfo models.UpdateProfile) error

		GetPresentation(ctx context.Context, studentID uuid.UUID) (models.ReportData, error)
		GetStudentLoad(ctx context.Context, studentID uuid.UUID, actSem int32) ([]models.PedagogicalWork, error)
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
		TokenExists(ctx context.Context, token string) (bool, error)
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

	RecommendationCacheService interface {
		GetCachedRecommendations(ctx context.Context, studentID uuid.UUID, semester int32, topN int) ([]Recommendation, *time.Time, error)
		SaveCachedRecommendations(ctx context.Context, studentID uuid.UUID, semester int32, topN int, recs []Recommendation) error
	}
)

type StudentHandler struct {
	student             StudentService
	dissertation        DissertationService
	scientific          ScientificWorksService
	load                TeachingLoadService
	mark                MarksService
	recommendationCache RecommendationCacheService
	authenticator       Authenticator
	email               EmailService
	enum                EnumService
	admin               AdminService
	report              ReportService
}

func NewHandler(
	student *student.Service,
	authenticator Authenticator,
	email EmailService,
	enum EnumService,
	admin AdminService,
	recommendationCache RecommendationCacheService,
) *StudentHandler {
	return &StudentHandler{
		student:             student,
		dissertation:        student,
		scientific:          student,
		load:                student,
		authenticator:       authenticator,
		mark:                student,
		email:               email,
		enum:                enum,
		admin:               admin,
		report:              student,
		recommendationCache: recommendationCache,
	}
}

func (h *StudentHandler) authenticate(ctx *gin.Context) (*model.Users, error) {
	token := helpers.GetToken(ctx)

	user, err := h.authenticator.AuthenticateWithUserType(ctx, token, model.UserType_Student.String())
	if err != nil {
		// Если администратор, возвращаем пользователя без ошибки
		if errors.Is(err, models.ErrWrongUserType) {
			user, err = h.authenticator.AuthenticateWithUserType(ctx, token, model.UserType_Admin.String())
			if err == nil {
				return user, nil
			}
		}
		return nil, err
	}

	return user, nil
}

// ValidateToken проверяет, что переданный в запросе токен есть в таблице authorization_token
// Если токен невалиден — прервёт контекст с соответствующим статусом
func (h *StudentHandler) ValidateToken(ctx *gin.Context) error {
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

type PgRecommendationCache struct {
	DB *pgxpool.Pool
}

func NewPgRecommendationCache(db *pgxpool.Pool) *PgRecommendationCache {
	return &PgRecommendationCache{DB: db}
}

func (c *PgRecommendationCache) GetCachedRecommendations(
	ctx context.Context, studentID uuid.UUID, semester int32, topN int,
) ([]Recommendation, *time.Time, error) {
	const query = `
		SELECT recommendations
		FROM recommendations_cache
		WHERE student_id = $1 AND semester = $2 AND top_n = $3
		LIMIT 1
	`
	var data []byte
	var updatedAt time.Time
	err := c.DB.QueryRow(ctx, query, studentID, semester, topN).Scan(&data, &updatedAt)
	if err != nil {
		// Если кэша нет — это не ошибка, просто возвращаем nil, nil
		if err.Error() == "no rows in result set" { // pgx.ErrNoRows
			return nil, nil, nil
		}
		return nil, nil, err
	}

	var recs []Recommendation
	if err := json.Unmarshal(data, &recs); err != nil {
		return nil, nil, err
	}
	return recs, &updatedAt, nil
}

func (c *PgRecommendationCache) SaveCachedRecommendations(
	ctx context.Context, studentID uuid.UUID, semester int32, topN int, recs []Recommendation,
) error {
	const query = `
		INSERT INTO recommendations_cache(student_id, semester, top_n, recommendations, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (student_id, semester, top_n)
		DO UPDATE SET recommendations = EXCLUDED.recommendations, updated_at = EXCLUDED.updated_at
	`
	data, err := json.Marshal(recs)
	if err != nil {
		return err
	}
	_, err = c.DB.Exec(ctx, query, studentID, semester, topN, data, time.Now())
	return err
}
