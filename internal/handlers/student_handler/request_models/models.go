package request_models

import (
	"mime/multipart"

	"uir_draft/internal/pkg/models"

	"github.com/google/uuid"
)

type DeleteIDs struct {
	IDs      []uuid.UUID `json:"ids,omitempty"`
	Semester int32       `json:"semester,omitempty"`
}

type UpsertAdditionalLoadRequest struct {
	// Дополнительные нагрузки, которые нужно добавить/изменить. Если запрос на добавление, то
	// ID дополнительной нагрузки (load_id) не заполнять (делать null), в противном случае - иначе
	Loads []models.AdditionalLoad `json:"loads,omitempty"`
	// Семестр
	Semester int32 `json:"semester,omitempty"`
}

type UpsertClassroomLoadRequest struct {
	// Аудиторные нагрузки, которые нужно добавить/изменить. Если запрос на добавление, то
	// ID нагрузки (load_id) не заполнять (делать null), в противном случае - иначе
	Loads []models.ClassroomLoad `json:"loads,omitempty"`
	// Семестр
	Semester int32 `json:"semester,omitempty"`
}

type UpsertConferencesRequest struct {
	// Научные конференции, которые нужно добавить/изменить. Если запрос на добавление, то
	// ID конференции (conference_id) не заполнять (делать null), в противном случае - иначе
	Conferences []models.Conference `json:"conferences,omitempty"`
	// Семестр
	Semester int32 `json:"semester,omitempty"`
}

type UpsertIndividualLoadRequest struct {
	// Индивидуальные нагрузки, которые нужно добавить/изменить. Если запрос на добавление, то
	// ID индивидуальной нагрузки (load_id) не заполнять (делать null), в противном случае - иначе
	Loads    []models.IndividualStudentsLoad `json:"loads,omitempty"`
	Semester int32                           `json:"semester,omitempty"`
}

type UpsertPublicationsRequest struct {
	// Научные публикации, которые нужно добавить/изменить. Если запрос на добавление, то
	// ID публикации (publication_id) не заполнять (делать null), в противном случае - иначе
	Publications []models.Publication `json:"publications,omitempty"`
	Semester     int32                `json:"semester,omitempty"`
}

type UpsertPatentsRequest struct {
	// Патенты, которые нужно добавить/изменить. Если запрос на добавление, то
	// ID патента (patent_id) не заполнять (делать null), в противном случае - иначе
	Patents  []models.Patent `json:"patents,omitempty"`
	Semester int32           `json:"semester,omitempty"`
}

type UpsertResearchProjectsRequest struct {
	// Научные проекты, которые нужно добавить/изменить. Если запрос на добавление, то
	// ID проекта (project_id) не заполнять (делать null), в противном случае - иначе
	Projects []models.ResearchProject `json:"projects,omitempty"`
	Semester int32                    `json:"semester,omitempty"`
}

type UpsertProgressRequest struct {
	Progresses []models.SemesterProgressRequest `json:"progresses"`
}

type DownloadDissertationRequest struct {
	// Семестр
	Semester int32 `json:"semester,omitempty"`
}

type DownloadPresentationRequest struct {
	// Семестр
	Semester int32 `json:"semester,omitempty"`
}

type UploadDissertationRequest struct {
	// Семестр
	Semester int32                 `form:"semester" binding:"required"`
	File     *multipart.FileHeader `form:"upload" binding:"required" swaggerignore:"true"`
}

type ToReviewRequest struct {
	Semester int32 `json:"semester"`
}

type UpsertDissertationTitleRequest struct {
	Title           string `json:"title,omitempty"`
	ResearchObject  string `json:"research_object"`
	ResearchSubject string `json:"research_subject"`
}

type UpdateProgressivenessRequest struct {
	Progress int32 `json:"progress"`
}

type AllToReviewRequest struct {
	Commentary *string `json:"commentary"`
}

type UpsertReportCommentsRequest struct {
	DissertationComment models.DissertationCommentRequest `json:"dissertation_comment"`
	DissertationPlan    models.DissertationPlanRequest    `json:"dissertation_plan"`
	Semester            int32                             `json:"semester"`
}

type UpsertExamResultsRequest struct {
	Marks []models.ExamRequest `json:"marks,omitempty"`
}
