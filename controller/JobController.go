package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/probesvc/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"github.com/segmentio/ksuid"
)

type JobHandlers struct {
	Service service.JobService
}

func getJobId(jobIdParam string) (string, api_error.ApiErr) {
	jobId, err := ksuid.Parse(jobIdParam)
	if err != nil {
		logger.Error("User Id should be a ksuid", err)
		return "", api_error.NewBadRequestError("user id should be a ksuid")
	}
	return jobId.String(), nil
}

func (jh *JobHandlers) GetAllJobs(c *gin.Context) {
	logger.Debug("Processing get all jobs request")
	jobs, err := jh.Service.GetAllJobs()
	if err != nil {
		logger.Error("Service error while getting all jobs", err)
		c.JSON(err.StatusCode(), err.Message())
		return
	}
	c.JSON(http.StatusOK, jobs)
	logger.Debug("Done processing get all jobs request")
}

func (jh *JobHandlers) GetJobById(c *gin.Context) {
	logger.Debug("Processing get job by id request")
	jobId, err := getJobId(c.Param("job_id"))
	if err != nil {
		c.JSON(err.StatusCode(), err.Message())
		return
	}
	job, err := jh.Service.GetJobById(jobId)
	if err != nil {
		logger.Error("Service error while getting job by id", err)
		c.JSON(err.StatusCode(), err.Message())
		return
	}
	c.JSON(http.StatusOK, job)
	logger.Debug("Done processing get job by id request")
}
