package service

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/XxThunderBlastxX/neoshare/internal/config"
)

type s3Service struct {
	config *config.S3Config
	client *s3.Client
}

type S3Service interface {
	UploadFile(key *string, object io.Reader) error
	DownloadFile(key *string) ([]byte, error)
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
	uploader := manager.NewUploader(s.client, func(u *manager.Uploader) {
		u.Concurrency = 5
		u.S3 = s.client
		u.PartSize = 20 * 1024 * 1024
	})

	_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: &s.config.Bucket,
		Key:    key,
		Body:   object,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *s3Service) DownloadFile(key *string) ([]byte, error) {
	downloader := manager.NewDownloader(s.client)

	buff := manager.NewWriteAtBuffer([]byte{})

	_, err := downloader.Download(context.Background(), buff, &s3.GetObjectInput{
		Bucket: &s.config.Bucket,
		Key:    key,
	})
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
