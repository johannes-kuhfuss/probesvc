package service

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/johannes-kuhfuss/probesvc/config"
	realdomain "github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/dto"
	"github.com/johannes-kuhfuss/probesvc/mocks/domain"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/stretchr/testify/assert"
)

var (
	fileCtrl        *gomock.Controller
	mockFileRepo    *domain.MockFileRepository
	fileService     FileService
	jobFileCtrl     *gomock.Controller
	mockJobFileRepo *domain.MockJobRepository
	jobFileService  JobService
	probePath       string = "ffprobe.exe"
)

func setupFile(t *testing.T) func() {
	jobFileCtrl = gomock.NewController(t)
	mockJobFileRepo = domain.NewMockJobRepository(jobFileCtrl)
	jobFileService = NewJobService(mockJobFileRepo)
	fileCtrl = gomock.NewController(t)
	mockFileRepo = domain.NewMockFileRepository(fileCtrl)
	fileService = NewFileService(mockFileRepo, jobFileService)
	return func() {
		fileService = nil
		fileCtrl.Finish()
	}
}

func Test_startJob_Returns_NotFoundError(t *testing.T) {
	teardown := setupFile(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	jobReq := newJob.ToDto()

	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	mockJobFileRepo.EXPECT().FindById(id).Return(nil, apiError)

	err := fileService.startJob(&jobReq)

	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("Job with id %v does not exist", id), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_startJob_Returns_NoError(t *testing.T) {
	teardown := setupFile(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	jobReq := newJob.ToDto()
	jobStatus := dto.JobStatusUpdateRequest{
		Status: "running",
		ErrMsg: "",
	}
	jobStatusReq, _ := realdomain.ParseStatusRequest(jobStatus)

	mockJobFileRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobFileRepo.EXPECT().SetStatus(id, *jobStatusReq).Return(nil)

	err := fileService.startJob(&jobReq)

	assert.Nil(t, err)
}

func Test_failJob_Returns_NotFoundError(t *testing.T) {
	teardown := setupFile(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	jobReq := newJob.ToDto()

	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	failErr := api_error.NewBadRequestError("bad request")
	mockJobFileRepo.EXPECT().FindById(id).Return(nil, apiError)

	err := fileService.failJob(&jobReq, failErr)

	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("Job with id %v does not exist", id), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_failJob_Returns_NoError(t *testing.T) {
	teardown := setupFile(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	jobReq := newJob.ToDto()
	jobStatus := dto.JobStatusUpdateRequest{
		Status: "failed",
		ErrMsg: "Error while analyzing file",
	}
	jobStatusReq, _ := realdomain.ParseStatusRequest(jobStatus)
	failErr := api_error.NewBadRequestError("bad request")

	mockJobFileRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobFileRepo.EXPECT().SetStatus(id, *jobStatusReq).Return(nil)

	err := fileService.failJob(&jobReq, failErr)

	assert.Nil(t, err)
}

func Test_finishJob_Returns_NotFoundError(t *testing.T) {
	teardown := setupFile(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	jobReq := newJob.ToDto()

	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	mockJobFileRepo.EXPECT().FindById(id).Return(nil, apiError)

	err := fileService.finishJob(&jobReq)

	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("Job with id %v does not exist", id), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_finishJob_Returns_NoError(t *testing.T) {
	teardown := setupFile(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	jobReq := newJob.ToDto()
	jobStatus := dto.JobStatusUpdateRequest{
		Status: "finished",
		ErrMsg: "",
	}
	jobStatusReq, _ := realdomain.ParseStatusRequest(jobStatus)

	mockJobFileRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobFileRepo.EXPECT().SetStatus(id, *jobStatusReq).Return(nil)

	err := fileService.finishJob(&jobReq)

	assert.Nil(t, err)
}

func Test_addResultToJob_Returns_NotFoundError(t *testing.T) {
	teardown := setupFile(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	jobReq := newJob.ToDto()

	apiError := api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	mockJobFileRepo.EXPECT().FindById(id).Return(nil, apiError)

	err := fileService.addResultToJob(&jobReq, "")

	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("Job with id %v does not exist", id), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_addResultToJob_Returns_NoError(t *testing.T) {
	teardown := setupFile(t)
	defer teardown()
	newJob, _ := realdomain.NewJob("job 1", "url1")
	id := newJob.Id.String()
	jobReq := newJob.ToDto()
	result := "result data"

	mockJobFileRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobFileRepo.EXPECT().SetResult(id, result).Return(nil)

	err := fileService.addResultToJob(&jobReq, result)

	assert.Nil(t, err)
}

func Test_runProbe_Returns_RunError(t *testing.T) {
	config.FfprobePath = probePath
	ctx := context.Background()
	ffArgs := []string{"-loglevel", "warning"}
	cmd := exec.CommandContext(ctx, probePath, ffArgs...)

	data, err := runProbe(cmd)

	assert.EqualValues(t, "", data)
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("error running %v [You have to specify one input file.\r\nUse -h to get full help or, even better, run 'man ffprobe'.\r\n]", probePath), err.Message())
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode())
}

func Test_runProbe_Returns_NoError(t *testing.T) {
	config.FfprobePath = probePath
	ctx := context.Background()
	ffArgs := []string{"-loglevel", "fatal", "-print_format", "json", "-show_format", "-show_streams", "../testmedia/BZ2A3738.MOV"}
	cmd := exec.CommandContext(ctx, probePath, ffArgs...)

	data, err := runProbe(cmd)

	assert.Nil(t, err)
	assert.NotNil(t, data)
}
