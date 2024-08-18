package service

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/XxThunderBlastxX/neoshare/internal/config"
	"github.com/XxThunderBlastxX/neoshare/internal/model"
	"github.com/XxThunderBlastxX/neoshare/internal/utils"
)

type s3Service struct {
	config *config.S3Config
	client *s3.Client
}

type S3Service interface {
	UploadFile(key string, contentType string, fileName string, object []byte) error
	DownloadFile(key string) ([]byte, error)
	GetFiles() ([]model.File, error)
	GetFileNameAndType(key string) (string, string, error)
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

func (s *s3Service) UploadFile(key, contentType, fileName string, object []byte) error {
	uploader := manager.NewUploader(s.client, func(u *manager.Uploader) {
		u.Concurrency = 5
		u.S3 = s.client
		u.PartSize = 20 * 1024 * 1024
	})

	// Defining metadata for the file.
	metaData := map[string]string{
		"filename": utils.RemoveNonASCIIValue(fileName),
	}

	_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(s.config.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(object),
		ContentType: aws.String(contentType),
		Metadata:    metaData,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *s3Service) DownloadFile(key string) ([]byte, error) {
	downloader := manager.NewDownloader(s.client)

	buff := manager.NewWriteAtBuffer([]byte{})

	_, err := downloader.Download(context.Background(), buff, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (s *s3Service) GetFiles() ([]model.File, error) {
	// TODO : Implement to fetch from DB
	var files []model.File

	objects, err := s.client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: &s.config.Bucket,
	})
	if err != nil {
		return nil, err
	}

	for _, obj := range objects.Contents {
		filename, _, _ := s.GetFileNameAndType(*obj.Key)
		files = append(files, model.File{
			Name:         filename,
			Key:          *obj.Key,
			Size:         *obj.Size,
			LastModified: *obj.LastModified,
		})
	}

	return files, nil
}

func (s *s3Service) GetFileNameAndType(key string) (fileName, contentType string, err error) {
	metaData, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", "", err
	}

	fileName = metaData.Metadata["filename"]
	contentType = *metaData.ContentType

	return fileName, contentType, nil
}
