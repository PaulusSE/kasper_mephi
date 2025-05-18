package admin

import (
	"context"
	"math/rand"
	"net/mail"
	"strings"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/administator_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) GetStudentSupervisorPairs(ctx context.Context) ([]models.StudentSupervisorPair, error) {
	pairs := make([]models.StudentSupervisorPair, 0)

	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		dPairs, err := s.usersRepo.GetStudentSupervisorPairsTx(ctx, tx)
		if err != nil {
			return err
		}

		pairs = dPairs

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "GetStudentSupervisorPairs()")
	}

	return pairs, nil
}

func (s *Service) ChangeSupervisor(ctx context.Context, pairs []models.ChangeSupervisor) error {
	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		for _, pair := range pairs {
			err := s.usersRepo.SetNewSupervisorTx(ctx, tx, pair.StudentID, pair.SupervisorID)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "ChangeSupervisor()")
	}

	return nil
}

func (s *Service) SetStudentFlags(ctx context.Context, students []models.SetStudentsFlags) error {
	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		for _, student := range students {
			var dStatus model.StudentStatus
			if err := dStatus.Scan(strings.TrimSpace(student.StudyingStatus)); err != nil {
				return err
			}

			if err := s.usersRepo.SetStudentFlags(ctx, tx, dStatus, student.CanEdit, student.StudentID); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "SetStudentFlags()")
	}

	return nil
}

func (s *Service) GetSupervisors(ctx context.Context) ([]models.Supervisor, error) {
	sups := make([]models.Supervisor, 0)

	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		dSups, err := s.usersRepo.GetSupervisorsTx(ctx, tx)
		if err != nil {
			return err
		}

		sups = dSups

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "GetSupervisors()")
	}

	return sups, nil
}

func (s *Service) GetStudentsList(ctx context.Context) ([]models.Student, error) {
	students := make([]models.Student, 0)

	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		dStudents, err := s.usersRepo.GetStudentsList(ctx, tx)
		if err != nil {
			return err
		}

		students = dStudents

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "GetSupervisorsStudents()")
	}

	return students, nil
}

func (s *Service) UpsertAttestationMarks(ctx context.Context, marks []models.AttestationMarkRequest) error {
	dMarks := make([]model.Marks, 0, len(marks))
	for _, mark := range marks {
		if mark.Mark < 0 {
			return errors.Wrap(models.ErrInvalidValue, "UpsertAttestationMarks()")
		}

		dMark := model.Marks{
			StudentID: mark.StudentID,
			Mark:      mark.Mark,
			Semester:  mark.Semester,
		}

		dMarks = append(dMarks, dMark)
	}

	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		if err := s.marksRepo.UpsertAttestationMarksTx(ctx, tx, dMarks); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "UpsertAttestationMarks()")
	}

	return nil
}

func (s *Service) AddUsers(ctx context.Context, users request_models.AddUsersRequest, userType model.UserType) ([]models.UsersCredentials, error) {
	strEmails := users.UsersString

	if strings.ContainsAny(strEmails, ",") || strings.ContainsAny(strEmails, "\n") {
		return nil, errors.Wrap(models.ErrInvalidFormat, "AddUsers()")
	}

	emails := strings.Split(strEmails, ";")
	userCreds := make([]models.UsersCredentials, 0, len(emails))
	domainUsers := make([]model.Users, 0, len(emails))

	for i := 0; i < len(emails); i++ {
		email := strings.TrimSpace(emails[i])

		if email == "" {
			continue
		}

		_, err := mail.ParseAddress(email)
		if err != nil {
			return nil, errors.Wrap(err, "AddUsers()")
		}

		password := randPassword(passwordLength)

		userCred := models.UsersCredentials{
			Email:    email,
			Password: password,
		}

		userCreds = append(userCreds, userCred)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		domainUser := model.Users{
			UserID:     uuid.New(),
			Email:      email,
			Password:   string(hashedPassword),
			KasperID:   uuid.New(),
			UserType:   userType,
			Registered: false,
		}

		domainUsers = append(domainUsers, domainUser)
	}

	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		err := s.clientRepo.InsertUsersTx(ctx, tx, domainUsers)

		return err
	}); err != nil {
		return nil, errors.Wrap(err, "AddUsers()")
	}

	return userCreds, nil
}

func (s *Service) ArchiveSupervisor(ctx context.Context, supervisors []models.SupervisorStatus) error {
	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		return s.usersRepo.ArchiveSupervisor(ctx, tx, supervisors)
	}); err != nil {
		return errors.Wrap(err, "GetSupervisorsStudents()")
	}

	return nil
}

func (s *Service) GetNotRegisteredUsers(ctx context.Context) ([]models.UserInfo, error) {
	users := make([]models.UserInfo, 0)
	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		var err error
		users, err = s.clientRepo.GetNotRegisteredUsers(ctx, tx)
		return err
	}); err != nil {
		return nil, errors.Wrap(err, "GetNotRegisteredUsers()")
	}

	return users, nil
}

func (s *Service) DeleteNotRegisteredUsers(ctx context.Context, userIDs []uuid.UUID) error {
	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		err := s.clientRepo.DeleteNotRegisteredUsers(ctx, tx, userIDs)

		return err
	}); err != nil {
		return errors.Wrap(err, "DeleteNotRegisteredUsers()")
	}

	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const passwordLength = 10

func randPassword(passwordLength int) string {
	b := make([]byte, passwordLength)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
