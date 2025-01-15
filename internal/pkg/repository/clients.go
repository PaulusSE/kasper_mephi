package repository

import (
	"context"
	"log"
	"time"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/generated/new_kasper/new_uir/public/table"
	"uir_draft/internal/pkg/models"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type ClientRepository struct{}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{}
}

func (r *ClientRepository) GetStudentTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) (model.Students, error) {
	stmt, args := table.Students.
		SELECT(table.Students.AllColumns).
		WHERE(table.Students.StudentID.EQ(postgres.UUID(studentID))).
		Sql()

	row := tx.QueryRow(ctx, stmt, args...)
	student := model.Students{}

	if err := scanStudent(row, &student); err != nil {
		return model.Students{}, errors.Wrap(err, "GetStudentTx()")
	}

	return student, nil
}

func (r *ClientRepository) GetStudentStatusTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) (models.Student, error) {
	stmt, args := table.Students.
		SELECT(
			table.Students.AllColumns.Except(table.Students.UserID, table.Students.SpecID, table.Students.GroupID),
			table.Specializations.Title,
			table.Groups.GroupName,
		).
		FROM(table.Students.
			LEFT_JOIN(table.Groups, table.Students.GroupID.EQ(table.Groups.GroupID)).
			LEFT_JOIN(table.Specializations, table.Students.SpecID.EQ(table.Specializations.SpecID)),
		).
		WHERE(table.Students.StudentID.EQ(postgres.UUID(studentID))).
		Sql()

	log.Printf("Executing SQL: %s with args: %v", stmt, args)

	row := tx.QueryRow(ctx, stmt, args...)
	student := models.Student{}
	if err := scanStudentStatus(row, &student); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Student{}, errors.New("no student found with the given ID")
		}
		return models.Student{}, errors.Wrap(err, "GetStudentStatusTx()")
	}

	return student, nil
}

func (r *ClientRepository) GetStudentsList(ctx context.Context, tx pgx.Tx) ([]models.Student, error) {
	stmt, args := table.Students.
		SELECT(
			table.Students.AllColumns.Except(table.Students.UserID, table.Students.SpecID, table.Students.GroupID),
			table.Specializations.Title,
			table.Groups.GroupName,
		).
		FROM(table.Students.
			INNER_JOIN(table.Groups, table.Students.GroupID.EQ(table.Groups.GroupID)).
			INNER_JOIN(table.Specializations, table.Students.SpecID.EQ(table.Specializations.SpecID)),
		).
		Sql()

	rows, err := tx.Query(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	students := make([]models.Student, 0)
	for rows.Next() {
		student := models.Student{}
		if err := scanStudentStatus(rows, &student); err != nil {
			return nil, errors.Wrap(err, "GetStudentStatusTx()")
		}

		students = append(students, student)
	}

	return students, nil
}

func (r *ClientRepository) InsertStudentTx(ctx context.Context, tx pgx.Tx, student model.Students) error {
	stmt, args := table.Students.
		INSERT(
			table.Students.AllColumns.
				Except(
					table.Students.StudyingStatus,
					table.Students.Status,
					table.Students.CanEdit,
				),
		).
		VALUES(
			student.StudentID,
			student.UserID,
			student.FullName,
			//student.Department,
			student.SpecID,
			student.ActualSemester,
			student.Years,
			student.StartDate,
			student.GroupID,
			student.Phone,
			student.Category,
			student.EndDate,
		).
		Sql()

	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		return errors.Wrap(err, "InsertStudentTx()")
	}

	return nil
}

func (r *ClientRepository) SetStudentStatusTx(ctx context.Context, tx pgx.Tx, status model.ApprovalStatus, studyingStatus model.StudentStatus, studentID uuid.UUID) error {
	stmt, args := table.Students.
		UPDATE(
			table.Students.Status,
			table.Students.StudyingStatus,
		).
		SET(
			status,
			studyingStatus,
		).
		WHERE(table.Students.StudentID.EQ(postgres.UUID(studentID))).
		Sql()

	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		return errors.Wrap(err, "SetStudentStatusTx()")
	}

	return nil
}

