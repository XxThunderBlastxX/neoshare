package service

import (
	"context"
	"fmt"
	"github.com/XxThunderBlastxX/neoshare/internal/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

type s3Service struct {
	config *config.S3Config
	client *s3.Client
}

type S3Service interface {
	UploadFile(key *string, object io.Reader) error
}

func New(c *config.S3Config) S3Service {
	option := s3.Options{
		BaseEndpoint:       &c.Endpoint,
		Credentials:        credentials.NewStaticCredentialsProvider(c.AccessKey, c.SecretKey, ""),
		EndpointResolverV2: s3.NewDefaultEndpointResolverV2(),
		Region:             "auto",
	}

	return &s3Service{
		config: c,
		client: s3.New(option),
	}
}

func (s *s3Service) UploadFile(key *string, object io.Reader) error {
	uploader := manager.NewUploader(s.client)

	opt, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.config.Bucket,
		Key:    key,
		Body:   object,
	})
	if err != nil {
		return err
	}

	fmt.Println(opt)

	return nil
}
