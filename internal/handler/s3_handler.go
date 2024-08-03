package handler

import (
	"github.com/XxThunderBlastxX/neoshare/internal/service"
	"github.com/gofiber/fiber/v2"
)

type s3Handler struct {
	s3Service service.S3Service
}

type S3Handler interface {
	UploadHandler() fiber.Handler
	DownloadHandler() fiber.Handler
}

func NewS3Handler(s3Service service.S3Service) S3Handler {
	return &s3Handler{
		s3Service: s3Service,
	}
}

func (s *s3Handler) UploadHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		defer file.Close()

		key := "kuchbhi8"

		err = s.s3Service.UploadFile(&key, file)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "File uploaded successfully",
		})
	}
}

func (s *s3Handler) DownloadHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		key := ctx.Params("key")

		file, err := s.s3Service.DownloadFile(&key)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// TODO: Set the content type based on the file type
		ctx.Set("Content-Type", "image/png")

		return ctx.Status(fiber.StatusOK).Send(file)
	}

}
