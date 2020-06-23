package images

import (
	"context"
	"media_service_grpc/models"
)

type Repository interface {
	CreateImage(ctx context.Context, im *models.Image) (*models.Image, error)
	GetAllImages(ctx context.Context, userId string) ([]*models.Image, error)
	GetImageById(ctx context.Context, imgId string) (*models.Image, error)
}
