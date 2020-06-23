package services

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/labstack/gommon/log"

	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const bucketImgPath = "/images/"

func getCredentialsFromFile() *aws.Config {
	creds := credentials.NewSharedCredentials("config/credentials.sh", "default")

	return &aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: creds,
	}
}

func SaveImageToS3(file io.Reader, filename string) (url string, err error) {
	// The session the S3 Uploader will use
	config := getCredentialsFromFile()
	sess := session.Must(session.NewSession(config))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(viper.GetString("bucket")),
		Key:    aws.String(bucketImgPath + filename),
		Body:   file,
	})
	if err != nil {
		log.Error(err)
		return "", err
	}
	return result.Location, nil
}

func GetImageFromS3(file *os.File, filename string) (err error) {
	// The session the S3 Downloader will use
	config := getCredentialsFromFile()
	sess := session.Must(session.NewSession(config))

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Write the contents of S3 Object to the file
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("bucket")),
		Key:    aws.String(bucketImgPath + filename),
	})
	if err != nil {
		return errors.New("failed to download file")
	}
	return nil
}

func DeleteImageFromS3(filename string) (err error) {
	config := getCredentialsFromFile()

	svc := s3.New(session.New(config))

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(viper.GetString("bucket")),
		Key:    aws.String(bucketImgPath + filename),
	}

	_, err = svc.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	return
}
