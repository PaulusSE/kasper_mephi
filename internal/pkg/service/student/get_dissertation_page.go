package student

import (
	"context"

	"uir_draft/internal/generated/kasper/uir_draft/public/model"
	"uir_draft/internal/pkg/models"
	"uir_draft/internal/pkg/service/student/mapping"

	"github.com/pkg/errors"
)

type StudentDissertationPlan struct {
	First  bool `json:"id1,omitempty"`
	Second bool `json:"id2,omitempty"`
	Third  bool `json:"id3,omitempty"`
	Forth  bool `json:"id4,omitempty"`
	Fifth  bool `json:"id5,omitempty"`
	Sixth  bool `json:"id6,omitempty"`
}

type DissertationPage struct {
	DissertationPlan map[string]*StudentDissertationPlan `json:"dissertationPlan"`
	CommonInfo       models.StudentCommonInformation     `json:"commonInfo"`
	IDs              []*mapping.DissertationIDs          `json:"ids"`
}

func (s *Service) GetDissertationPage(ctx context.Context, token string) (*DissertationPage, error) {
	session, err := s.tokenRepo.Authenticate(ctx, token, s.db)
	if err != nil {
		return nil, errors.Wrap(err, "authentication error")
	}

	if session.TokenStatus != model.TokenStatus_Active {
		return nil, ErrNonValidToken
	}

	commonInfo, err := s.studRepo.GetStudentCommonInfo(ctx, s.db, session.KasperID)
	if err != nil {
		return nil, errors.Wrap(err, "GetStudentCommonInfo()")
	}

	plans, err := s.semesterRepo.GetStudentDissertationPlan(ctx, s.db, session.KasperID)
	if err != nil {
		return nil, errors.Wrap(err, "GetStudentDissertationPlan()")
	}

	domainIDs, err := s.dRepo.GetDissertationIDs(ctx, s.db, session.KasperID)
	if err != nil {
		return nil, err
	}

	ids := mapping.MapIDsFromDomain(domainIDs)

	planMap := make(map[string]*StudentDissertationPlan, len(plans))

	for _, semester := range plans {
		plan := &StudentDissertationPlan{
			First:  semester.First,
			Second: semester.Second,
			Third:  semester.Third,
			Forth:  semester.Forth,
			Fifth:  semester.Fifth,
			Sixth:  semester.Sixth,
		}
		planMap[semester.Name] = plan
	}

	return &DissertationPage{
		DissertationPlan: planMap,
		CommonInfo:       *commonInfo,
		IDs:              ids,
	}, nil
}
