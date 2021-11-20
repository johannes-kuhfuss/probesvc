package app

import (
	"github.com/johannes-kuhfuss/probesvc/controller"
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/service"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

func mapUrls() {
	logger.Debug("Mapping URLs")
	//jh := controller.JobHandlers{Service: service.NewJobService(domain.NewJobRepositoryStub())}
	jh := controller.JobHandlers{Service: service.NewJobService(domain.NewJobRepositoryDb())}
	router.GET("/jobs", jh.GetAllJobs)
	logger.Debug("Done mapping URLs")
}
