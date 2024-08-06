package routes

import "github.com/XxThunderBlastxX/neoshare/internal/handler"

func (r *Router) DashboardRouter() {
	api := r.app.Group("/api")
	view := r.app.Group("/")

	h := handler.NewDashboardHandler(r.s3service)

	// API Routes
	api.Post("/upload", h.UploadHandler())
	api.Get("/download/:key", h.DownloadHandler())

	// View Routes
	view.Get("/dashboard", h.DashboardView())
	view.Get("/files", h.FilesView())
	view.Get("/:key", h.DownloadHandler())
}
