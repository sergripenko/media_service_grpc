package mongo

import (
	"context"
	"media_service_grpc/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type Image struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	UserId    primitive.ObjectID `bson:"userId"`
	Filename  string             `bson:"filename"`
	Height    int32              `bson:"height"`
	Width     int32              `bson:"width"`
	UniqId    string             `bson:"uniqId"`
	OrigImgId string             `bson:"OrigImgId"`
	Url       string             `bson:"url"`
}

type ImageRepository struct {
	db *mongo.Collection
}

func NewImageRepository(db *mongo.Database, collection string) *ImageRepository {
	return &ImageRepository{
		db: db.Collection(collection),
	}
}

func (r ImageRepository) CreateImage(ctx context.Context, im *models.Image) (*models.Image, error) {
	model := toModel(im)
	res, err := r.db.InsertOne(ctx, model)

	if err != nil {
		return nil, err
	}
	im.Id = res.InsertedID.(primitive.ObjectID).Hex()
	return im, nil
}

func (r ImageRepository) GetAllImages(ctx context.Context, userId string) ([]*models.Image, error) {
	uid, _ := primitive.ObjectIDFromHex(userId)

	cur, err := r.db.Find(ctx, bson.M{
		"userId": uid,
	})
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}
	out := make([]*Image, 0)

	for cur.Next(ctx) {
		im := new(Image)
		err = cur.Decode(im)

		if err != nil {
			return nil, err
		}
		out = append(out, im)
	}
	if err = cur.Err(); err != nil {
		return nil, err
	}
	return toImages(out), nil
}

func (r ImageRepository) GetImageById(ctx context.Context, imgId string) (*models.Image, error) {
	uid, _ := primitive.ObjectIDFromHex(imgId)
	var im models.Image

	err := r.db.FindOne(ctx, bson.M{
		"_id": uid,
	}).Decode(&im)

	return &im, err
}

func toModel(i *models.Image) *Image {
	uid, _ := primitive.ObjectIDFromHex(i.UserId)

	return &Image{
		UserId:    uid,
		Filename:  i.Filename,
		Height:    i.Height,
		Width:     i.Width,
		UniqId:    i.UniqId,
		OrigImgId: i.OrigImgId,
		Url:       i.Url,
	}
}

func toImage(i *Image) *models.Image {
	return &models.Image{
		Id:        i.Id.Hex(),
		UserId:    i.UserId.Hex(),
		Filename:  i.Filename,
		Height:    i.Height,
		Width:     i.Width,
		UniqId:    i.UniqId,
		OrigImgId: i.OrigImgId,
		Url:       i.Url,
	}
}

func toImages(ims []*Image) []*models.Image {
	out := make([]*models.Image, len(ims))

	for i, img := range ims {
		out[i] = toImage(img)
	}
	return out
}
