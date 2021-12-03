package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/johannes-kuhfuss/probesvc/mocks/domain"
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

func Test_startJob(t *testing.T) {
	teardown := setupFile(t)
	defer teardown()

}
