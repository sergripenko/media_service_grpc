package images

import (
	"context"
	mediaservice "media_service_grpc/proto"
)

type UseCase interface {
	CreateImage(ctx context.Context, data *mediaservice.CreateImageData) (*mediaservice.Response, error)
	GetAllImages(ctx context.Context, userId *mediaservice.UserId) (*mediaservice.Response, error)
	UpdateImage(ctx context.Context, data *mediaservice.UpdateImageData) (*mediaservice.Response, error)
}
