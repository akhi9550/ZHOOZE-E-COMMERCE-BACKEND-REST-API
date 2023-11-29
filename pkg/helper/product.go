package helper

import (
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

func GetImageMimeType(filename string) string {
	extension := strings.ToLower(strings.Split(filename, ".")[len(strings.Split(filename, "."))-1])

	imageMimeTypes := map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"bmp":  "image/bmp",
		"webp": "image/webp",
	}

	if mimeType, ok := imageMimeTypes[extension]; ok {
		return mimeType
	}

	return "application/octet-stream"
}

func AddImageToS3(file *multipart.FileHeader) (string, error) {
	f, openErr := file.Open()
	if openErr != nil {

		fmt.Println("opening error:", openErr)
		return "", openErr
	}
	defer f.Close()
	if err := godotenv.Load(); err != nil {
		fmt.Println("error 1", err)
		return "", err
	}
	mimeType := GetImageMimeType(file.Filename)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		fmt.Println("error in session config", err)
		return "", err
	}
	// Create an S3 uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	BucketName := "zhooze"
	//upload data(video or image)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(BucketName),
		Key:         aws.String(file.Filename),
		Body:        f,
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		fmt.Println("error 2", err)
		return "", err
	}
	url := fmt.Sprintf("https://d2jkb5eqmpty2t.cloudfront.net/%s", file.Filename)
	return url, nil
}
