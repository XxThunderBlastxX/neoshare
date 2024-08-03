package routes

import "github.com/XxThunderBlastxX/neoshare/internal/handler"

func (r *Router) S3Router() {
	api := r.app.Group("/api")

	s3Handler := handler.NewS3Handler(r.s3service)

	// API
	api.Post("/s3/upload", s3Handler.UploadHandler())
}