func (r *ClientRepository) SetStudentFlags(ctx context.Context, tx pgx.Tx, studyingStatus model.StudentStatus, canEdit bool, studentID uuid.UUID) error {
	stmt, args := table.Students.
		UPDATE(
			table.Students.StudyingStatus,
			table.Students.CanEdit,
		).
		SET(
			studyingStatus,
			canEdit,
		).
		WHERE(table.Students.StudentID.EQ(postgres.UUID(studentID))).
		Sql()

	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		return errors.Wrap(err, "SetStudentFlags()")
	}

	return nil
}

func (r *ClientRepository) GetSupervisorsStudentsTx(ctx context.Context, tx pgx.Tx, supervisorID uuid.UUID) ([]models.Student, error) {
	stmt, args := table.Students.
		SELECT(
			table.Students.AllColumns.
				Except(
					table.Students.UserID,
					table.Students.SpecID,
					table.Students.GroupID,
					table.Supervisors.Archived,
				),
			table.Specializations.Title,
			table.Groups.GroupName,
		).DISTINCT(table.Students.StudentID).
		FROM(table.Students.
			INNER_JOIN(table.StudentsSupervisors, table.Students.StudentID.EQ(table.StudentsSupervisors.StudentID)).
			INNER_JOIN(table.Groups, table.Students.GroupID.EQ(table.Groups.GroupID)).
			INNER_JOIN(table.Specializations, table.Students.SpecID.EQ(table.Specializations.SpecID)),
		).
		WHERE(table.StudentsSupervisors.SupervisorID.EQ(postgres.UUID(supervisorID)).
			AND(table.StudentsSupervisors.EndAt.IS_NULL())).
		Sql()

	rows, err := tx.Query(ctx, stmt, args...)
	if err != nil {
		return nil, errors.Wrap(err, "GetSupervisorsStudentsTx()")
	}
	defer rows.Close()

	list := make([]models.Student, 0)

	for rows.Next() {
		el := models.Student{}
		if err = scanStudentStatus(rows, &el); err != nil {
			return nil, errors.Wrap(err, "GetSupervisorsStudentsTx(): scanning rows")
		}

		list = append(list, el)
	}

	return list, err
}

func (r *ClientRepository) GetStudentSupervisorPairsTx(ctx context.Context, tx pgx.Tx) ([]models.StudentSupervisorPair, error) {
	stmt, args := table.Students.
		SELECT(
			table.Students.AllColumns.
				Except(
					table.Students.UserID,
					table.Students.SpecID,
					table.Students.GroupID,
					table.Supervisors.Archived,
					table.Supervisors.Rank,
					table.Supervisors.Position,
				),
			table.Specializations.Title,
			table.Groups.GroupName,
			table.Supervisors.SupervisorID,
			table.Supervisors.FullName,
		).
		FROM(table.Students.
			LEFT_JOIN(table.StudentsSupervisors, table.StudentsSupervisors.StudentID.EQ(table.Students.StudentID)).
			LEFT_JOIN(table.Supervisors, table.StudentsSupervisors.SupervisorID.EQ(table.Supervisors.SupervisorID)).
			INNER_JOIN(table.Groups, table.Students.GroupID.EQ(table.Groups.GroupID)).
			INNER_JOIN(table.Specializations, table.Students.SpecID.EQ(table.Specializations.SpecID)),
		).
		WHERE(table.StudentsSupervisors.EndAt.IS_NULL()).
		Sql()

	rows, err := tx.Query(ctx, stmt, args...)
	if err != nil {
		return nil, errors.Wrap(err, "GetStudentSupervisorPairsTx()")
	}
	defer rows.Close()

	list := make([]models.StudentSupervisorPair, 0)

	for rows.Next() {
		el := models.StudentSupervisorPair{}
		if err = scanStudentSupervisorPair(rows, &el); err != nil {
			return nil, errors.Wrap(err, "GetStudentSupervisorPairsTx(): scanning rows")
		}

		list = append(list, el)
	}

	return list, nil
}

