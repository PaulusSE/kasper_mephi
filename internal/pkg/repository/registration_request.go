package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type RegistrationRequest struct {
	RequestID    uuid.UUID
	Email        string
	PasswordHash string
	UserType     string
	Payload      []byte
	Status       string
	CreatedAt    time.Time
}

type RegistrationRequestsRepository interface {
	InsertRequestTx(ctx context.Context, tx pgx.Tx, r *RegistrationRequest) error
	ListRequestsTx(ctx context.Context, tx pgx.Tx, userType string) ([]RegistrationRequest, error)
	UpdateStatusTx(ctx context.Context, tx pgx.Tx, requestID uuid.UUID, status string) error
	GetByIDTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (RegistrationRequest, error)
}

func NewRegistrationRequestsRepository() RegistrationRequestsRepository {
	return &registrationRequestsRepo{}
}

type registrationRequestsRepo struct{}

func (r *registrationRequestsRepo) InsertRequestTx(ctx context.Context, tx pgx.Tx, req *RegistrationRequest) error {
	_, err := tx.Exec(ctx, `
        INSERT INTO registration_requests
          (request_id, email, password_hash, user_type, payload, status, created_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7)
    `,
		req.RequestID, req.Email, req.PasswordHash,
		req.UserType, req.Payload, req.Status, req.CreatedAt,
	)
	return err
}

//func (r *registrationRequestsRepo) ListRequestsTx(ctx context.Context, tx pgx.Tx, userType string) ([]RegistrationRequest, error) {
//	var rows pgx.Rows
//	var err error
//	if userType == "admin" {
//		rows, err = tx.Query(ctx, `SELECT * FROM registration_requests WHERE status='pending' ORDER BY created_at`)
//	} else {
//		// TODO: пока руководитель видит заявки других руководителей
//		rows, err = tx.Query(ctx, `SELECT * FROM registration_requests WHERE status='pending' ORDER BY created_at`)
//	}
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var res []RegistrationRequest
//	for rows.Next() {
//		var r RegistrationRequest
//		if err := rows.Scan(
//			&r.RequestID, &r.Email, &r.PasswordHash, &r.UserType,
//			&r.Payload, &r.Status, &r.CreatedAt,
//		); err != nil {
//			return nil, err
//		}
//		res = append(res, r)
//	}
//	return res, nil
//}

const listSQL = `
SELECT
    request_id,
    email,
    password_hash,
    user_type::text,         -- кастуем enum к text
    payload,
    status::text,
    created_at
FROM registration_requests
WHERE status = 'pending'
ORDER BY created_at;
`

func (r *registrationRequestsRepo) ListRequestsTx(ctx context.Context, tx pgx.Tx, userType string) ([]RegistrationRequest, error) {
	rows, err := tx.Query(ctx, listSQL)
	if err != nil {
		return nil, err
	}

	var res []RegistrationRequest
	for rows.Next() {
		var rr RegistrationRequest
		if err = rows.Scan(
			&rr.RequestID,
			&rr.Email,
			&rr.PasswordHash,
			&rr.UserType,
			&rr.Payload,
			&rr.Status,
			&rr.CreatedAt,
		); err != nil {
			return nil, err // здесь раньше падало
		}
		res = append(res, rr)
	}
	return res, rows.Err()
}

// Апдейт статуса:
func (r *registrationRequestsRepo) UpdateStatusTx(ctx context.Context, tx pgx.Tx, requestID uuid.UUID, status string) error {
	_, err := tx.Exec(ctx,
		`UPDATE registration_requests SET status=$1 WHERE request_id=$2`,
		status, requestID,
	)
	return err
}

func (r *registrationRequestsRepo) GetByIDTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (RegistrationRequest, error) {
	var rr RegistrationRequest
	row := tx.QueryRow(ctx, `
      SELECT request_id, email, password_hash, user_type, payload, status, created_at
      FROM registration_requests
      WHERE request_id=$1`, id)
	err := row.Scan(
		&rr.RequestID,
		&rr.Email,
		&rr.PasswordHash,
		&rr.UserType,
		&rr.Payload,
		&rr.Status,
		&rr.CreatedAt,
	)
	return rr, err
}
