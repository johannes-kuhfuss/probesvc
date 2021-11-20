package app

import (
	"github.com/johannes-kuhfuss/probesvc/controller"
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/domain/service"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

func mapUrls() {
	logger.Debug("Mapping URLs")
	jh := controller.JobHandlers{service.NewJobService(domain.NewJobRepositoryStub())}
	router.GET("/jobs", jh.GetAllJobs)
	logger.Debug("Done mapping URLs")
}
