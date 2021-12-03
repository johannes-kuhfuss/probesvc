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
	"github.com/stretchr/testify/assert"
)

var (
	fileCtrl        *gomock.Controller
	mockFileRepo    *domain.MockFileRepository
	fileService     FileService
	jobFileCtrl     *gomock.Controller
	mockJobFileRepo *domain.MockJobRepository
	jobFileService  JobService
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
	_ = jobStatus

	mockJobFileRepo.EXPECT().FindById(id).Return(newJob, nil)
	mockJobFileRepo.EXPECT().SetStatus(id, jobStatus.Status).Return(nil)

	err := fileService.startJob(&jobReq)

	assert.Nil(t, err)
}
