package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"uir_draft/internal/generated/kasper/uir_draft/public/model"
	"uir_draft/internal/pkg/service/student/mapping"
)

func (s *Service) GetTeachingLoad(ctx context.Context, token string, studentID uuid.UUID) (*mapping.TeachingLoad, error) {
	session, err := s.tokenRepo.Authenticate(ctx, token, s.db)
	if err != nil {
		return nil, errors.Wrap(err, "[Student]")
	}

	if session.TokenStatus != model.TokenStatus_Active {
		return nil, ErrNonValidToken
	}

	load, err := s.loadRepo.GetStudentsTeachingLoad(ctx, s.db, studentID)
	if err != nil {
		return nil, errors.Wrap(err, "[Student]")
	}

	loads := mapping.MapTeachingLoadFromDomain(load)

	return &loads, nil
}
