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

func (i *ImageUseCase) SaveImage(ctx context.Context, fileName string, image *models.Image) (*models.Image, error) {
    if fileName == "" {
        return nil, errors.New("file name cannot be empty")
    }
    if image == nil {
        return nil, errors.New("image cannot be nil")
    }

    ext := filepath.Ext(fileName)
    allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
    if !allowedExts[ext] {
        return nil, errors.New("unsupported file type")
    }

    propertyID := fmt.Sprintf("%d", image.PropertyID)
    newFileName := uuid.New().String() + ext

    propertyDir := filepath.Join("/app/uploads/properties", propertyID)
    if err := os.MkdirAll(propertyDir, os.ModePerm); err != nil {
        return nil, errors.New("could not create upload directory")
    }

    image.Path = fmt.Sprintf("/uploads/properties/%s/%s", propertyID, newFileName)

    savedImage, err := i.imageRepo.SaveImage(ctx, image)
    if err != nil {
        return nil, err
    }
    return savedImage, nil
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
	return i.imageRepo.DeleteImage(ctx, id)
}
