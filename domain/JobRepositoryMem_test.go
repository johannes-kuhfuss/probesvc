package domain

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

var (
	jobRepo JobRepositoryMem
)

func setup() func() {
	jobRepo = NewJobRepositoryMem()
	return func() {
		jobRepo.jobList = nil
	}
}

func Test_FindAll_NoJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setup()
	defer teardown()

	jList, err := jobRepo.FindAll("")

	assert.Nil(t, jList)
	assert.NotNil(t, err)
	assert.EqualValues(t, "no jobs in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_FindAll_NoJobsAfterFilter_Returns_NotFoundError(t *testing.T) {
	teardown := setup()
	defer teardown()
	fillJobList()

	jList, err := jobRepo.FindAll("finished")

	assert.Nil(t, jList)
	assert.NotNil(t, err)
	assert.EqualValues(t, "no jobs with status finished in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_FindAll_NoFilter_Returns_NoError(t *testing.T) {
	teardown := setup()
	defer teardown()
	fillJobList()

	jList, err := jobRepo.FindAll("")

	assert.NotNil(t, jList)
	assert.Nil(t, err)
	assert.Equal(t, convertMapToSlice(jobRepo.jobList), jList)
	assert.EqualValues(t, 2, len(*jList))
}

func Test_FindAll_WithFilter_Returns_NoError(t *testing.T) {
	teardown := setup()
	defer teardown()
	fillJobList()

	jList, err := jobRepo.FindAll("running")

	assert.NotNil(t, jList)
	assert.Nil(t, err)
	assert.NotEqual(t, jobRepo.jobList, jList)
	assert.EqualValues(t, 1, len(*jList))
}

func Test_FindById_NoJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setup()
	defer teardown()

	job, err := jobRepo.FindById("")

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, "no jobs in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_FindById_NoJobsAfterFilter_Returns_NotFoundError(t *testing.T) {
	teardown := setup()
	defer teardown()
	fillJobList()
	id := ksuid.New().String()

	job, err := jobRepo.FindById(id)

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("no job with id %v in joblist", id), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_FindById_Returns_NoError(t *testing.T) {
	teardown := setup()
	defer teardown()
	id := fillJobList()

	job, err := jobRepo.FindById(id)

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, id, job.Id.String())
}

func fillJobList() (id string) {
	job1, _ := NewJob("job 1", "url 1")
	job2, _ := NewJob("job 2", "url 2")
	job2.Status = JobStatusRunning
	id1 := job1.Id.String()
	id2 := job2.Id.String()
	jList := make(map[string]Job)
	jList[id1] = *job1
	jList[id2] = *job2
	jobRepo.mu.Lock()
	defer jobRepo.mu.Unlock()
	jobRepo.jobList = jList
	return id1
}

func Test_Create_Returns_NoError(t *testing.T) {
	teardown := setup()
	defer teardown()
	job, _ := NewJob("job 1", "url 1")

	err := jobRepo.Save(*job)

	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(jobRepo.jobList))
}

func Test_DeleteById_NoJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setup()
	defer teardown()

	err := jobRepo.DeleteById("")

	assert.NotNil(t, err)
	assert.EqualValues(t, "no jobs in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_DeleteById_NoJobWithId_Returns_NotFoundError(t *testing.T) {
	teardown := setup()
	defer teardown()
	fillJobList()
	id := ksuid.New().String()

	err := jobRepo.DeleteById(id)

	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("no job with id %v in joblist", id), err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_DeleteById_Returns_NoError(t *testing.T) {
	teardown := setup()
	defer teardown()
	id := fillJobList()

	deletErr := jobRepo.DeleteById(id)
	job, findErr := jobRepo.FindById(id)

	assert.Nil(t, deletErr)
	assert.NotNil(t, findErr)
	assert.Nil(t, job)
	assert.Equal(t, 1, len(jobRepo.jobList))
}

func Test_GetNext_NoJobs_Returns_NotFoundError(t *testing.T) {
	teardown := setup()
	defer teardown()

	job, err := jobRepo.GetNext()

	assert.Nil(t, job)
	assert.NotNil(t, err)
	assert.EqualValues(t, "no jobs in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_GetNext_Returns_NoError(t *testing.T) {
	teardown := setup()
	defer teardown()
	createdId := fillJobList()

	job, err := jobRepo.GetNext()

	assert.NotNil(t, job)
	assert.Nil(t, err)
	assert.EqualValues(t, createdId, job.Id.String())
}

func Test_SetStatus_NoJob_Returns_NotFoundError(t *testing.T) {
	teardown := setup()
	defer teardown()
	newStatus := JobStatusUpdate{}
	err := jobRepo.SetStatus("", newStatus)

	assert.NotNil(t, err)
	assert.EqualValues(t, "no jobs in joblist", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_SetStatus_Returns_NoError(t *testing.T) {
	teardown := setup()
	defer teardown()
	id := fillJobList()
	newStatus := JobStatusUpdate{
		newStatus: JobStatusFailed,
		errMsg:    "why-did-I-fail",
	}
	err := jobRepo.SetStatus(id, newStatus)
	job, _ := jobRepo.FindById(id)

	assert.Nil(t, err)
	assert.EqualValues(t, newStatus.newStatus, job.Status)
	assert.EqualValues(t, newStatus.errMsg, job.ErrorMsg)
}
