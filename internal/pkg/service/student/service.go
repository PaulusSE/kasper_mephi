package student

import (
	"context"
	"time"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/pkg/models"
	"uir_draft/internal/pkg/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	TokenRepository interface {
		GetUserIDByTokenTx(ctx context.Context, tx pgx.Tx, token string) (uuid.UUID, error)
	}

	UsersRepository interface {
		GetUserTx(ctx context.Context, tx pgx.Tx, userID uuid.UUID) (model.Users, error)
		SetUserRegisteredTx(ctx context.Context, tx pgx.Tx, userID uuid.UUID) error
		ChangeUsersEmail(ctx context.Context, tx pgx.Tx, userID uuid.UUID, email string) error
	}

	MarksRepository interface {
		GetStudentsAttestationMarksTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.Marks, error)
		GetStudentsExamResults(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.Exams, error)
		UpsertExamResults(ctx context.Context, tx pgx.Tx, models []model.Exams) error
		GetStudentsSupervisorMarks(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.SupervisorMarks, error)
		DeleteExamMark(ctx context.Context, tx pgx.Tx, semester int32, ids []uuid.UUID) error
	}

	StudentRepository interface {
		GetStudentTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) (model.Students, error)
		SetStudentStatusTx(ctx context.Context, tx pgx.Tx, status model.ApprovalStatus, studyingStatus model.StudentStatus, studentID uuid.UUID) error
		GetStudentStatusTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) (models.Student, error)
		InsertStudentTx(ctx context.Context, tx pgx.Tx, student model.Students) error
		SetNewSupervisorTx(ctx context.Context, tx pgx.Tx, studentID, supervisorID uuid.UUID) error
		GetAllStudentsSupervisors(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]models.SupervisorFull, error)
		GetStudentProfile(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) (models.StudentProfile, error)
		UpdateStudent(ctx context.Context, tx pgx.Tx, student model.Students) error
		GetStudentsActualSupervisorTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) (models.Supervisor, error)
	}

	DissertationRepository interface {
		SetSemesterProgressStatusTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID, status model.ApprovalStatus, acceptedAt *time.Time) error
		SetDissertationStatusTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID, status model.ApprovalStatus, semester int32) error
		SetDissertationTitleStatusTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID, status model.ApprovalStatus, semester int32, acceptedAt *time.Time) error

		GetSemesterProgressTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.SemesterProgress, error)
		UpsertSemesterProgressTx(ctx context.Context, tx pgx.Tx, progresses []model.SemesterProgress) error

		GetDissertationDataBySemester(ctx context.Context, tx pgx.Tx, studentID uuid.UUID, semester int32) (model.Dissertations, error)
		GetDissertationsTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.Dissertations, error)
		UpsertDissertationTx(ctx context.Context, tx pgx.Tx, model model.Dissertations) error

		GetDissertationTitlesTx(ctx context.Context, tx pgx.Tx, dissertationID uuid.UUID) ([]model.DissertationTitles, error)
		InsertDissertationTitleTx(ctx context.Context, tx pgx.Tx, title model.DissertationTitles) error

		GetFeedbackTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.Feedback, error)

		GetStudentsProgressiveness(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.Progressiveness, error)
		UpsertStudentsProgressiveness(ctx context.Context, tx pgx.Tx, progress model.Progressiveness) error
	}

	ScientificRepository interface {
		SetScientificWorkStatusTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID, status model.ApprovalStatus, semester int32, acceptedAt *time.Time) error

		GetScientificWorksStatusTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.ScientificWorksStatus, error)

		InsertPublicationsTx(ctx context.Context, tx pgx.Tx, publications []model.Publications) error
		UpdatePublicationsTx(ctx context.Context, tx pgx.Tx, publications []model.Publications) error
		DeletePublicationsTx(ctx context.Context, tx pgx.Tx, publicationsIDs []uuid.UUID) error

		InsertConferencesTx(ctx context.Context, tx pgx.Tx, conferences []model.Conferences) error
		UpdateConferencesTx(ctx context.Context, tx pgx.Tx, conferences []model.Conferences) error
		DeleteConferencesTx(ctx context.Context, tx pgx.Tx, conferencesIDs []uuid.UUID) error

		InsertResearchProjectsTx(ctx context.Context, tx pgx.Tx, projects []model.ResearchProjects) error
		UpdateResearchProjectsTx(ctx context.Context, tx pgx.Tx, projects []model.ResearchProjects) error
		DeleteResearchProjectsTx(ctx context.Context, tx pgx.Tx, projectsIDs []uuid.UUID) error

		InitScientificWorkStatusTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID, years int32) error
		GetScientificWorksStatusBySemesterTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID, semester int32) (model.ScientificWorksStatus, error)

		GetPublicationsTx(ctx context.Context, tx pgx.Tx, worksIDs []uuid.UUID) ([]model.Publications, error)
		GetConferencesTx(ctx context.Context, tx pgx.Tx, worksIDs []uuid.UUID) ([]model.Conferences, error)
		GetResearchProjectsTx(ctx context.Context, tx pgx.Tx, worksIDs []uuid.UUID) ([]model.ResearchProjects, error)
		GetScientificWorksStatusIDs(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]uuid.UUID, error)

		InsertPatents(ctx context.Context, tx pgx.Tx, patents []model.Patents) error
		GetPatents(ctx context.Context, tx pgx.Tx, worksIDs []uuid.UUID) ([]model.Patents, error)
		UpdatePatents(ctx context.Context, tx pgx.Tx, patents []model.Patents) error
		DeletePatents(ctx context.Context, tx pgx.Tx, patentIDs []uuid.UUID) error
	}

	TeachingLoadRepository interface {
		SetTeachingLoadStatusTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID, status model.ApprovalStatus, semester int32, acceptedAt *time.Time) error

		GetTeachingLoadStatusTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.TeachingLoadStatus, error)

		InsertClassroomLoadsTx(ctx context.Context, tx pgx.Tx, loads []model.ClassroomLoad) error
		UpdateClassroomLoadsTx(ctx context.Context, tx pgx.Tx, loads []model.ClassroomLoad) error
		DeleteClassroomLoadsTx(ctx context.Context, tx pgx.Tx, classroomsIDs []uuid.UUID) error

		InsertIndividualLoadsTx(ctx context.Context, tx pgx.Tx, loads []model.IndividualStudentsLoad) error
		UpdateIndividualLoadsTx(ctx context.Context, tx pgx.Tx, loads []model.IndividualStudentsLoad) error
		DeleteIndividualStudentsLoadsTx(ctx context.Context, tx pgx.Tx, individualsIDs []uuid.UUID) error

		InsertAdditionalLoadsTx(ctx context.Context, tx pgx.Tx, loads []model.AdditionalLoad) error
		UpdateAdditionalLoadsTx(ctx context.Context, tx pgx.Tx, loads []model.AdditionalLoad) error
		DeleteAdditionalLoadsTx(ctx context.Context, tx pgx.Tx, additionalIDs []uuid.UUID) error

		InitTeachingLoadsStatusTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID, years int32) error
		GetTeachingLoadStatusBySemesterTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID, semester int32) (model.TeachingLoadStatus, error)

		GetClassroomLoadsTx(ctx context.Context, tx pgx.Tx, loadsIDs []uuid.UUID) ([]model.ClassroomLoad, error)
		GetAdditionalLoadsTx(ctx context.Context, tx pgx.Tx, loadsIDs []uuid.UUID) ([]model.AdditionalLoad, error)
		GetIndividualLoadsTx(ctx context.Context, tx pgx.Tx, loadsIDs []uuid.UUID) ([]model.IndividualStudentsLoad, error)
		GetTeachingLoadStatusIDs(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]uuid.UUID, error)
	}

	CommentRepository interface {
		UpsertStudentsComment(ctx context.Context, tx pgx.Tx, comment model.StudentsCommentary) error
		GetStudentsCommentaries(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.StudentsCommentary, error)

		UpsertDissertationComment(ctx context.Context, tx pgx.Tx, comment model.DissertationCommentary) error
		GetDissertationComments(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.DissertationCommentary, error)

		UpsertPlanComment(ctx context.Context, tx pgx.Tx, plan model.DissertationPlans) error
		GetPlanComments(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]model.DissertationPlans, error)
	}
)

type Service struct {
	dissertationRepo DissertationRepository
	loadRepo         TeachingLoadRepository
	scienceRepo      ScientificRepository
	marksRepo        MarksRepository
	studRepo         StudentRepository
	tokenRepo        TokenRepository
	userRepo         UsersRepository
	commentRepo      CommentRepository
	db               *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		dissertationRepo: repository.NewDissertationRepository(),
		loadRepo:         repository.NewTeachingLoadRepository(),
		scienceRepo:      repository.NewScientificRepository(),
		marksRepo:        repository.NewMarksRepository(),
		studRepo:         repository.NewClientRepository(),
		tokenRepo:        repository.NewTokenRepository(),
		userRepo:         repository.NewUsersRepository(),
		commentRepo:      repository.NewCommentaryRepository(),
		db:               db,
	}
}
