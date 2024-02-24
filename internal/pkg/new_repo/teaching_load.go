package new_repo

import (
	"context"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/generated/new_kasper/new_uir/public/table"
	"uir_draft/internal/pkg/domain"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type TeachingLoadStatusRepository struct{}

func NewTeachingLoadStatusRepository() *TeachingLoadStatusRepository {
	return &TeachingLoadStatusRepository{}
}

func (r *TeachingLoadStatusRepository) GetTeachingLoadStatusTx(ctx context.Context, tx *pgxpool.Tx, studentID uuid.UUID) ([]model.TeachingLoadStatus, error) {
	stmt, args := table.TeachingLoadStatus.
		SELECT(table.TeachingLoadStatus.AllColumns).
		WHERE(table.TeachingLoadStatus.StudentID.EQ(postgres.UUID(studentID))).
		Sql()

	rows, err := tx.Query(ctx, stmt, args...)
	if err != nil {
		return nil, errors.Wrap(err, "GetTeachingLoadStatusTx()")
	}

	loads := make([]model.TeachingLoadStatus, 0, 8)

	for rows.Next() {
		load := model.TeachingLoadStatus{}

		if err := scanTeachingLoadStatusStatus(rows, &load); err != nil {
			return nil, errors.Wrap(err, "GetTeachingLoadStatusTx()")
		}

		loads = append(loads, load)
	}

	return loads, nil
}

func (r *TeachingLoadStatusRepository) UpdateTeachingLoadStatusTx(ctx context.Context, tx *pgxpool.Tx, loads []model.TeachingLoadStatus) error {
	for _, load := range loads {
		stmt, args := table.TeachingLoadStatus.
			UPDATE(
				table.TeachingLoadStatus.Status,
				table.TeachingLoadStatus.UpdatedAt,
				table.TeachingLoadStatus.AcceptedAt,
			).
			SET(
				load.Status,
				load.UpdatedAt,
				load.AcceptedAt,
			).
			WHERE(table.TeachingLoadStatus.LoadsID.EQ(postgres.UUID(load.LoadsID)).
				AND(table.TeachingLoadStatus.Semester.EQ(postgres.Int32(load.Semester)))).
			Sql()

		if _, err := tx.Exec(ctx, stmt, args...); err != nil {
			return errors.Wrap(err, "UpdateTeachingLoadStatusTx()")
		}
	}

	return nil
}

func (r *TeachingLoadStatusRepository) InsertClassroomLoadsTx(ctx context.Context, tx *pgxpool.Pool, loads []model.ClassroomLoad) error {
	stmt, args := table.ClassroomLoad.
		INSERT().
		MODELS(loads).
		Sql()

	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		return errors.Wrap(err, "InsertClassroomLoadTx()")
	}

	return nil
}

func (r *TeachingLoadStatusRepository) UpdateClassroomLoadsTx(ctx context.Context, tx *pgxpool.Tx, loads []model.ClassroomLoad) error {
	for _, load := range loads {
		stmt, args := table.ClassroomLoad.
			UPDATE(
				table.ClassroomLoad.Hours,
				table.ClassroomLoad.LoadType,
				table.ClassroomLoad.MainTeacher,
				table.ClassroomLoad.GroupName,
				table.ClassroomLoad.SubjectName,
			).
			SET(
				load.Hours,
				load.LoadType,
				load.MainTeacher,
				load.GroupName,
				load.SubjectName,
			).
			WHERE(table.ClassroomLoad.LoadID.EQ(postgres.UUID(load.LoadID))).
			Sql()

		if _, err := tx.Exec(ctx, stmt, args...); err != nil {
			return errors.Wrap(err, "UpdateClassroomLoadTx()")
		}
	}

	return nil
}

func (r *ScientificRepository) DeleteClassroomLoadsTx(ctx context.Context, tx *pgxpool.Tx, classroomsIDs []uuid.UUID) error {
	var exps []postgres.Expression
	for _, id := range classroomsIDs {
		exp := postgres.Expression(postgres.UUID(id))

		exps = append(exps, exp)
	}

	stmt, args := table.ClassroomLoad.
		DELETE().
		WHERE(table.ClassroomLoad.LoadID.IN(exps...)).
		Sql()

	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		return errors.Wrap(err, "DeleteClassroomLoadTx()")
	}

	return nil
}

func (r *TeachingLoadStatusRepository) InsertIndividualLoadsTx(ctx context.Context, tx *pgxpool.Tx, loads []model.IndividualStudentsLoad) error {
	stmt, args := table.IndividualStudentsLoad.
		INSERT().
		MODELS(loads).
		Sql()

	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		return errors.Wrap(err, "InsertIndividualLoadTx()")
	}

	return nil
}

func (r *TeachingLoadStatusRepository) UpdateIndividualLoadsTx(ctx context.Context, tx *pgxpool.Tx, loads []model.IndividualStudentsLoad) error {
	for _, load := range loads {
		stmt, args := table.IndividualStudentsLoad.
			UPDATE(
				table.IndividualStudentsLoad.StudentsAmount,
				table.IndividualStudentsLoad.Comment,
			).
			SET(
				load.StudentsAmount,
				load.Comment,
			).
			WHERE(table.IndividualStudentsLoad.LoadID.EQ(postgres.UUID(load.LoadID))).
			Sql()

		if _, err := tx.Exec(ctx, stmt, args...); err != nil {
			return errors.Wrap(err, "UpdateIndividualLoadTx()")
		}
	}

	return nil
}

