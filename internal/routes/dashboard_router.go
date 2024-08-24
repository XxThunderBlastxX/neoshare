package routes

import "github.com/XxThunderBlastxX/neoshare/internal/handler"

func (r *Router) DashboardRouter() {
	api := r.app.Group("/api")
	view := r.app.Group("/")

	h := handler.NewDashboardHandler(r.s3service, r.fileService)

	// API Routes
	api.Post("/upload", r.middleware.VerifyToken(), h.UploadHandler())
	api.Get("/download/:key", h.DownloadHandler())
	api.Delete("/delete/:key", h.DeleteFileHandler())

	// View Routes
	view.Get("/dashboard", r.middleware.VerifyToken(), h.DashboardView())
	view.Get("/files", r.middleware.VerifyToken(), h.FilesView())
	view.Get("/v/:key", h.DownloadHandler())
}
