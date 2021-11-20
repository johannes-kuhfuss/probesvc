package app

import (
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/handler"
	"github.com/johannes-kuhfuss/probesvc/service"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

func mapUrls() {
	logger.Debug("Mapping URLs")
	jh := handler.JobHandlers{Service: service.NewJobService(domain.NewJobRepositoryDb())}
	router.GET("/jobs", jh.GetAllJobs)
	router.GET("jobs/:job_id", jh.GetJobById)
	logger.Debug("Done mapping URLs")
}
