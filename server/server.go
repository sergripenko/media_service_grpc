package main

import (
	"context"
	"media_service_grpc/config"
	ms "media_service_grpc/proto"
	"net"
	"time"

	"github.com/labstack/gommon/log"

	imMongo "media_service_grpc/images/repository/mongo"
	imUC "media_service_grpc/images/usecase"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

func main() {
	if err := config.Init(); err != nil {
		log.Error(err)
	}
	db := initDB()
	imageRepo := imMongo.NewImageRepository(db, viper.GetString("mongo.images_collection"))
	imageUC := imUC.NewImageUseCase(imageRepo)

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Error("cant listet port", err)
	}
	server := grpc.NewServer()
	ms.RegisterMediaServiceServer(server, imageUC)

	log.Info("starting server at :8081")

	if err = server.Serve(lis); err != nil {
		log.Error(err)
	}
}

func initDB() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		log.Error("Error occured while establishing connection to mongoDB")
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		log.Error(err)
		return nil
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		log.Error(err)
		return nil
	}
	return client.Database(viper.GetString("mongo.name"))
}
