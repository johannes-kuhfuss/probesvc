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

func Test_FindAll_Returns_Full_NoError(t *testing.T) {
	teardown := setup()
	defer teardown()
	fillJobList()

	jList, err := jobRepo.FindAll("")

	assert.NotNil(t, jList)
	assert.Nil(t, err)
	assert.Equal(t, jobRepo.jobList, jList)
	assert.EqualValues(t, 2, len(*jList))
}

func Test_FindAll_Returns_Partial_NoError(t *testing.T) {
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

func fillJobList() (id1 string) {
	job1, _ := NewJob("job 1", "url 1")
	job2, _ := NewJob("job 2", "url 2")
	job2.Status = JobStatusRunning
	jList := make([]Job, 0)
	jList = append(jList, *job1)
	jList = append(jList, *job2)
	jobRepo.jobList = &jList
	return job1.Id.String()
}
