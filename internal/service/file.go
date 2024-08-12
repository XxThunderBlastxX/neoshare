package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/XxThunderBlastxX/neoshare/internal/model"
	"github.com/XxThunderBlastxX/neoshare/internal/repository"
)

type fileService struct {
	ctx       context.Context
	query     *repository.Queries
	db        *sql.DB
	s3Service S3Service
}

type FileService interface {
	UploadFile(key string, contentType string, fileName string, object []byte) error
	DownloadFile(key string) ([]byte, error)
	SyncFileWithDB(file model.File) error
	GetFiles() ([]model.File, error)
}

func NewFileService(ctx context.Context, query *repository.Queries, db *sql.DB, service S3Service) FileService {
	return &fileService{
		ctx:       ctx,
		query:     query,
		s3Service: service,
		db:        db,
	}
}

func (f *fileService) UploadFile(key string, contentType string, fileName string, object []byte) error {
	tx, _ := f.db.BeginTx(f.ctx, nil)
	f.query.WithTx(tx)

	// Note: Not handling error because we are considering there won't
	// any error as the db is locally hosted and error can only occur during uploading
	// the file to the s3 bucket.
	_, err := f.query.CreateFile(f.ctx, repository.CreateFileParams{
		Name:         fileName,
		Key:          key,
		Size:         int32(len(object)),
		LastModified: time.Now(),
	})
	if err != nil {
		// TODO: Handle the error in better way.
		log.Println(err)
	}

	// If the error occurs during the file upload, we will rollback the transaction.
	err = f.s3Service.UploadFile(key, contentType, fileName, object)
	if err != nil {
		// Rollback the transaction if the file is not uploaded successfully.
		_ = tx.Rollback() // TODO: Handle the error in better way.
		return err
	}

	// If the file is uploaded successfully, we will commit the transaction.
	return tx.Commit()
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
