package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
	"ds-backend/internal/infrastructure/db"
)

type TokenRepository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

func NewTokenRepository(db *sql.DB) ports.TokenRepository {
	return &TokenRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

func (r *TokenRepository) SaveToken(ctx context.Context, token *models.RefreshToken) error {
	query := r.qb.Insert("refresh_tokens").
		Columns("id", "user_id", "token", "expires_at", "created_at").
		Values(token.ID, token.UserID, token.Token, token.ExpiresAt, time.Now())

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to save refresh token")
		return err
	}

	logrus.Info("New refresh token saved successfully")
	return nil
}

func (r *TokenRepository) DeleteToken(ctx context.Context, tokenID string) error {
	query := r.qb.Delete("refresh_tokens").
		Where(squirrel.Eq{"id": tokenID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		if db.GetDBErrorNoRows(err) {
			logrus.Warn("No refresh token found with the provided ID")
			return nil
		}
		logrus.WithError(err).Error("Failed to delete refresh token")
		return err
	}

	logrus.Info("Refresh token deleted successfully")
	return nil
}

func (r *TokenRepository) GetTokenByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
    query := r.qb.Select("id", "user_id", "token", "expires_at", "created_at").
        From("refresh_tokens").
        Where(squirrel.Eq{"token": token})

    sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var tokenResponse models.RefreshToken
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(&tokenResponse.ID, &tokenResponse.UserID, &tokenResponse.Token, &tokenResponse.ExpiresAt, &tokenResponse.CreatedAt)
	if err != nil {
		if db.GetDBErrorNoRows(err) {
			logrus.Warn("No refresh token found with the provided token")
			return nil, errors.New("no token found")
		}
		return nil, err
	}

	return &tokenResponse, nil
}

func (r *TokenRepository) GetTokenByUserID(ctx context.Context, userID uint) (*models.RefreshToken, error) {
	query := r.qb.Select("id", "user_id", "token", "expires_at", "created_at").
		From("refresh_tokens").
		Where(squirrel.Eq{"user_id": userID})

	sql, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query")
		return nil, err
	}

	var tokenResponse models.RefreshToken
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(&tokenResponse.ID, &tokenResponse.UserID, &tokenResponse.Token, &tokenResponse.ExpiresAt, &tokenResponse.CreatedAt)
	if err != nil {
		if db.GetDBErrorNoRows(err) {
			logrus.Warn("No refresh token found for the provided user ID")
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to execute query")
		return nil, err
	}

	return &tokenResponse, nil
}

func (r *TokenRepository) DeleteExpiredTokensByUserID(ctx context.Context, userID uint) error {
	query := r.qb.Delete("refresh_tokens").
        Where(squirrel.Eq{"user_id": userID}).
        Where(squirrel.Lt{"expires_at": time.Now().Unix()})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete expired refresh tokens")
		return err
	}

	logrus.Infof("Expired refresh tokens deleted successfully for user ID: %d", userID)
	return nil
}
