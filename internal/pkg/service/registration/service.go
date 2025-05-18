package registration

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"uir_draft/internal/handlers/authorization_handler/request_models"
	"uir_draft/internal/pkg/repository"
)

type Service struct {
	db   *pgxpool.Pool
	repo repository.RegistrationRequestsRepository
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		db:   db,
		repo: repository.NewRegistrationRequestsRepository(),
	}
}

// CreateStudentRequest сохраняет форму регистрации в registration_requests
func (s *Service) CreateStudentRequest(ctx context.Context, hashedPassword []byte, data request_models.FirstStudentRegistry) error {
	// упаковываем все поля кроме email/password в JSON
	payload, err := json.Marshal(map[string]interface{}{
		"full_name":         data.FullName,
		"group_number":      data.GroupID,
		"specialization_id": data.SpecializationID,
		"actual_semester":   data.ActualSemester,
		"start_date":        data.StartDate,
		"phone":             data.Phone,
		"number_of_years":   data.NumberOfYears,
		"supervisor_id":     data.SupervisorID,
		"category":          data.Category,
	})
	if err != nil {
		return err
	}

	req := &repository.RegistrationRequest{
		RequestID:    uuid.New(),
		Email:        data.Email,
		PasswordHash: string(hashedPassword),
		UserType:     "student",
		Payload:      payload,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	return s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		return s.repo.InsertRequestTx(ctx, tx, req)
	})
}

// CreateSupervisorRequest сохраняет форму регистрации в registration_requests
func (s *Service) CreateSupervisorRequest(ctx context.Context, hashedPassword []byte, data request_models.FirstSupervisorRegistry) error {
	payload, err := json.Marshal(map[string]interface{}{
		"full_name": data.FullName,
		"phone":     data.Phone,
	})
	if err != nil {
		return err
	}

	req := &repository.RegistrationRequest{
		RequestID:    uuid.New(),
		Email:        *data.Email,
		PasswordHash: string(hashedPassword),
		UserType:     "supervisor",
		Payload:      payload,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	return s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		if err := s.repo.InsertRequestTx(ctx, tx, req); err != nil {
			// временный лог
			log.Printf("CreateSupervisorRequest: insert failed: %v", err)
			return err
		}
		return nil
	})
}

func (s *Service) ListRequests(ctx context.Context, userType string) ([]repository.RegistrationRequest, error) {
	var reqs []repository.RegistrationRequest
	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		var err error
		reqs, err = s.repo.ListRequestsTx(ctx, tx, userType)
		return err
	})
	return reqs, err
}

func (s *Service) ReviewRequest(ctx context.Context, requestID uuid.UUID, approve bool) error {
	status := "rejected"
	if approve {
		status = "approved"
	}
	return s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		return s.repo.UpdateStatusTx(ctx, tx, requestID, status)
	})
}

func (s *Service) GetRequestByID(ctx context.Context, id uuid.UUID) (repository.RegistrationRequest, error) {
	var req repository.RegistrationRequest
	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		var err error
		req, err = s.repo.GetByIDTx(ctx, tx, id)
		return err
	})
	return req, err
}