func (r *ClientRepository) SetNewSupervisorTx(ctx context.Context, tx pgx.Tx, studentID, supervisorID uuid.UUID) error {
	insertStmt, instArgs := table.StudentsSupervisors.
		INSERT(
			table.StudentsSupervisors.ID,
			table.StudentsSupervisors.SupervisorID,
			table.StudentsSupervisors.StudentID,
		).
		VALUES(
			uuid.New(),
			supervisorID,
			studentID,
		).
		Sql()

	updateStmt, updArgs := table.StudentsSupervisors.
		UPDATE(table.StudentsSupervisors.EndAt).
		SET(time.Now()).
		WHERE(
			table.StudentsSupervisors.EndAt.IS_NULL().
				AND(table.StudentsSupervisors.StudentID.EQ(postgres.UUID(studentID))),
		).
		Sql()

	_, err := tx.Exec(ctx, updateStmt, updArgs...)
	if err != nil {
		return errors.Wrap(err, "SetNewSupervisorTx(): update")
	}

	_, err = tx.Exec(ctx, insertStmt, instArgs...)
	if err != nil {
		return errors.Wrap(err, "SetNewSupervisorTx(): insert")
	}

	return nil
}

func (r *ClientRepository) GetSupervisorsTx(ctx context.Context, tx pgx.Tx) ([]models.Supervisor, error) {
	stmt, args := table.Supervisors.
		SELECT(
			table.Supervisors.AllColumns.Except(table.Supervisors.UserID),
		).
		Sql()

	rows, err := tx.Query(ctx, stmt, args...)
	if err != nil {
		return nil, errors.Wrap(err, "GetSupervisorsTx()")
	}
	defer rows.Close()

	supervisors := make([]models.Supervisor, 0)

	for rows.Next() {
		el := models.Supervisor{}
		if err = scanSupervisor(rows, &el); err != nil {
			return nil, errors.Wrap(err, "GetSupervisorsTx()")
		}

		supervisors = append(supervisors, el)
	}

	return supervisors, nil
}

func (r *ClientRepository) GetStudentsActualSupervisorTx(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) (models.Supervisor, error) {
	stmt, args := table.Supervisors.
		SELECT(
			table.Supervisors.AllColumns.Except(table.Supervisors.UserID),
		).
		FROM(
			table.StudentsSupervisors.
				INNER_JOIN(table.Supervisors, table.StudentsSupervisors.SupervisorID.EQ(table.Supervisors.SupervisorID)),
		).
		WHERE(table.StudentsSupervisors.StudentID.EQ(postgres.UUID(studentID)).
			AND(table.StudentsSupervisors.EndAt.IS_NULL())).
		Sql()

	row := tx.QueryRow(ctx, stmt, args...)
	supervisor := models.Supervisor{}
	if err := scanSupervisor(row, &supervisor); err != nil {
		return models.Supervisor{}, errors.Wrap(err, "GetStudentsActualSupervisorTx()")
	}

	return supervisor, nil
}

func (r *ClientRepository) GetSupervisorTx(ctx context.Context, tx pgx.Tx, supervisorID uuid.UUID) (models.Supervisor, error) {
	stmt, args := table.Supervisors.
		SELECT(
			table.Supervisors.AllColumns.Except(table.Supervisors.UserID),
		).
		WHERE(table.Supervisors.SupervisorID.EQ(postgres.UUID(supervisorID))).
		Sql()

	row := tx.QueryRow(ctx, stmt, args...)
	supervisor := models.Supervisor{}
	if err := scanSupervisor(row, &supervisor); err != nil {
		return models.Supervisor{}, errors.Wrap(err, "GetSupervisorTx()")
	}

	return supervisor, nil
}

func (r *ClientRepository) GetAllStudentsSupervisors(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) ([]models.SupervisorFull, error) {
	stmt, args := table.Supervisors.
		SELECT(
			table.Supervisors.SupervisorID,
			table.Supervisors.FullName,
			table.StudentsSupervisors.StartAt,
			table.StudentsSupervisors.EndAt,
		).
		FROM(table.Supervisors.
			INNER_JOIN(table.StudentsSupervisors, table.Supervisors.SupervisorID.EQ(table.StudentsSupervisors.SupervisorID))).
		WHERE(table.StudentsSupervisors.StudentID.EQ(postgres.UUID(studentID))).
		Sql()

	rows, err := tx.Query(ctx, stmt, args...)
	if err != nil {
		return nil, errors.Wrap(err, "GetAllStudentsSupervisors()")
	}
	defer rows.Close()

	supervisors := make([]models.SupervisorFull, 0)

	for rows.Next() {
		supervisor := models.SupervisorFull{}
		if err := scanSupervisorFull(rows, &supervisor); err != nil {
			return nil, errors.Wrap(err, "GetAllStudentsSupervisors(): scanning")
		}

		supervisors = append(supervisors, supervisor)
	}

	return supervisors, nil
}

