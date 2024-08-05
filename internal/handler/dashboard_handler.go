package handler

import (
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/XxThunderBlastxX/neoshare/cmd/web/page"
	"github.com/XxThunderBlastxX/neoshare/internal/model"
	"github.com/XxThunderBlastxX/neoshare/internal/service"
	"github.com/XxThunderBlastxX/neoshare/internal/utils"
)

type dashboardHandler struct {
	s3Service service.S3Service
}

type DashboardHandler interface {
	UploadHandler() fiber.Handler
	DownloadHandler() fiber.Handler

	DashboardView() fiber.Handler
}

func NewDashboardHandler(s3Service service.S3Service) DashboardHandler {
	return &dashboardHandler{
		s3Service: s3Service,
	}
}

func (d *dashboardHandler) DashboardView() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		render := adaptor.HTTPHandler(templ.Handler(page.DashboardPage()))

		return render(ctx)
	}
}

func (d *dashboardHandler) UploadHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[*model.ErrorResponse]{
				Error:   err.Error(),
				Success: false,
			})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[*model.ErrorResponse]{
				Error:   err.Error(),
				Success: false,
			})
		}
		defer file.Close()

		key := utils.GenerateUID(strconv.FormatInt(time.Now().UnixNano(), 10))

		err = d.s3Service.UploadFile(&key, file)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[*model.ErrorResponse]{
				Error:   err.Error(),
				Success: false,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.SuccessResponse]{
			Success: true,
		})
	}
}

func (d *dashboardHandler) DownloadHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		key := ctx.Params("key")

		file, err := d.s3Service.DownloadFile(&key)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[*model.ErrorResponse]{
				Error:   err.Error(),
				Success: false,
			})
		}

		// TODO: Set the content type based on the file type
		ctx.Set("Content-Type", "image/png")

		return ctx.Status(fiber.StatusOK).Send(file)
	}
}
