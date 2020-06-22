package images

import (
	"context"
	mediaservice "media_service_grpc/proto"
)

type UseCase interface {
	CreateImage(context.Context, *mediaservice.CreateImageData) (*mediaservice.Response, error)
	GetAllImages(context.Context, *mediaservice.UserId) (*mediaservice.Response, error)
	UpdateImage(context.Context, *mediaservice.UpdateImageData) (*mediaservice.Response, error)
}
