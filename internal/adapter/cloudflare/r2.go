package cloudflare

import (
	"bwanews/config"
	"bwanews/internal/core/domain/entity"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2/log"
)

var code string

type CloudflareR2Adapter interface {
	UploadImage(req *entity.FileUploadEntity) (string, error)
}

type cloudflareR2Adapter struct {
	Client  *s3.Client
	Bucket  string
	Baseurl string
}

func (c *cloudflareR2Adapter) UploadImage(req *entity.FileUploadEntity) (string, error) {
	openedFile, err := os.Open(req.Path)
	if err != nil {
		code = "[CLOUDFLARE R2] UploadImage - 1"
		log.Errorw(code, err)
		return "", err
	}

	defer openedFile.Close()

	buffer := make([]byte, 512)
	n, err := openedFile.Read(buffer)
	if err != nil && err != io.EOF {
		log.Errorw("[CLOUDFLARE R2] ReadBuffer Error", err)
		return "", err
	}

	contentType := http.DetectContentType(buffer[:n])

	_, err = openedFile.Seek(0, 0)
	if err != nil {
		log.Errorw("[CLOUDFLARE R2] Seek Error", err)
		return "", err
	}

	_, err = c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(c.Bucket),
		Key:         aws.String(req.Name),
		Body:        openedFile,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		code = "[CLOUDFLARE R2] UploadImage - 2"
		log.Errorw(code, err)
		return "", err
	}

	return fmt.Sprintf("%s/%s", c.Baseurl, req.Name), nil
}

func NewCloudflareR2Adapter(client *s3.Client, cfg *config.Config) CloudflareR2Adapter {
	clientBase := s3.NewFromConfig(cfg.LoadAwsConfig(), func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2.AccountID))
	})

	return &cloudflareR2Adapter{
		Client:  clientBase,
		Bucket:  cfg.R2.Name,
		Baseurl: cfg.R2.PublicURL,
	}
}
