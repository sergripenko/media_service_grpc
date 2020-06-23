package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	ms "media_service_grpc/proto"

	"github.com/labstack/gommon/log"

	"google.golang.org/grpc"
)

func main() {
	// connect to grpc
	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Error("cant connect to grpc")
		return
	}
	defer grcpConn.Close()

	mediaManager := ms.NewMediaServiceClient(grcpConn)
	ctx := context.Background()

	// create image
	createImgFile, err := ioutil.ReadFile("client/create_image.json")

	if err != nil {
		log.Error(err)
		return
	}
	var newImg ms.CreateImageData

	if err = json.Unmarshal(createImgFile, &newImg); err != nil {
		log.Error(err)
		return
	}
	resp, err := mediaManager.CreateImage(ctx, &newImg)

	if err != nil {
		log.Error(err)
		return
	}
	log.Info(resp)

	// get all images by userId
	resp, err = mediaManager.GetAllImages(ctx, &ms.UserId{UserId: "5ef1dc2822c13a0d3ab53fba"})
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(resp)

	// update image
	data := &ms.UpdateImageData{
		ImageId: resp.Images[0].Id,
		ImageData: &ms.CreateImageData{
			Height: 400,
			Width:  400,
		},
	}
	resp, err = mediaManager.UpdateImage(ctx, data)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(resp)
}
