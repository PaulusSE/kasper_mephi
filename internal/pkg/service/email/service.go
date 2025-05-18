package email

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

type (
	UsersRepository interface {
		GetUserByKasperIDTx(ctx context.Context, tx pgx.Tx, kasperID uuid.UUID) (model.Users, error)
	}

	ClientRepository interface {
		GetStudentTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) (model.Students, error)
		GetStudentsActualSupervisorTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) (models.Supervisor, error)
		GetSupervisorTx(ctx context.Context, tx pgx.Tx, supervisorID uuid.UUID) (models.Supervisor, error)
	}
)

type Service struct {
	sender     string
	password   string
	host       string
	db         *pgxpool.Pool
	userRepo   UsersRepository
	clientRepo ClientRepository
}

func NewService(
	sender string,
	password string,
	host string,
	db *pgxpool.Pool,
	userRepo UsersRepository,
	clientRepo ClientRepository,
) *Service {
	return &Service{
		sender:     sender,
		password:   password,
		host:       host,
		db:         db,
		userRepo:   userRepo,
		clientRepo: clientRepo,
	}
}

type SupervisorData struct {
	SupervisorName string
	StudentName    string
	Status         string
	Type           string
	Date           string
}

type StudentData struct {
	SupervisorName string
	StudentName    string
	Type           string
	Date           string
}

type InviteData struct {
	Email    string
	Password string
}

var approveStatusMap = map[model.ApprovalStatus]string{
	model.ApprovalStatus_Todo:       "На доработку",
	model.ApprovalStatus_OnReview:   "Ожидает проверки",
	model.ApprovalStatus_Failed:     "Отклонено",
	model.ApprovalStatus_Approved:   "Принято",
	model.ApprovalStatus_InProgress: "В процессе",
}

func (s *Service) SendInviteEmails(_ context.Context, credentials []models.UsersCredentials, templatePath string) error {
	mails := make([]*gomail.Message, 0, len(credentials))

	for _, cred := range credentials {
		var body bytes.Buffer
		data := InviteData{
			Email:    cred.Email,
			Password: cred.Password,
		}
		t, err := template.ParseFiles(templatePath)
		if err != nil {
			return errors.Wrap(err, "parsing template")
		}

		_ = t.Execute(&body, data)

		m := gomail.NewMessage()
		m.SetHeader("From", s.sender)
		m.SetHeader("To", cred.Email)
		m.SetHeader("Subject", "Приглашение на регистрацию в системе учета деятельности аспирантов")
		m.SetBody("text/html", body.String())

		mails = append(mails, m)
	}

	d := gomail.NewDialer(s.host, 587, s.sender, s.password)

	if err := d.DialAndSend(mails...); err != nil {
		return errors.Wrap(err, "sending email")
	}

	return nil
}

func (s *Service) SendMailToStudent(ctx context.Context, studentID, supervisorID uuid.UUID, templatePath, tt, status string) error {
	var data SupervisorData
	var (
		receiver string
	)

	var approveStatus model.ApprovalStatus
	if err := approveStatus.Scan(strings.TrimSpace(status)); err != nil {
		return errors.Wrap(err, models.ErrInvalidEnumValue.Error())
	}

	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		uStud, err := s.userRepo.GetUserByKasperIDTx(ctx, tx, studentID)
		if err != nil {
			return err
		}

		stud, err := s.clientRepo.GetStudentTx(ctx, tx, studentID)
		if err != nil {
			return err
		}

		sup, err := s.clientRepo.GetSupervisorTx(ctx, tx, supervisorID)
		if err != nil {
			return err
		}

		data = SupervisorData{
			SupervisorName: sup.FullName,
			StudentName:    stud.FullName,
			Status:         approveStatusMap[approveStatus],
			Type:           tt,
			Date:           time.Now().Format("02.01.2006"),
		}

		receiver = uStud.Email

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "SendMailToStudent()")
	}

	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return errors.Wrap(err, "parsing template")
	}

	_ = t.Execute(&body, data)

	m := gomail.NewMessage()
	m.SetHeader("From", s.sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", fmt.Sprintf("Уведомление от Научного руководителя %s", data.SupervisorName))
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(s.host, 587, s.sender, s.password)

	if err = d.DialAndSend(m); err != nil {
		return errors.Wrap(err, "sending email")
	}

	return nil
}

func (s *Service) SendMailToSupervisor(ctx context.Context, studentID uuid.UUID, templatePath, tt string) error {
	var data StudentData
	var (
		receiver string
	)

	err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		supervisor, err := s.clientRepo.GetStudentsActualSupervisorTx(ctx, tx, studentID)
		if err != nil {
			return err
		}

		userSupervisor, err := s.userRepo.GetUserByKasperIDTx(ctx, tx, supervisor.SupervisorID)
		if err != nil {
			return err
		}

		student, err := s.clientRepo.GetStudentTx(ctx, tx, studentID)
		if err != nil {
			return err
		}

		data = StudentData{
			SupervisorName: supervisor.FullName,
			StudentName:    student.FullName,
			Type:           tt,
			Date:           time.Now().Format("02.01.2006"),
		}

		receiver = userSupervisor.Email
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "SendMailToStudent()")
	}
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return errors.Wrap(err, "parsing template")
	}

	_ = t.Execute(&body, data)

	mail := gomail.NewMessage()
	mail.SetHeader("From", s.sender)
	mail.SetHeader("To", receiver)
	mail.SetHeader("Subject", fmt.Sprintf("Уведомление от Студента %s", data.StudentName))
	mail.SetBody("text/html", body.String())

	d := gomail.NewDialer(s.host, 587, s.sender, s.password)

	if err = d.DialAndSend(mail); err != nil {
		return errors.Wrap(err, "sending email")
	}

	return nil
}
