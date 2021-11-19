package domain

import (
	"net/http"
	"testing"
	"time"

	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

const (
	validSrcUrl string = "https://server/path/file.ext"
)

func isValidKSUID(id string) bool {
	_, parseErr := ksuid.Parse(id)
	return parseErr == nil
}

func isValidDate(dateStr string) bool {
	_, parseErr := time.Parse(date.ApiDateLayout, dateStr)
	return parseErr == nil
}

func TestConstants(t *testing.T) {
	assert.EqualValues(t, JobStatusCreated, "created")
	assert.EqualValues(t, JobStatusQueued, "queued")
	assert.EqualValues(t, JobStatusRunning, "running")
	assert.EqualValues(t, JobStatusPaused, "paused")
	assert.EqualValues(t, JobStatusFinished, "finished")
	assert.EqualValues(t, JobStatusFailed, "failed")
}

func Test_NewJob_NoSrUrl_ReturnsBadRequestErr(t *testing.T) {
	newJob, err := NewJob(" ", " ")
	assert.Nil(t, newJob)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
	assert.EqualValues(t, "Job must have a source URL", err.Message())
}

func Test_NewJob_NoName_Returns_NewJob(t *testing.T) {
	newJob, err := NewJob("", validSrcUrl)
	assert.NotNil(t, newJob)
	assert.Nil(t, err)
	assert.True(t, isValidKSUID(newJob.Id))
	assert.Contains(t, newJob.Name, "new job @")
	assert.True(t, isValidDate(newJob.CreatedAt))
	assert.Empty(t, newJob.CreatedBy)
	assert.Empty(t, newJob.ModifiedAt)
	assert.Empty(t, newJob.ModifiedBy)
	assert.EqualValues(t, validSrcUrl, newJob.SrcUrl)
	assert.EqualValues(t, JobStatusCreated, newJob.Status)
	assert.Empty(t, newJob.ErrorMsg)
	assert.Empty(t, newJob.TechInfo)
}

func Test_NewJob_WithName_Returns_NewJob(t *testing.T) {
	newJob, err := NewJob("my new job", validSrcUrl)
	assert.NotNil(t, newJob)
	assert.Nil(t, err)
	assert.EqualValues(t, "my new job", newJob.Name)
}
