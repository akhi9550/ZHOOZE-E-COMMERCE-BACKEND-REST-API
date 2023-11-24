package helper

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func AddImageToS3(file *multipart.FileHeader) (string, error) {
	f, openErr := file.Open()
	if openErr != nil {
		fmt.Println("opening error:", openErr)
		return "", openErr
	}
	defer f.Close()
	if err := godotenv.Load(); err != nil {
		return "", err
	}

	// Create a new AWS session with the loaded access keys
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
		// Add more configurations if needed.
	})
	if err != nil {
		return "", err
	}
	// Create an S3 uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	BucketName := "zhooze"
	// Upload the video data to S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(file.Filename),
		Body:   f,
	})
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", BucketName, file.Filename)
	return url, nil
}
