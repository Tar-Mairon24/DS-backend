package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
)

type ImageUseCase struct {
	imageRepo ports.ImageRepository
}

func NewImageUseCase(imageRepo ports.ImageRepository) ports.ImageUseCase {
	return &ImageUseCase{
		imageRepo: imageRepo,
	}
}

func (i *ImageUseCase) SaveImage(ctx context.Context, image *models.Image) (*models.Image, error) {
    if image == nil {
        return nil, errors.New("image cannot be nil")
    }
    return i.imageRepo.SaveImage(ctx, image)
}

func (i *ImageUseCase) GeneratePath(fileName string, propertyID uint) (diskPath string, urlPath string, err error) {
    ext := filepath.Ext(fileName)
    allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
    if !allowedExts[ext] {
        return "", "", errors.New("unsupported file type")
    }

    propertyDir := filepath.Join("/app/uploads/properties", fmt.Sprintf("%d", propertyID))
    if err := os.MkdirAll(propertyDir, os.ModePerm); err != nil {
        return "", "", errors.New("could not create upload directory")
    }

    newFileName := uuid.New().String() + ext
    diskPath = filepath.Join(propertyDir, newFileName)
    urlPath = fmt.Sprintf("/uploads/properties/%d/%s", propertyID, newFileName)
    return diskPath, urlPath, nil
}

func (i *ImageUseCase) GetImageByID(ctx context.Context, id uint) (*models.Image, error) {
	return i.imageRepo.GetImageByID(ctx, id)
}

func (i *ImageUseCase) GetImagesByPropertyID(ctx context.Context, propertyID uint) ([]models.Image, error) {
	return i.imageRepo.GetImagesByPropertyID(ctx, propertyID)
}

func (i *ImageUseCase) GetMainImageByPropertyID(ctx context.Context, propertyID uint) (*models.Image, error) {
	return i.imageRepo.GetMainImageByPropertyID(ctx, propertyID)
}

func (i *ImageUseCase) UpdateMainImageStatus(ctx context.Context, propertyID uint, imageID uint) error {
	return i.imageRepo.UpdateMainImageStatus(ctx, propertyID, imageID)
}

func (i *ImageUseCase) UpdateImage(ctx context.Context, image *models.Image) (*models.Image, error) {
	return i.imageRepo.UpdateImage(ctx, image)
}

func (i *ImageUseCase) DeleteImage(ctx context.Context, id uint) error {
    image, err := i.imageRepo.GetImageByID(ctx, id)
    if err != nil {
        return err
    }
    if image == nil {
        return errors.New("image not found")
    }

    if err := i.imageRepo.DeleteImage(ctx, id); err != nil {
        return err
    }

    if image.MainImage {
        images, err := i.imageRepo.GetImagesByPropertyID(ctx, image.PropertyID)
        if err != nil {
            return err
        }
        if len(images) > 0 {
            if err := i.imageRepo.UpdateMainImageStatus(ctx, image.PropertyID, images[0].ID); err != nil {
                return err
            }
        }
    }
	return nil
}
