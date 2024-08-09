package service

import (
	"context"
	"io"

	"github.com/XxThunderBlastxX/neoshare/internal/model"
	"github.com/XxThunderBlastxX/neoshare/internal/repository"
)

type fileService struct {
	ctx       context.Context
	query     *repository.Queries
	s3Service S3Service
}

type FileService interface {
	UploadFile(key string, contentType string, fileName string, object io.Reader) error
	DownloadFile(key string) ([]byte, error)
	SyncFileWithDB(file model.File) error
	GetFiles() ([]model.File, error)
}

func NewFileService(ctx context.Context, query *repository.Queries, service S3Service) FileService {
	return &fileService{
		ctx:       ctx,
		query:     query,
		s3Service: service,
	}
}

func (f *fileService) UploadFile(key string, contentType string, fileName string, object io.Reader) error {
	return f.s3Service.UploadFile(key, contentType, fileName, object)
}

func (f *fileService) DownloadFile(key string) ([]byte, error) {
	return f.s3Service.DownloadFile(key)
}

func (f *fileService) SyncFileWithDB(file model.File) error {
	fileParams := repository.CreateFileParams{
		Name:         file.Name,
		Key:          file.Key,
		Size:         int32(file.Size),
		LastModified: file.LastModified,
	}

	_, err := f.query.CreateFile(f.ctx, fileParams)
	if err != nil {
		return err
	}

	return nil
}

func (f *fileService) GetFiles() ([]model.File, error) {
	res, err := f.query.ListFiles(f.ctx)
	if err != nil {
		return nil, err
	}

	files := make([]model.File, len(res))
	for idx, repoFile := range res {
		files[idx].FromRepositoryFile(repoFile)
	}

	return files, nil
}
