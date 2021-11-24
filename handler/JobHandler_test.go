package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/dto"
	"github.com/johannes-kuhfuss/probesvc/mocks/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	jh          JobHandlers
	router      *gin.Engine
	mockService *service.MockJobService
	recorder    *httptest.ResponseRecorder
)

func Test_getJobId_NonKsuid_Returns_BadRequestError(t *testing.T) {
	testParam := "wrong_id"
	jobId, err := getJobId(testParam)
	assert.NotNil(t, err)
	assert.EqualValues(t, "", jobId)
	assert.EqualValues(t, "user id should be a ksuid", err.Message())
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
}

func Test_getJobId_WithKsuid_Returns_String(t *testing.T) {
	testParam, _ := ksuid.NewRandom()
	jobId, err := getJobId(testParam.String())
	assert.NotNil(t, jobId)
	assert.Nil(t, err)
	assert.EqualValues(t, testParam.String(), jobId)
}

func setupTest(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = service.NewMockJobService(ctrl)
	jh = JobHandlers{mockService}
	router = gin.Default()
	recorder = httptest.NewRecorder()
	return func() {
		router = nil
		ctrl.Finish()
	}
}

func Test_GetAllJobs_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	dummyJobList := createDummyJobList()
	mockService.EXPECT().GetAllJobs("").Return(&dummyJobList, nil)
	router.GET("/jobs", jh.GetAllJobs)
	request, _ := http.NewRequest(http.MethodGet, "/jobs", nil)

	router.ServeHTTP(recorder, request)
	dummyJobListJson, _ := json.Marshal(dummyJobList)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.EqualValues(t, dummyJobListJson, recorder.Body.String())
}

func createDummyJobList() []dto.JobResponse {
	job1, _ := domain.NewJob("Job 1", "https://server1/path1/file1.ext")
	job2, _ := domain.NewJob("Job 2", "https://server2/path2/file2.ext")
	job1Dto := job1.ToDto()
	job2Dto := job2.ToDto()
	dummyJobList := []dto.JobResponse{}
	dummyJobList = append(dummyJobList, job1Dto)
	dummyJobList = append(dummyJobList, job2Dto)
	return dummyJobList
}

func Test_GetAllJobs_Returns_Error(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("database error")
	mockService.EXPECT().GetAllJobs("").Return(nil, apiError)
	router.GET("/jobs", jh.GetAllJobs)
	request, _ := http.NewRequest(http.MethodGet, "/jobs", nil)

	router.ServeHTTP(recorder, request)
	errorJson, _ := json.Marshal(apiError)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}
