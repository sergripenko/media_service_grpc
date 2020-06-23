package usecase

import (
	"bytes"
	"context"
	"errors"
	"image"
	"io/ioutil"
	"media_service_grpc/images"
	"media_service_grpc/models"
	ms "media_service_grpc/proto"
	services "media_service_grpc/services/amazon"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/gommon/log"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const localImgPath = "images/"

type ImageUseCase struct {
	imageRepo images.Repository
}

func NewImageUseCase(imageRepo images.Repository) *ImageUseCase {
	return &ImageUseCase{
		imageRepo: imageRepo,
	}
}

func (i ImageUseCase) CreateImage(ctx context.Context, data *ms.CreateImageData) (*ms.Response, error) {
	// create uniq filename for original image
	var uniqFilename uuid.UUID
	var err error

	if uniqFilename, err = uuid.NewUUID(); err != nil {
		log.Error(err)
		return nil, err
	}
	filenameSlice := strings.Split(data.Filename, ".")
	imgFileName := uniqFilename.String() + "." + filenameSlice[len(filenameSlice)-1]
	var imgConfig image.Config

	// decode image parameters
	if imgConfig, _, err = image.DecodeConfig(bytes.NewReader(data.File)); err != nil {
		log.Error(err)
		return nil, err
	}
	// save img to s3
	var url string

	if url, err = services.SaveImageToS3(bytes.NewReader(data.File), imgFileName); err != nil {
		log.Error(err)
		return nil, err
	}
	// save original image in repo
	origImg := &models.Image{
		UserId:   data.UserId,
		Filename: data.Filename,
		Height:   int32(imgConfig.Height),
		Width:    int32(imgConfig.Width),
		UniqId:   uniqFilename.String() + "." + filenameSlice[len(filenameSlice)-1],
		Url:      url,
	}
	var origImgObj *models.Image

	if origImgObj, err = i.imageRepo.CreateImage(ctx, origImg); err != nil {
		log.Error(err)
		return nil, err
	}
	// save original image in directory
	if err = ioutil.WriteFile(localImgPath+imgFileName, data.File, 0644); err != nil {
		return nil, err
	}
	// open original image for resize
	var src image.Image

	if src, err = imaging.Open(localImgPath + imgFileName); err != nil {
		log.Error(err)
		return nil, err
	}
	// resize and save image
	resized := imaging.Resize(src, int(data.Width), int(data.Height), imaging.Lanczos)
	resizedFileName := uniqFilename.String() + "_" + strconv.Itoa(int(data.Width)) + "x" +
		strconv.Itoa(int(data.Height)) + "." + filenameSlice[len(filenameSlice)-1]

	if err = imaging.Save(resized, localImgPath+resizedFileName); err != nil {
		log.Error(err)
		return nil, err
	}
	// open resized image
	var resizedFile []byte

	if resizedFile, err = ioutil.ReadFile(localImgPath + resizedFileName); err != nil {
		log.Error(err)
		return nil, err
	}
	var resizedUrl string

	if resizedUrl, err = services.SaveImageToS3(bytes.NewReader(resizedFile), resizedFileName); err != nil {
		log.Error(err)
		return nil, err
	}
	// create and save resized image object
	resizedImg := &models.Image{
		UserId: data.UserId,
		Filename: filenameSlice[0] + "_" + strconv.Itoa(int(data.Width)) + "x" +
			strconv.Itoa(int(data.Height)) + "." + filenameSlice[len(filenameSlice)-1],
		Height:    data.Height,
		Width:     data.Width,
		UniqId:    resizedFileName,
		OrigImgId: origImgObj.Id,
		Url:       resizedUrl,
	}
	var resizedImgObj *models.Image
	if resizedImgObj, err = i.imageRepo.CreateImage(ctx, resizedImg); err != nil {
		log.Error(err)
		return nil, err
	}
	// delete files from directory
	err = os.Remove(localImgPath + imgFileName)
	err = os.Remove(localImgPath + resizedFileName)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp := &ms.Response{
		Images: toResp(origImgObj, resizedImgObj),
	}
	return resp, nil
}

func (i ImageUseCase) GetAllImages(ctx context.Context, userId *ms.UserId) (*ms.Response, error) {
	imagesObj, err := i.imageRepo.GetAllImages(ctx, userId.UserId)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp := &ms.Response{
		Images: toResp(imagesObj...),
	}
	return resp, nil
}

func (i ImageUseCase) UpdateImage(ctx context.Context, data *ms.UpdateImageData) (*ms.Response, error) {
	var imageObject *models.Image
	var err error

	if imageObject, err = i.imageRepo.GetImageById(ctx, data.ImageId); err != nil {
		log.Error(err)
		return nil, err
	}

	// Create a file to write the S3 Object contents to.
	file, err := os.Create(localImgPath + imageObject.UniqId)
	defer file.Close()

	if err != nil {
		err = errors.New("failed to create file")
		log.Error(err)
		return nil, err
	}
	if err = services.GetImageFromS3(file, imageObject.UniqId); err != nil {
		log.Error(err)
		return nil, err
	}
	// open original image for resize
	var src image.Image

	if src, err = imaging.Open(localImgPath + imageObject.UniqId); err != nil {
		log.Error(err)
		return nil, err
	}
	// resize and save image
	resized := imaging.Resize(src, int(data.ImageData.Width), int(data.ImageData.Height), imaging.Lanczos)
	fullFilenameSlice := strings.Split(imageObject.Filename, ".")
	// create uniq filename for original image
	var uniqFilename uuid.UUID

	if uniqFilename, err = uuid.NewUUID(); err != nil {
		log.Error(err)
		return nil, err
	}
	resizedFileName := uniqFilename.String() + "_" + strconv.Itoa(int(data.ImageData.Width)) + "x" +
		strconv.Itoa(int(data.ImageData.Height)) + "." + fullFilenameSlice[len(fullFilenameSlice)-1]

	if err = imaging.Save(resized, localImgPath+resizedFileName); err != nil {
		log.Error(err)
		return nil, err
	}
	// open resized image
	var resizedFile []byte

	if resizedFile, err = ioutil.ReadFile(localImgPath + resizedFileName); err != nil {
		log.Error(err)
		return nil, err
	}
	var url string

	if url, err = services.SaveImageToS3(bytes.NewReader(resizedFile), resizedFileName); err != nil {
		log.Error(err)
		return nil, err
	}
	// delete local files
	err = os.Remove(localImgPath + resizedFileName)
	err = os.Remove(localImgPath + imageObject.UniqId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	newResizedImg := &models.Image{
		UserId: imageObject.UserId,
		Filename: fullFilenameSlice[0] + "_" + strconv.Itoa(int(data.ImageData.Width)) + "x" +
			strconv.Itoa(int(data.ImageData.Height)) + "." + fullFilenameSlice[len(fullFilenameSlice)-1],
		Height:    data.ImageData.Height,
		Width:     data.ImageData.Width,
		UniqId:    resizedFileName,
		OrigImgId: data.ImageId,
		Url:       url,
	}
	var newImgObj *models.Image

	if newImgObj, err = i.imageRepo.CreateImage(ctx, newResizedImg); err != nil {
		log.Error(err)
		return nil, err
	}
	resp := &ms.Response{
		Images: toResp(newImgObj),
	}
	return resp, nil
}

func toResp(images ...*models.Image) []*ms.Image {
	out := make([]*ms.Image, len(images))

	for i, img := range images {
		out[i] = &ms.Image{
			Id:        img.Id,
			UserId:    img.UserId,
			Filename:  img.Filename,
			Height:    img.Height,
			Width:     img.Width,
			UniqId:    img.UniqId,
			OrigImgId: img.OrigImgId,
			Url:       img.Url,
		}
	}
	return out
}