func (r *ClientRepository) GetSupervisorProfile(ctx context.Context, tx pgx.Tx, supervisorID uuid.UUID) (models.SupervisorProfile, error) {
	stmt, args := table.Supervisors.
		SELECT(
			table.Supervisors.AllColumns.Except(table.Supervisors.UserID),
			table.Users.Email,
		).
		FROM(
			table.Supervisors.INNER_JOIN(table.Users, table.Supervisors.SupervisorID.EQ(table.Users.KasperID)),
		).
		WHERE(table.Supervisors.SupervisorID.EQ(postgres.UUID(supervisorID))).
		Sql()

	row := tx.QueryRow(ctx, stmt, args...)
	supervisor := models.SupervisorProfile{}
	if err := scanSupervisorProfile(row, &supervisor); err != nil {
		return models.SupervisorProfile{}, errors.Wrap(err, "GetSupervisorProfile()")
	}

	return supervisor, nil
}

func (r *ClientRepository) GetStudentProfile(ctx context.Context, tx pgx.Tx, studentID uuid.UUID) (models.StudentProfile, error) {
	stmt, args := table.Students.
		SELECT(
			table.Students.AllColumns.Except(table.Students.UserID, table.Students.SpecID, table.Students.GroupID),
			table.Specializations.Title,
			table.Groups.GroupName,
			table.Users.Email,
		).
		FROM(table.Students.
			INNER_JOIN(table.Groups, table.Students.GroupID.EQ(table.Groups.GroupID)).
			INNER_JOIN(table.Specializations, table.Students.SpecID.EQ(table.Specializations.SpecID)).
			INNER_JOIN(table.Users, table.Students.StudentID.EQ(table.Users.KasperID)),
		).
		WHERE(table.Students.StudentID.EQ(postgres.UUID(studentID))).
		Sql()

	row := tx.QueryRow(ctx, stmt, args...)
	student := models.StudentProfile{}
	if err := scanStudentProfile(row, &student); err != nil {
		return models.StudentProfile{}, errors.Wrap(err, "GetStudentProfile()")
	}

	return student, nil
}

func (r *ClientRepository) UpsertSupervisor(ctx context.Context, tx pgx.Tx, supervisor model.Supervisors) error {
	stmt, args := table.Supervisors.
		INSERT(table.Supervisors.AllColumns.Except(table.Supervisors.Archived)).
		MODEL(supervisor).
		ON_CONFLICT(table.Supervisors.SupervisorID).
		DO_UPDATE(postgres.
			SET(
				table.Supervisors.Degree.SET(postgres.String(lo.FromPtr(supervisor.Degree))),
				table.Supervisors.Faculty.SET(postgres.String(lo.FromPtr(supervisor.Faculty))),
				table.Supervisors.Department.SET(postgres.String(lo.FromPtr(supervisor.Department))),
				table.Supervisors.FullName.SET(postgres.String(lo.FromPtr(supervisor.FullName))),
				table.Supervisors.Phone.SET(postgres.String(supervisor.Phone)),
				table.Supervisors.Rank.SET(postgres.String(lo.FromPtr(supervisor.Rank))),
				table.Supervisors.Position.SET(postgres.String(lo.FromPtr(supervisor.Position))),
			),
		).
		Sql()

	_, err := tx.Exec(ctx, stmt, args...)
	if err != nil {
		return errors.Wrap(err, "UpsertSupervisor()")
	}

	return nil
}

