package main

import (
	"context"
	"fmt"
	"log"
	media_service "media_service_grpc/proto"
	"net"
	"time"

	immongo "media_service_grpc/images/repository/mongo"
	imuc "media_service_grpc/images/usecase"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

func main() {
	db := initDB()
	imageRepo := immongo.NewImageRepository(db, viper.GetString("mongo.images_collection"))
	imageUC := imuc.NewImageUseCase(imageRepo)

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("cant listet port", err)
	}

	server := grpc.NewServer()
	media_service.RegisterMediaServiceServer(server, imageUC)

	fmt.Println("starting server at :8081")
	server.Serve(lis)
}

func initDB() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		log.Fatalf("Error occured while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(viper.GetString("mongo.name"))
}
