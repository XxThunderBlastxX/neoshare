package handler

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/XxThunderBlastxX/neoshare/cmd/web/component"
	"github.com/XxThunderBlastxX/neoshare/cmd/web/page"
	"github.com/XxThunderBlastxX/neoshare/internal/model"
	"github.com/XxThunderBlastxX/neoshare/internal/service"
	"github.com/XxThunderBlastxX/neoshare/internal/utils"
)

type dashboardHandler struct {
	s3Service   service.S3Service
	fileService service.FileService
}

// DashboardHandler is an interface that defines the methods for the dashboard handler
type DashboardHandler interface {
	UploadHandler() fiber.Handler
	DownloadHandler() fiber.Handler
	DeleteFileHandler() fiber.Handler

	DashboardView() fiber.Handler
	FilesView() fiber.Handler
}

// NewDashboardHandler is a factory function that returns instance of the DashboardHandler
func NewDashboardHandler(s3Service service.S3Service, fileService service.FileService) DashboardHandler {
	return &dashboardHandler{
		s3Service:   s3Service,
		fileService: fileService,
	}
}

func (d *dashboardHandler) DashboardView() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Rendering the dashboard view page
		render := adaptor.HTTPHandler(templ.Handler(page.DashboardPage()))
		return render(ctx)
	}
}

func (d *dashboardHandler) UploadHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Parsing the form file with the key "file"
		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			log.Printf("Error occurred while parsing the form file: %v", err)
			render := adaptor.HTTPHandler(templ.Handler(page.UploadSection(model.WebResponse{
				Success:    false,
				StatusCode: fiber.StatusNotFound,
				Message:    "No file selected to upload",
			})))
			return render(ctx)
		}

		// Opening the file with the fileHeader
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Error occurred while opening the file: %v", err)
			render := adaptor.HTTPHandler(templ.Handler(page.UploadSection(model.WebResponse{
				Message:    "Error occurred while opening the file",
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			})))
			return render(ctx)
		}
		defer file.Close()

		// Reading the file to get the file buffer
		fileBuff, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Error occurred while reading the file: %v", err)
			render := adaptor.HTTPHandler(templ.Handler(page.UploadSection(model.WebResponse{
				Message:    "Error occurred while reading the file",
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			})))
			return render(ctx)
		}

		// Generating the key and content type for the file
		var (
			key         = utils.GenerateUID(strconv.FormatInt(time.Now().UnixNano(), 10))
			contentType = fileHeader.Header.Get("Content-Type")
		)

		// Uploading the file to the s3 bucket and syncing it with database
		err = d.fileService.UploadFile(key, contentType, fileHeader.Filename, fileBuff)
		if err != nil {
			log.Printf("Error occurred while uploading the file: %v", err)
			render := adaptor.HTTPHandler(templ.Handler(page.UploadSection(model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			})))
			return render(ctx)
		}

		// Rendering the short link view page
		render := adaptor.HTTPHandler(templ.Handler(component.ShortLinkView(ctx.BaseURL() + "/v/" + key)))
		return render(ctx)
	}
}

func (d *dashboardHandler) DownloadHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Getting the key from the url params
		key := ctx.Params("key")

		// Downloading the file from the s3 bucket
		file, err := d.s3Service.DownloadFile(key)
		if err != nil {
			log.Printf("Error occurred while downloading the file: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			})
		}

		// Getting the content type of the file
		filename, contentType, err := d.s3Service.GetFileNameAndType(key)
		if err != nil {
			log.Printf("Error occurred while fetching the file metadata: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			})
		}
		ctx.Set("Content-Type", contentType)
		ctx.Set("content-Disposition", fmt.Sprintf("filename=\"%q\"", filename))

		return ctx.Status(fiber.StatusOK).Send(file)
	}
}

func (d *dashboardHandler) FilesView() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Getting the list of files from the database
		files, err := d.fileService.GetFiles()
		if err != nil {
			log.Printf("Error occurred while fetching the files: %v", err)
			render := adaptor.HTTPHandler(templ.Handler(page.FilesPage([]model.File{}, model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			})))
			return render(ctx)
		}

		// Rendering the files view page
		render := adaptor.HTTPHandler(templ.Handler(page.FilesPage(files)))
		return render(ctx)
	}
}

func (d *dashboardHandler) DeleteFileHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fileKey := ctx.Params("key")

		if fileKey == "" {
			log.Printf("File key is empty")
			render := adaptor.HTTPHandler(templ.Handler(page.FilesPage([]model.File{}, model.WebResponse{
				Message:    "File key is empty",
				StatusCode: fiber.StatusBadRequest,
				Success:    false,
			})))
			return render(ctx)
		}

		err := d.fileService.DeleteFile(fileKey)
		if err != nil {
			log.Printf("Error occurred while deleting the file: %v", err)
			render := adaptor.HTTPHandler(templ.Handler(page.FilesPage([]model.File{}, model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			})))
			return render(ctx)
		}

		files, err := d.fileService.GetFiles()
		if err != nil {
			log.Printf("Error occurred while fetching the files: %v", err)
			render := adaptor.HTTPHandler(templ.Handler(page.FilesPage([]model.File{}, model.WebResponse{
				Message:    err.Error(),
				StatusCode: fiber.StatusInternalServerError,
				Success:    false,
			})))
			return render(ctx)
		}

		render := adaptor.HTTPHandler(templ.Handler(page.FilesSection(files, model.WebResponse{
			Message:    "File deleted successfully",
			StatusCode: fiber.StatusOK,
			Success:    true,
		})))
		return render(ctx)
	}
}
