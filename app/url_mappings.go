package app

import (
	"github.com/johannes-kuhfuss/services_utils/logger"
)

func mapUrls() {
	logger.Debug("Mapping URLs")
	router.GET("/jobs", jobHandler.GetAllJobs)
	router.GET("jobs/:job_id", jobHandler.GetJobById)
	logger.Debug("Done mapping URLs")
}
