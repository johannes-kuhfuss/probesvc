package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/dto"
	"github.com/johannes-kuhfuss/probesvc/mocks/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/date"
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
	assert.EqualValues(t, "User id should be a ksuid", err.Message())
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

func Test_GetJobById_Returns_InvalidIdError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("User id should be a ksuid")
	errorJson, _ := json.Marshal(apiError)
	router.GET("/jobs/:job_id", jh.GetJobById)
	request, _ := http.NewRequest(http.MethodGet, "/jobs/not_a_ksuid", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

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
	mockService.EXPECT().GetJobById(id.String()).Return(&newJob, nil)
	router.GET("/jobs/:job_id", jh.GetJobById)
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jobs/%v", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.EqualValues(t, bodyJson, recorder.Body.String())
}

func Test_CreateJob_Returns_InvalidJsonError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("invalid json body")
	errorJson, _ := json.Marshal(apiError)
	router.POST("/jobs", jh.CreateJob)
	request, _ := http.NewRequest(http.MethodPost, "/jobs", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_CreateJob_Returns_ServiceError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewInternalServerError("database error", nil)
	errorJson, _ := json.Marshal(apiError)
	jobReq := dto.NewJobRequest{
		Name:   "my new job",
		SrcUrl: "http://server/path/file.ext",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	mockService.EXPECT().CreateJob(jobReq).Return(nil, apiError)
	router.POST("/jobs", jh.CreateJob)
	request, _ := http.NewRequest(http.MethodPost, "/jobs", strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusInternalServerError, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_CreateJob_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	jobReq := dto.NewJobRequest{
		Name:   "my new job",
		SrcUrl: "http://server/path/file.ext",
	}
	jobReqJson, _ := json.Marshal(jobReq)
	jobResp := dto.JobResponse{
		Id:         ksuid.New().String(),
		Name:       jobReq.Name,
		CreatedAt:  date.GetNowUtc(),
		CreatedBy:  "",
		ModifiedAt: date.GetNowUtc(),
		ModifiedBy: "",
		SrcUrl:     jobReq.SrcUrl,
		Status:     "created",
		ErrorMsg:   "",
		TechInfo:   "",
	}
	bodyJson, _ := json.Marshal(jobResp)
	mockService.EXPECT().CreateJob(jobReq).Return(&jobResp, nil)
	router.POST("/jobs", jh.CreateJob)
	request, _ := http.NewRequest(http.MethodPost, "/jobs", strings.NewReader(string(jobReqJson)))

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusCreated, recorder.Code)
	assert.EqualValues(t, bodyJson, recorder.Body.String())
}

func Test_DeleteJobById_Returns_InvalidIdError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewBadRequestError("User id should be a ksuid")
	errorJson, _ := json.Marshal(apiError)
	router.DELETE("/jobs/:job_id", jh.DeleteJobById)
	request, _ := http.NewRequest(http.MethodDelete, "/jobs/not_a_ksuid", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusBadRequest, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_DeleteJobById_Returns_NotFoundError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	apiError := api_error.NewNotFoundError(fmt.Sprintf("job with id %v not found", id))
	errorJson, _ := json.Marshal(apiError)
	mockService.EXPECT().DeleteJobById(id.String()).Return(apiError)
	router.DELETE("/jobs/:job_id", jh.DeleteJobById)
	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/jobs/%v", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusNotFound, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_DeleteJobById_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	id := ksuid.New()
	mockService.EXPECT().DeleteJobById(id.String()).Return(nil)
	router.DELETE("/jobs/:job_id", jh.DeleteJobById)
	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/jobs/%v", id), nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
}

func Test_GetNextJob_Returns_NotFoundError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	apiError := api_error.NewNotFoundError("No next job found")
	errorJson, _ := json.Marshal(apiError)
	mockService.EXPECT().GetNextJob().Return(nil, apiError)
	router.GET("/jobs/next", jh.GetNextJob)
	request, _ := http.NewRequest(http.MethodGet, "/jobs/next", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusNotFound, recorder.Code)
	assert.EqualValues(t, errorJson, recorder.Body.String())
}

func Test_GetNextJob_Returns_NoError(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()
	jobResp := dto.JobResponse{
		Id:         ksuid.New().String(),
		Name:       "job 1",
		CreatedAt:  date.GetNowUtc(),
		CreatedBy:  "",
		ModifiedAt: date.GetNowUtc(),
		ModifiedBy: "",
		SrcUrl:     "http://server/path/file.ext",
		Status:     "created",
		ErrorMsg:   "",
		TechInfo:   "",
	}
	bodyJson, _ := json.Marshal(jobResp)
	mockService.EXPECT().GetNextJob().Return(&jobResp, nil)
	router.GET("/jobs/next", jh.GetNextJob)
	request, _ := http.NewRequest(http.MethodGet, "/jobs/next", nil)

	router.ServeHTTP(recorder, request)

	assert.EqualValues(t, http.StatusOK, recorder.Code)
	assert.EqualValues(t, bodyJson, recorder.Body.String())
}
