package request_models

import (
	"uir_draft/internal/pkg/models"

	"github.com/google/uuid"
)

type UpsertFeedbackRequest struct {
	Feedback  models.FeedbackRequest `json:"feedback"`
	StudentID uuid.UUID              `json:"student_id,omitempty"`
}

type GetByStudentID struct {
	StudentID uuid.UUID `json:"student_id,omitempty"`
}

type ToStatusRequest struct {
	StudentID uuid.UUID `json:"student_id,omitempty"`
	Status    string    `json:"status,omitempty"`
}

type DownloadDissertationRequestSup struct {
	// ID студента
	StudentID uuid.UUID `json:"student_id,omitempty"`
	// Семестр
	Semester int32 `json:"semester,omitempty"`
}

type DownloadDissertationRequestSupOnlySemester struct {
	Semester int32 `json:"semester,omitempty"`
}

type UpsertSupervisorMarkRequest struct {
	StudentID uuid.UUID `json:"student_id,omitempty"`
	Semester  int32     `json:"semester,omitempty"`
	Mark      int32     `json:"mark,omitempty"`
}
