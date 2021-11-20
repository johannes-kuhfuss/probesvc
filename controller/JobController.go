package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/probesvc/service"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type JobHandlers struct {
	Service service.JobService
}

func (jh *JobHandlers) GetAllJobs(c *gin.Context) {
	logger.Debug("Processing job create request")
	jobs, _ := jh.Service.GetAllJobs()
	c.JSON(http.StatusOK, jobs)
	logger.Debug("Done processing job create request")
}
