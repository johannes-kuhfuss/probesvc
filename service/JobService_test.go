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
	ctrl     *gomock.Controller
	mockRepo *domain.MockJobRepository
	service  JobService
)

func setup(t *testing.T) func() {
	ctrl = gomock.NewController(t)
	mockRepo = domain.NewMockJobRepository(ctrl)
	service = NewJobService(mockRepo)
	return func() {
		service = nil
		ctrl.Finish()
	}
}

func Test_GetAllJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	apiError := api_error.NewNotFoundError("no jobs found")
	mockRepo.EXPECT().FindAll("").Return(nil, apiError)

	result, err := service.GetAllJobs("")

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, apiError.Message(), err.Message())
}

func Test_GetAllJobs_Returns_NoError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	job1, _ := realdomain.NewJob("job 1", "url1")
	job2, _ := realdomain.NewJob("job 2", "url2")
	jobs := make([]realdomain.Job, 0)
	jobs = append(jobs, *job1)
	jobs = append(jobs, *job2)
	jobResult := make([]dto.JobResponse, 0)
	jobResult = append(jobResult, job1.ToDto())
	jobResult = append(jobResult, job2.ToDto())

	mockRepo.EXPECT().FindAll("").Return(&jobs, nil)

	result, err := service.GetAllJobs("")

	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.Equal(t, result, &jobResult)
}

func Test_GetJobById_Returns_NotFoundError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	id := ksuid.New().String()
	apiError := api_error.NewNotFoundError(fmt.Sprintf("job with id %v not found", id))
	mockRepo.EXPECT().FindById(id).Return(nil, apiError)

	result, err := service.GetJobById(id)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
	assert.EqualValues(t, apiError.Message(), err.Message())
}

func Test_GetJobById_Returns_NoError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	jobResp := newJob.ToDto()
	id := newJob.Id.String()
	mockRepo.EXPECT().FindById(id).Return(newJob, nil)

	result, err := service.GetJobById(id)

	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.Equal(t, result, &jobResp)
}

func Test_CreateJob_Returns_BaqRequestError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	jobReq := dto.NewJobRequest{
		Name:   "job 1",
		SrcUrl: "",
	}
	result, err := service.CreateJob(jobReq)
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Job must have a source URL", err.Message())
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
}

func Test_CreateJob_Returns_InternalServerError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	jobReq := dto.NewJobRequest{
		Name:   "job 1",
		SrcUrl: "url 1",
	}
	apiError := api_error.NewInternalServerError("database error", nil)
	mockRepo.EXPECT().Save(gomock.Any()).Return(apiError)

	result, err := service.CreateJob(jobReq)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "database error", err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
}

func Test_CreateJob_Returns_NoError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	jobReq := dto.NewJobRequest{
		Name:   "job 1",
		SrcUrl: "url 1",
	}
	mockRepo.EXPECT().Save(gomock.Any()).Return(nil)

	result, err := service.CreateJob(jobReq)

	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.EqualValues(t, jobReq.Name, result.Name)
	assert.EqualValues(t, jobReq.SrcUrl, result.SrcUrl)
	assert.EqualValues(t, "created", result.Status)
}

func Test_DeleteJobById_Returns_NotFoundError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	id := ksuid.New().String()
	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	mockRepo.EXPECT().FindById(id).Return(nil, apiError)

	err := service.DeleteJobById(id)

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_DeleteJobById_Returns_InternalServerError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	mockRepo.EXPECT().FindById(id).Return(newJob, nil)
	apiError := api_error.NewInternalServerError("database error", nil)
	mockRepo.EXPECT().DeleteById(id).Return(apiError)

	err := service.DeleteJobById(id)

	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
}

func Test_DeleteJobById_Returns_NoError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	mockRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockRepo.EXPECT().DeleteById(id).Return(nil)

	err := service.DeleteJobById(id)

	assert.Nil(t, err)
}

func Test_GetNextJob_Returns_NotFoundError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	apiError := api_error.NewNotFoundError("No next job found")
	mockRepo.EXPECT().GetNext().Return(nil, apiError)

	job, err := service.GetNextJob()

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, apiError.Message(), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_GetNextJob_Returns_NoError(t *testing.T) {
	teardown := setup(t)
	defer teardown()
	nextJob, _ := realdomain.NewJob("job 1", "url 1")
	mockRepo.EXPECT().GetNext().Return(nextJob, nil)

	job, err := service.GetNextJob()

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, nextJob.Name, job.Name)
	assert.EqualValues(t, nextJob.SrcUrl, job.SrcUrl)
}