func (r *ScientificRepository) DeleteIndividualStudentsLoadsTx(ctx context.Context, tx *pgxpool.Tx, individualsIDs []uuid.UUID) error {
	var exps []postgres.Expression
	for _, id := range individualsIDs {
		exp := postgres.Expression(postgres.UUID(id))

		exps = append(exps, exp)
	}

	stmt, args := table.IndividualStudentsLoad.
		DELETE().
		WHERE(table.IndividualStudentsLoad.LoadID.IN(exps...)).
		Sql()

	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		return errors.Wrap(err, "DeleteIndividualStudentsLoadTx()")
	}

	return nil
}

func (r *TeachingLoadStatusRepository) InsertAdditionalLoadsTx(ctx context.Context, tx *pgxpool.Tx, loads []model.AdditionalLoad) error {
	stmt, args := table.AdditionalLoad.
		INSERT().
		MODELS(loads).
		Sql()

	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		return errors.Wrap(err, "InsertAdditionalLoadTx()")
	}

	return nil
}

func (r *TeachingLoadStatusRepository) UpdateAdditionalLoadsTx(ctx context.Context, tx *pgxpool.Tx, loads []model.AdditionalLoad) error {
	for _, load := range loads {
		stmt, args := table.AdditionalLoad.
			UPDATE(
				table.AdditionalLoad.Name,
				table.AdditionalLoad.Volume,
				table.AdditionalLoad.Comment,
			).
			SET(
				load.Name,
				load.Volume,
				load.Comment,
			).
			WHERE(table.AdditionalLoad.LoadID.EQ(postgres.UUID(load.LoadID))).
			Sql()

		if _, err := tx.Exec(ctx, stmt, args...); err != nil {
			return errors.Wrap(err, "UpdateAdditionalLoadTx()")
		}
	}

	return nil
}

func (r *ScientificRepository) DeleteAdditionalLoadsTx(ctx context.Context, tx *pgxpool.Tx, additionalIDs []uuid.UUID) error {
	var exps []postgres.Expression
	for _, id := range additionalIDs {
		exp := postgres.Expression(postgres.UUID(id))

		exps = append(exps, exp)
	}

	stmt, args := table.AdditionalLoad.
		DELETE().
		WHERE(table.AdditionalLoad.LoadID.IN(exps...)).
		Sql()

	if _, err := tx.Exec(ctx, stmt, args...); err != nil {
		return errors.Wrap(err, "DeleteAdditionalLoadTx()")
	}

	return nil
}

func (r *TeachingLoadStatusRepository) GetTeachingLoadsTx(ctx context.Context, tx *pgxpool.Tx, studentID uuid.UUID) ([]domain.TeachingLoad, error) {
	stmt, args := table.TeachingLoadStatus.
		SELECT(
			table.TeachingLoadStatus.LoadsID,
			table.TeachingLoadStatus.StudentID,
			table.TeachingLoadStatus.Semester,
			table.TeachingLoadStatus.Status.AS("teaching_load.approval_status"),
			table.TeachingLoadStatus.UpdatedAt,
			table.TeachingLoadStatus.AcceptedAt,
			table.ClassroomLoad.AllColumns.Except(table.ClassroomLoad.TLoadID),
			table.IndividualStudentsLoad.AllColumns.Except(table.IndividualStudentsLoad.TLoadID),
			table.AdditionalLoad.AllColumns.Except(table.AdditionalLoad.TLoadID),
		).
		FROM(table.TeachingLoadStatus.
			INNER_JOIN(table.ClassroomLoad, table.TeachingLoadStatus.LoadsID.EQ(table.ClassroomLoad.TLoadID)).
			INNER_JOIN(table.IndividualStudentsLoad, table.TeachingLoadStatus.LoadsID.EQ(table.IndividualStudentsLoad.TLoadID)).
			INNER_JOIN(table.AdditionalLoad, table.TeachingLoadStatus.LoadsID.EQ(table.AdditionalLoad.TLoadID)),
		).
		WHERE(table.TeachingLoadStatus.StudentID.EQ(postgres.UUID(studentID))).
		Sql()

	rows, err := tx.Query(ctx, stmt, args...)
	if err != nil {
		return nil, errors.Wrap(err, "GetTeachingLoadTx()")
	}

	loads := make([]domain.TeachingLoad, 0, 10)

	for rows.Next() {
		load := domain.TeachingLoad{}

		if err := scanTeachingLoadStatus(rows, &load); err != nil {
			return nil, errors.Wrap(err, "GetTeachingLoadTx(): scanning row")
		}

		loads = append(loads, load)
	}

	return loads, nil
}

func scanTeachingLoadStatusStatus(row pgx.Row, target *model.TeachingLoadStatus) error {
	return row.Scan(
		&target.LoadsID,
		&target.StudentID,
		&target.Semester,
		&target.Status,
		&target.UpdatedAt,
		&target.AcceptedAt,
	)
}

func scanTeachingLoadStatus(row pgx.Row, target *domain.TeachingLoad) error {
	return row.Scan(
		&target.LoadsID,
		&target.StudentID,
		&target.Semester,
		&target.ApprovalStatus,
		&target.UpdatedAt,
		&target.AcceptedAt,
		&target.ClassroomLoad.LoadID,
		&target.ClassroomLoad.Hours,
		&target.ClassroomLoad.LoadType,
		&target.ClassroomLoad.MainTeacher,
		&target.ClassroomLoad.GroupName,
		&target.ClassroomLoad.SubjectName,
		&target.IndividualStudentsLoad.LoadID,
		&target.IndividualStudentsLoad.StudentsAmount,
		&target.IndividualStudentsLoad.Comment,
		&target.AdditionalLoad.LoadID,
		&target.AdditionalLoad.Name,
		&target.AdditionalLoad.Volume,
		&target.AdditionalLoad.Comment,
	)
}
