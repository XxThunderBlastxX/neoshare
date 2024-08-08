package handler

import (
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/sujit-baniya/flash"

	"github.com/XxThunderBlastxX/neoshare/cmd/web/component"
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
	FilesView() fiber.Handler
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
		render := adaptor.HTTPHandler(templ.Handler(page.DashboardPage()))

		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			render = adaptor.HTTPHandler(templ.Handler(page.UploadSection(model.WebResponse{
				Success:    false,
				StatusCode: fiber.StatusNotFound,
				Message:    "No file selected to upload",
			})))
			return render(ctx)
		}

		file, err := fileHeader.Open()
		if err != nil {
			render = adaptor.HTTPHandler(templ.Handler(page.UploadSection(model.WebResponse{
				Message:    "Error occurred while opening the file",
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			})))
			return render(ctx)
		}
		defer file.Close()

		key := utils.GenerateUID(strconv.FormatInt(time.Now().UnixNano(), 10))

		contentType := fileHeader.Header.Get("Content-Type")

		err = d.s3Service.UploadFile(key, contentType, fileHeader.Filename, file)
		if err != nil {
			render = adaptor.HTTPHandler(templ.Handler(page.UploadSection(model.WebResponse{
				Message:    "Error occurred while uploading the file",
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			})))
			return render(ctx)
		}

		render = adaptor.HTTPHandler(templ.Handler(component.ShortLinkView(ctx.BaseURL() + "/" + key)))

		return render(ctx)
	}
}

func (d *dashboardHandler) DownloadHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		key := ctx.Params("key")

		file, err := d.s3Service.DownloadFile(key)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			})
		}

		// TODO: Set the content type based on the file type
		ctx.Set("Content-Type", "image/png")

		return ctx.Status(fiber.StatusOK).Send(file)
	}
}

func (d *dashboardHandler) FilesView() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		res := flash.Get(ctx)
		if len(res) != 0 {
			var resData model.WebResponse
			resData.ConvertToStruct(res)
			render := adaptor.HTTPHandler(templ.Handler(page.FilesPage([]model.File{}, resData)))
			return render(ctx)
		}

		files, err := d.s3Service.GetFiles()
		if err != nil {
			errRes := model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			}
			return flash.WithError(ctx, errRes.ConvertToMap()).Redirect("/files")
		}
		render := adaptor.HTTPHandler(templ.Handler(page.FilesPage(files)))

		return render(ctx)
	}
}
