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

type ImageRepository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

func NewImageRepository(db *sql.DB) ports.ImageRepository {
	return &ImageRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

func (r *ImageRepository) SaveImage(ctx context.Context, image *models.Image) (*models.Image, error) {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    if image.MainImage {
        unsetQuery := r.qb.Update("images").
            Set("main_image", false).
            Where(squirrel.Eq{"property_id": image.PropertyID}).
            Where(squirrel.Eq{"main_image": true}).
            Where(squirrel.Expr("deleted_at IS NULL"))

        sql, args, err := unsetQuery.ToSql()
        if err != nil {
            return nil, err
        }
        _, err = tx.ExecContext(ctx, sql, args...)
        if err != nil {
            return nil, err
        }
    }

    insertQuery := r.qb.Insert("images").
        Columns("property_id", "path", "description", "main_image").
        Values(image.PropertyID, image.Path, image.Description, image.MainImage)

    sql, args, err := insertQuery.ToSql()
    if err != nil {
        return nil, err
    }

    result, err := tx.ExecContext(ctx, sql, args...)
    if err != nil {
        return nil, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }
    image.ID = uint(id)

    if err = tx.Commit(); err != nil {
        return nil, err
    }

    return image, nil
}

func (r *ImageRepository) GetImageByID(ctx context.Context, id uint) (*models.Image, error) {
	query := r.qb.Select("id", "property_id", "path", "description", "main_image", "created_at", "updated_at").
		From("images").
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.Expr("deleted_at IS NULL"))
	
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var image models.Image
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&image.ID,
		&image.PropertyID,
		&image.Path,
		&image.Description,
		&image.MainImage,
		&image.CreatedAt,
		&image.UpdatedAt,
	)
	if err != nil {
		if db.GetDBErrorNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &image, nil
}

func (r *ImageRepository) GetImagesByPropertyID(ctx context.Context, propertyID uint) ([]models.Image, error) {
	query := r.qb.Select("id", "property_id", "path", "description", "main_image", "created_at", "updated_at").
		From("images").
		Where(squirrel.Eq{"property_id": propertyID}).
		Where(squirrel.Expr("deleted_at IS NULL"))
	
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var image models.Image
		err := rows.Scan(
			&image.ID,
			&image.PropertyID,
			&image.Path,
			&image.Description,
			&image.MainImage,
			&image.CreatedAt,
			&image.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ImageRepository) GetMainImageByPropertyID(ctx context.Context, propertyID uint) (*models.Image, error) {
	query := r.qb.Select("id", "property_id", "path", "description", "main_image", "created_at", "updated_at").
		From("images").
		Where(squirrel.Eq{"property_id": propertyID}).
		Where(squirrel.Eq{"main_image": true}).
		Where(squirrel.Expr("deleted_at IS NULL"))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var image models.Image
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&image.ID,
		&image.PropertyID,
		&image.Path,
		&image.Description,
		&image.MainImage,
		&image.CreatedAt,
		&image.UpdatedAt,
	)
	if err != nil {
		if db.GetDBErrorNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &image, nil
}

func (r *ImageRepository) UpdateMainImageStatus(ctx context.Context, propertyID uint, imageID uint) error {
    image, err := r.GetImageByID(ctx, imageID)
    if err != nil {
		logrus.Errorf("Failed to get image by ID: %v", err)
        return err
    }
    if image == nil || image.PropertyID != propertyID {
		logrus.Errorf("Image does not belong to property: imageID=%d, propertyID=%d", imageID, propertyID)
        return errors.New("image does not belong to this property")
    }
	
	query := r.qb.Update("images").
        Set("main_image", false).
        Where(squirrel.Eq{"property_id": propertyID}).
        Where(squirrel.Eq{"main_image": true}).
        Where(squirrel.Expr("deleted_at IS NULL"))

    sql, args, err := query.ToSql()
    if err != nil {
        return err
    }

    _, err = r.db.ExecContext(ctx, sql, args...)
    if err != nil {
        return err
    }

    query = r.qb.Update("images").
        Set("main_image", true).
        Where(squirrel.Eq{"id": imageID}).
        Where(squirrel.Eq{"property_id": propertyID}).
        Where(squirrel.Expr("deleted_at IS NULL"))

    sql, args, err = query.ToSql()
    if err != nil {
        return err
    }

    _, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		logrus.Errorf("Failed to update main image status: %v", err)
		return err
	}
    return err
}

func (r *ImageRepository) DeleteImagesByPropertyID(ctx context.Context, propertyID uint) error {
	query := r.qb.Update("images").
		Set("deleted_at", time.Now()).
		Where(squirrel.Eq{"property_id": propertyID}).
		Where(squirrel.Expr("deleted_at IS NULL"))
	
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *ImageRepository) UpdateImage(ctx context.Context, image *models.Image) (*models.Image, error) {
	query := r.qb.Update("images").
		Set("path", image.Path).
		Set("description", image.Description).
		Where(squirrel.Eq{"id": image.ID}).
		Where(squirrel.Expr("deleted_at IS NULL"))
	
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (r *ImageRepository) DeleteImage(ctx context.Context, id uint) error {
	query := r.qb.Update("images").
		Set("deleted_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.Expr("deleted_at IS NULL"))
		
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *ImageRepository) HardDeleteImage(ctx context.Context, id uint) error {
    query := r.qb.Delete("images").Where(squirrel.Eq{"id": id})
    sql, args, err := query.ToSql()
    if err != nil {
        return err
    }
    _, err = r.db.ExecContext(ctx, sql, args...)
    return err
}