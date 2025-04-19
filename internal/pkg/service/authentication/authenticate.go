package authentication

import (
	"context"
	"strings"
	"time"

	"uir_draft/internal/generated/new_kasper/new_uir/public/model"
	"uir_draft/internal/handlers/authorization_handler/request_models"
	"uir_draft/internal/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) AuthenticateWithUserType(ctx context.Context, token, userType string) (*model.Users, error) {
	var user model.Users

	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		tokenModel, err := s.tokenRepo.GetByTokenNumberTx(ctx, tx, token)
		if err != nil {
			return errors.Wrap(err, "getting user_id by token")
		}

		if !tokenModel.IsActive {
			return models.ErrTokenExpired
		}

		user, err = s.userRepo.GetUserTx(ctx, tx, tokenModel.UserID)
		if err != nil {
			return errors.Wrap(err, "getting user info")
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "AuthenticateWithUserType()")
	}

	// Проверка типа пользователя
	if user.UserType.String() != userType {
		// Разрешаем администратору доступ ко всем типам
		if user.UserType == model.UserType_Admin {
			return &user, nil
		}
		// Дополнительно проверяем Supervisor
		if user.UserType == model.UserType_Admin && userType == model.UserType_Supervisor.String() {
			return &user, nil
		}
		return nil, models.ErrWrongUserType
	}

	return &user, nil
}

func (s *Service) Authenticate(ctx context.Context, token string) (*model.Users, error) {
	var user model.Users

	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		tokenModel, err := s.tokenRepo.GetByTokenNumberTx(ctx, tx, token)
		if err != nil {
			return errors.Wrap(err, "getting user_id by token")
		}

		if !tokenModel.IsActive {
			return models.ErrTokenExpired
		}

		user, err = s.userRepo.GetUserTx(ctx, tx, tokenModel.UserID)
		if err != nil {
			return errors.Wrap(err, "getting user info")
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "Authenticate()")
	}

	return &user, nil
}

func (s *Service) TokenCheck(ctx context.Context, token string) (*model.Users, error) {
	var user model.Users

	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		tokenModel, err := s.tokenRepo.GetByTokenNumberTx(ctx, tx, token)
		if err != nil {
			return errors.Wrap(err, "getting user_id by token")
		}

		if !tokenModel.IsActive {
			return models.ErrTokenExpired
		}

		user, err = s.userRepo.GetUserTx(ctx, tx, tokenModel.UserID)
		if err != nil {
			return errors.Wrap(err, "getting user info")
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "TokenCheck()")
	}

	return &user, nil
}

func (s *Service) Authorize(ctx context.Context, request models.AuthorizeRequest) (*models.AuthorizeResponse, bool, error) {
	var response models.AuthorizeResponse
	var isAuthorized bool

	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		user, err := s.userRepo.GetUserByEmailTx(ctx, tx, strings.TrimSpace(request.Email))
		if err != nil {
			return errors.Wrap(err, "getting user info")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
		if err != nil {
			isAuthorized = false
			return nil
		}

		isAuthorized = true

		token := model.AuthorizationToken{
			UserID:      user.UserID,
			IsActive:    true,
			TokenNumber: uuid.New().String(),
			CreatedAt:   time.Now(),
			UpdatedAt:   nil,
		}

		if err = s.tokenRepo.InsertTokenTx(ctx, tx, &token); err != nil {
			return err
		}

		response = models.AuthorizeResponse{
			UserType:   user.UserType.String(),
			Token:      token.TokenNumber,
			Registered: user.Registered,
		}

		return nil
	}); err != nil {
		return nil, false, errors.Wrap(err, "AuthenticateWithUserType()")
	}

	if !isAuthorized {
		return nil, false, nil
	}

	return &response, isAuthorized, nil
}

func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, request request_models.ChangePasswordRequest) error {
	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		user, err := s.userRepo.GetUserTx(ctx, tx, userID)
		if err != nil {
			return errors.Wrap(err, "getting user info")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
		if err != nil {
			return models.ErrWrongPassword
		}

		newPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		err = s.userRepo.ChangeUsersPasswordTx(ctx, tx, userID, string(newPassword))

		return err
	}); err != nil {
		return errors.Wrap(err, "ChangePassword()")
	}

	return nil
}

func (s *Service) GetUserProfile(ctx context.Context, userID uuid.UUID) (model.Users, error) {
	var user model.Users

	if err := s.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		dUser, err := s.userRepo.GetUserTx(ctx, tx, userID)
		if err != nil {
			return errors.Wrap(err, "getting user info")
		}
		user = dUser

		return nil
	}); err != nil {
		return model.Users{}, errors.Wrap(err, "GetUserProfile()")
	}

	return user, nil
}

func (s *Service) CreateAnonymousToken(ctx context.Context, token string) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO authorization_token (token_number, is_active)
        VALUES ($1, false)
    `, token)
	return err
}

func (s *Service) TokenExists(ctx context.Context, token string) (bool, error) {
	row := s.db.QueryRow(ctx,
		`SELECT COUNT(*) > 0 FROM authorization_token WHERE token_number = $1`, token)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Service) CreateUser(ctx context.Context, user *model.Users) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO users (user_id, email, password, kasper_id, user_type, registered)
       VALUES ($1,$2,$3,$4,$5,$6)`,
		user.UserID, user.Email, user.Password, user.KasperID, user.UserType.String(), user.Registered)
	return err
}

func (s *Service) AttachTokenToUser(ctx context.Context, token string, userID uuid.UUID) error {
	_, err := s.db.Exec(ctx,
		`UPDATE authorization_token
         SET user_id = $1, is_active = true
       WHERE token_number = $2`,
		userID, token)
	return err
}
