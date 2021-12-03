package service

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	realdomain "github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/dto"
	"github.com/johannes-kuhfuss/probesvc/mocks/domain"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	jobCtrl     *gomock.Controller
	mockJobRepo *domain.MockJobRepository
	jobService  JobService
)

func setupJob(t *testing.T) func() {
	jobCtrl = gomock.NewController(t)
	mockJobRepo = domain.NewMockJobRepository(jobCtrl)
	jobService = NewJobService(mockJobRepo)
	return func() {
		jobService = nil
		jobCtrl.Finish()
	}
}

func Test_GetAllJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	apiError := api_error.NewNotFoundError("no jobs found")
	mockJobRepo.EXPECT().FindAll("").Return(nil, apiError)

	result, err := jobService.GetAllJobs("")

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, apiError.Message(), err.Message())
}

func Test_GetAllJobs_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	job1, _ := realdomain.NewJob("job 1", "url1")
	job2, _ := realdomain.NewJob("job 2", "url2")
	jobs := make([]realdomain.Job, 0)
	jobs = append(jobs, *job1)
	jobs = append(jobs, *job2)
	jobResult := make([]dto.JobResponse, 0)
	jobResult = append(jobResult, job1.ToDto())
	jobResult = append(jobResult, job2.ToDto())

	mockJobRepo.EXPECT().FindAll("").Return(&jobs, nil)

	result, err := jobService.GetAllJobs("")

	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.Equal(t, result, &jobResult)
}

func Test_GetJobById_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	id := ksuid.New().String()
	apiError := api_error.NewNotFoundError(fmt.Sprintf("job with id %v not found", id))
	mockJobRepo.EXPECT().FindById(id).Return(nil, apiError)

	result, err := jobService.GetJobById(id)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, apiError.Message(), err.Message())
}

func Test_GetJobById_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	jobResp := newJob.ToDto()
	id := newJob.Id.String()
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)

	result, err := jobService.GetJobById(id)

	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.Equal(t, result, &jobResp)
}

func Test_CreateJob_Returns_BaqRequestError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	jobReq := dto.NewJobRequest{
		Name:   "job 1",
		SrcUrl: "",
	}
	result, err := jobService.CreateJob(jobReq)
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Job must have a source URL", err.Message())
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
}

func Test_CreateJob_Returns_InternalServerError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	jobReq := dto.NewJobRequest{
		Name:   "job 1",
		SrcUrl: "url 1",
	}
	apiError := api_error.NewInternalServerError("database error", nil)
	mockJobRepo.EXPECT().Save(gomock.Any()).Return(apiError)

	result, err := jobService.CreateJob(jobReq)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "database error", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
}

func Test_CreateJob_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	jobReq := dto.NewJobRequest{
		Name:   "job 1",
		SrcUrl: "url 1",
	}
	mockJobRepo.EXPECT().Save(gomock.Any()).Return(nil)

	result, err := jobService.CreateJob(jobReq)

	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.EqualValues(t, jobReq.Name, result.Name)
	assert.EqualValues(t, jobReq.SrcUrl, result.SrcUrl)
	assert.EqualValues(t, "created", result.Status)
}

func Test_DeleteJobById_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	id := ksuid.New().String()
	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	mockJobRepo.EXPECT().FindById(id).Return(nil, apiError)

	err := jobService.DeleteJobById(id)

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_DeleteJobById_Returns_InternalServerError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	apiError := api_error.NewInternalServerError("database error", nil)
	mockJobRepo.EXPECT().DeleteById(id).Return(apiError)

	err := jobService.DeleteJobById(id)

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
}

func Test_DeleteJobById_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobRepo.EXPECT().DeleteById(id).Return(nil)

	err := jobService.DeleteJobById(id)

	assert.Nil(t, err)
}

func Test_GetNextJob_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	apiError := api_error.NewNotFoundError("No next job found")
	mockJobRepo.EXPECT().GetNext().Return(nil, apiError)

	job, err := jobService.GetNextJob()

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_GetNextJob_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	nextJob, _ := realdomain.NewJob("job 1", "url 1")
	mockJobRepo.EXPECT().GetNext().Return(nextJob, nil)

	job, err := jobService.GetNextJob()

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, nextJob.Name, job.Name)
	assert.EqualValues(t, nextJob.SrcUrl, job.SrcUrl)
}

func Test_SetStatus_NoJobWithId_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	id := ksuid.New().String()
	updReq := dto.JobStatusUpdateRequest{
		Status: "",
		ErrMsg: "",
	}
	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	mockJobRepo.EXPECT().FindById(id).Return(nil, apiError)

	err := jobService.SetStatus(id, updReq)

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_SetStatus_WrongStatusValue_Returns_BadRequestError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	updReq := dto.JobStatusUpdateRequest{
		Status: "wrong_value",
		ErrMsg: "",
	}

	err := jobService.SetStatus(id, updReq)

	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("Could not parse status value %v", updReq.Status), err.Message())
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
}

func Test_SetStatus_StatusFailed_Returns_Error(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	updReq := dto.JobStatusUpdateRequest{
		Status: "failed",
		ErrMsg: "failure_reason",
	}
	updReqParsed, _ := realdomain.ParseStatusRequest(updReq)
	apiError := api_error.NewInternalServerError("something bad happened", nil)
	mockJobRepo.EXPECT().SetStatus(id, *updReqParsed).Return(apiError)

	err := jobService.SetStatus(id, updReq)

	assert.NotNil(t, err)
	assert.EqualValues(t, "something bad happened", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
}

func Test_SetStatus_StatusFailed_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	updReq := dto.JobStatusUpdateRequest{
		Status: "failed",
		ErrMsg: "failure_reason",
	}
	updReqParsed, _ := realdomain.ParseStatusRequest(updReq)
	mockJobRepo.EXPECT().SetStatus(id, *updReqParsed).Return(nil)

	err := jobService.SetStatus(id, updReq)

	assert.Nil(t, err)
}

func Test_SetResult_NoJobWithId_Returns_NotFoundError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	id := ksuid.New().String()
	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	mockJobRepo.EXPECT().FindById(id).Return(nil, apiError)

	err := jobService.SetResult(id, "new result data")

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_SetResult_Returns_Error(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	apiError := api_error.NewInternalServerError("something went wrong", nil)
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobRepo.EXPECT().SetResult(id, "new result data").Return(apiError)

	err := jobService.SetResult(id, "new result data")

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
}

func Test_SetResult_Returns_NoError(t *testing.T) {
	teardown := setupJob(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	mockJobRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobRepo.EXPECT().SetResult(id, "new result data").Return(nil)

	err := jobService.SetResult(id, "new result data")

	assert.Nil(t, err)
}