func (r *ClientRepository) UpdateStudent(ctx context.Context, tx pgx.Tx, student model.Students) error {
	stmt, args := table.Students.
		UPDATE(
			table.Students.Phone,
			table.Students.FullName,
			table.Students.Years,
			table.Students.Category,
			table.Students.StartDate,
			table.Students.GroupID,
		).
		SET(
			student.Phone,
			student.FullName,
			student.Years,
			student.Category,
			student.StartDate,
			student.GroupID,
		).
		WHERE(table.Students.StudentID.EQ(postgres.UUID(student.StudentID))).
		Sql()

	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		return errors.Wrap(err, "UpdateStudent()")
	}

	return nil
}

func (r *ClientRepository) ArchiveSupervisor(ctx context.Context, tx pgx.Tx, supervisors []models.SupervisorStatus) error {
	for _, supervisor := range supervisors {
		stmt, args := table.Supervisors.
			UPDATE(table.Supervisors.Archived).
			SET(supervisor.Archived).
			WHERE(table.Supervisors.SupervisorID.EQ(postgres.UUID(supervisor.SupervisorID))).
			Sql()

		_, err := tx.Exec(ctx, stmt, args...)
		if err != nil {
			return errors.Wrap(err, "ArchiveSupervisor()")
		}
	}

	return nil
}

func scanSupervisorProfile(row pgx.Row, target *models.SupervisorProfile) error {
	return row.Scan(
		&target.SupervisorID,
		&target.FullName,
		&target.Phone,
		&target.Archived,
		&target.Faculty,
		&target.Department,
		&target.Degree,
		&target.Rank,
		&target.Position,
		&target.Email,
	)
}

func scanStudentProfile(row pgx.Row, target *models.StudentProfile) error {
	return row.Scan(
		&target.StudentID,
		&target.FullName,
		//&target.Department,
		&target.ActualSemester,
		&target.Years,
		&target.StartDate,
		&target.StudyingStatus,
		&target.Status,
		&target.CanEdit,
		&target.Phone,
		&target.Category,
		&target.EndDate,
		&target.Specialization,
		&target.GroupName,
		&target.Email,
	)
}

func scanStudent(row pgx.Row, target *model.Students) error {
	return row.Scan(
		&target.StudentID,
		&target.UserID,
		&target.FullName,
		//&target.Department,
		&target.SpecID,
		&target.ActualSemester,
		&target.Years,
		&target.StartDate,
		&target.StudyingStatus,
		&target.GroupID,
		&target.Status,
		&target.CanEdit,
		&target.Phone,
		&target.Category,
		&target.EndDate,
	)
}

func scanStudentStatus(row pgx.Row, target *models.Student) error {
	return row.Scan(
		&target.StudentID,
		&target.FullName,
		//&target.Department,
		&target.ActualSemester,
		&target.Years,
		&target.StartDate,
		&target.StudyingStatus,
		&target.Status,
		&target.CanEdit,
		&target.Phone,
		&target.Category,
		&target.EndDate,
		&target.Specialization,
		&target.GroupName,
	)
}

func scanStudentSupervisorPair(row pgx.Row, target *models.StudentSupervisorPair) error {
	return row.Scan(
		&target.Student.StudentID,
		&target.Student.FullName,
		//&target.Student.Department,
		&target.Student.ActualSemester,
		&target.Student.Years,
		&target.Student.StartDate,
		&target.Student.StudyingStatus,
		&target.Student.Status,
		&target.Student.CanEdit,
		&target.Student.Phone,
		&target.Student.Category,
		&target.Student.EndDate,
		&target.Student.Specialization,
		&target.Student.GroupName,
		&target.Supervisor.SupervisorID,
		&target.Supervisor.FullName,
	)
}

func scanSupervisor(row pgx.Row, target *models.Supervisor) error {
	return row.Scan(
		&target.SupervisorID,
		&target.FullName,
		&target.Phone,
		&target.Archived,
		&target.Faculty,
		&target.Department,
		&target.Degree,
		&target.Rank,
		&target.Position,
	)
}

func scanSupervisorFull(row pgx.Row, target *models.SupervisorFull) error {
	return row.Scan(
		&target.SupervisorID,
		&target.FullName,
		&target.StartAt,
		&target.EndAt,
	)
}
