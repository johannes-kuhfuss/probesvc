package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	testParam := ksuid.New()
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
	dummyJobListJson, _ := json.Marshal(dummyJobList)
	mockService.EXPECT().GetAllJobs("").Return(&dummyJobList, nil)
	router.GET("/jobs", jh.GetAllJobs)
	request, _ := http.NewRequest(http.MethodGet, "/jobs", nil)

	router.ServeHTTP(recorder, request)

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
	errorJson, _ := json.Marshal(apiError)
	mockService.EXPECT().GetAllJobs("").Return(nil, apiError)
	router.GET("/jobs", jh.GetAllJobs)
	request, _ := http.NewRequest(http.MethodGet, "/jobs", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

/*
func Test_GetJobById_Returns_InvalidIdError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("User id should be a ksuid")
	errorJson, _ := json.Marshal(apiError)
	router.GET("/jobs/:job_id", jh.GetJobById)
	request, _ := http.NewRequest(http.MethodGet, "/jobs/:jod_id", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}
*/

func Test_GetJobById_Returns_NotFoundError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewNotFoundError(fmt.Sprintf("job with id %v not found", id))
	errorJson, _ := json.Marshal(apiError)
	mockService.EXPECT().GetJobById(gomock.Eq(id.String())).Return(nil, apiError)
	router.GET("/jobs/:job_id", jh.GetJobById)
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jobs/%v", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusNotFound, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_GetJobById_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	newJob := dto.JobResponse{
		Id:         id.String(),
		Name:       "",
		CreatedAt:  time.Time{},
		CreatedBy:  "",
		ModifiedAt: time.Time{},
		ModifiedBy: "",
		SrcUrl:     "",
		Status:     "",
		ErrorMsg:   "",
		TechInfo:   "",
	}
	bodyJson, _ := json.Marshal(newJob)
	mockService.EXPECT().GetJobById(gomock.Eq(id.String())).Return(&newJob, nil)
	router.GET("/jobs/:job_id", jh.GetJobById)
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jobs/%v", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.EqualValues(t, bodyJson, recorder.Body.String())
}

func Test_CreateNewJob_Returns_InvalidJsonError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("invalid json body")
	errorJson, _ := json.Marshal(apiError)
	router.POST("/jobs", jh.CreateNewJob)
	request, _ := http.NewRequest(http.MethodPost, "/jobs", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}
