package repository

import (
    "context"
    "database/sql"
    "errors"
    "time"

    "github.com/Masterminds/squirrel"
    "github.com/sirupsen/logrus"

    "ds-backend/internal/infrastructure/db"
    "ds-backend/internal/domain/ports"
)

type EmailRepository struct {
    db *sql.DB
    qb squirrel.StatementBuilderType
}

func NewEmailRepository(db *sql.DB) ports.EmailRepository {
    return &EmailRepository{
        db: db,
        qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
    }
}

// SaveVerificationCode inserts a new verification token into the database
func (r *EmailRepository) SaveVerificationCode(ctx context.Context, toEmail string, code string, motivo string) error {
    expirationDate := time.Now().Add(48 * time.Hour)

    query := r.qb.Insert("verification_tokens").
        Columns("token", "user_id", "expires_at", "created_at", "used", "motive", "resends").
        Values(code, squirrel.Expr("(SELECT id FROM users WHERE email = ?)", toEmail), expirationDate, time.Now(), false, motivo, 0)

    sql, args, err := query.ToSql()
    if err != nil {
        logrus.WithError(err).Error("Failed to build SQL query for SaveVerificationCode")
        return err
    }

    _, err = r.db.ExecContext(ctx, sql, args...)
    if err != nil {
        logrus.WithError(err).Error("Error inserting verification code in database")
        return err
    }

    logrus.Info("Verification code saved successfully")
    return nil
}

// GetIDFromEmail retrieves the user ID by email
func (r *EmailRepository) GetIDFromEmail(ctx context.Context, toEmail string) (int, error) {
    query := r.qb.Select("id").
        From("users").
        Where(squirrel.Eq{"email": toEmail})

    sql, args, err := query.ToSql()
    if err != nil {
        logrus.WithError(err).Error("Failed to build SQL query for GetIDFromEmail")
        return 0, err
    }

    var userID int

    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&userID)
    if err != nil {
        if db.GetDBErrorNoRows(err) {
            logrus.Warn("User not found with email:", toEmail)
            return 0, errors.New("user not found")
        }
        logrus.WithError(err).Error("Error fetching user ID")
        return 0, err
    }

    return userID, nil
}

// VerifyTokenAndUser checks if a verification token is valid and not used
func (r *EmailRepository) VerifyTokenAndUser(ctx context.Context, code string, toEmail string) (int, bool, error) {
    query := r.qb.Select("t.user_id", "t.used").
        From("verification_tokens t").
        Join("users u ON t.user_id = u.id").
        Where(squirrel.Eq{"t.token": code, "u.email": toEmail}).
        Where(squirrel.Gt{"t.expires_at": time.Now()})

    sql, args, err := query.ToSql()
    if err != nil {
        logrus.WithError(err).Error("Failed to build SQL query for VerifyTokenAndUser")
        return 0, false, err
    }

    var userID int
    var used bool

    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&userID, &used)
    if err != nil {
        if db.GetDBErrorNoRows(err) {
            logrus.Warn("No matching verification code or email found")
            return 0, false, errors.New("invalid verification code or email")
        }
        logrus.WithError(err).Error("Error fetching user by verification code")
        return 0, false, err
    }

    return userID, used, nil
}

// UpdateTokenAsUsed marks a token as used
func (r *EmailRepository) UpdateTokenAsUsed(ctx context.Context, code string, userID int) error {
    query := r.qb.Update("verification_tokens").
        Set("used", true).
        Set("used_at", time.Now()).
        Where(squirrel.Eq{"token": code, "user_id": userID})

    sql, args, err := query.ToSql()
    if err != nil {
        logrus.WithError(err).Error("Failed to build SQL query for UpdateTokenAsUsed")
        return err
    }


    _, err = r.db.ExecContext(ctx, sql, args...)
    if err != nil {
        logrus.WithError(err).Error("Error updating token status")
        return err
    }

    logrus.Info("Token marked as used successfully")
    return nil
}

// GetTokenVerificationStatus retrieves the verification status and motive of a token
func (r *EmailRepository) GetTokenVerificationStatus(ctx context.Context, userID int) (bool, string, error) {
    query := r.qb.Select("used", "motive").
        From("verification_tokens").
        Where(squirrel.Eq{"user_id": userID}).
        OrderBy("created_at DESC").
        Limit(1)

    sql, args, err := query.ToSql()
    if err != nil {
        logrus.WithError(err).Error("Failed to build SQL query for GetTokenVerificationStatus")
        return false, "", err
    }

    var used bool
    var motive string

    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&used, &motive)
    if err != nil {
        if db.GetDBErrorNoRows(err) {
            logrus.Warn("No verification token found for user ID:", userID)
            return false, "", nil
        }
        logrus.WithError(err).Error("Error fetching user verification status")
        return false, "", err
    }

    return used, motive, nil
}

// GetLatestTokenInfo retrieves the resend count and token of the latest verification token
func (r *EmailRepository) GetLatestTokenInfo(ctx context.Context, userID int) (int, string, error) {
    query := r.qb.Select("resends", "token").
        From("verification_tokens").
        Where(squirrel.Eq{"user_id": userID}).
        OrderBy("created_at DESC").
        Limit(1)

    sql, args, err := query.ToSql()
    if err != nil {
        logrus.WithError(err).Error("Failed to build SQL query for GetLatestTokenInfo")
        return 0, "", err
    }

    var resends int
    var token string

    err = r.db.QueryRowContext(ctx, sql, args...).Scan(&resends, &token)
    if err != nil {
        if db.GetDBErrorNoRows(err) {
            logrus.Warn("No verification token found for user ID:", userID)
            return 0, "", nil
        }
        logrus.WithError(err).Error("Error fetching resend count")
        return 0, "", err
    }

    return resends, token, nil
}

// UpdateTokenResendInfo updates the resend count and token for a verification token
func (r *EmailRepository) UpdateTokenResendInfo(ctx context.Context, userID int, oldToken string, newToken string) error {
    expirationDate := time.Now().Add(48 * time.Hour)

    query := r.qb.Update("verification_tokens").
        Set("resends", squirrel.Expr("resends + 1")).
        Set("token", newToken).
        Set("updated_at", time.Now()).
        Set("expires_at", expirationDate).
        Where(squirrel.Eq{"user_id": userID, "token": oldToken}).
        OrderBy("created_at DESC").
        Limit(1)

    sql, args, err := query.ToSql()
    if err != nil {
        logrus.WithError(err).Error("Failed to build SQL query for UpdateTokenResendInfo")
        return err
    }


    _, err = r.db.ExecContext(ctx, sql, args...)
    if err != nil {
        logrus.WithError(err).Error("Error updating resend count")
        return err
    }

    logrus.Info("Resend count updated successfully")
    return nil
}

// UpdateUserVerificationStatus marks a user as verified
func (r *EmailRepository) UpdateUserVerificationStatus(ctx context.Context, userID int) error {
    query := r.qb.Update("users").
        Set("verified", true).
        Where(squirrel.Eq{"id": userID})

    sql, args, err := query.ToSql()
    if err != nil {
        logrus.WithError(err).Error("Failed to build SQL query for UpdateUserVerificationStatus")
        return err
    }


    _, err = r.db.ExecContext(ctx, sql, args...)
    if err != nil {
        logrus.WithError(err).Error("Error updating user verification status")
        return err
    }

    logrus.Info("User verification status updated successfully")
    return nil
}