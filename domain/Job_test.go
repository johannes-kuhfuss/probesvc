package domain

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/johannes-kuhfuss/probesvc/dto"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

const (
	validSrcUrl string = "https://server/path/file.ext"
)

func isValidKSUID(id ksuid.KSUID) bool {
	return isValidKSUIDString(id.String())
}

func isValidKSUIDString(id string) bool {
	_, parseErr := ksuid.Parse(id)
	return parseErr == nil
}

func isNowDate(t1, t2 time.Time) bool {
	t1r := t1.Round(1 * time.Minute)
	t2r := t2.Round(1 * time.Minute)
	return t1r == t2r
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
	now := date.GetNowUtc()
	newJob, err := NewJob("", validSrcUrl)
	assert.NotNil(t, newJob)
	assert.Nil(t, err)
	assert.True(t, isValidKSUID(newJob.Id))
	assert.Contains(t, newJob.Name, "new job @")
	assert.True(t, isNowDate(newJob.CreatedAt, now))
	assert.Empty(t, newJob.CreatedBy)
	assert.True(t, isNowDate(newJob.ModifiedAt, now))
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

func Test_JobToDto_Returns_JobDto(t *testing.T) {
	now := date.GetNowUtc()
	newJob, _ := NewJob("my new job", validSrcUrl)
	newJobDto := newJob.ToDto()
	assert.NotNil(t, newJobDto)
	assert.True(t, isValidKSUIDString(newJobDto.Id))
	assert.EqualValues(t, "my new job", newJobDto.Name)
	assert.True(t, isNowDate(newJobDto.CreatedAt, now))
	assert.EqualValues(t, "", newJobDto.CreatedBy)
	assert.True(t, isNowDate(newJobDto.ModifiedAt, now))
	assert.EqualValues(t, "", newJobDto.ModifiedBy)
	assert.EqualValues(t, validSrcUrl, newJobDto.SrcUrl)
	assert.EqualValues(t, JobStatusCreated, newJobDto.Status)
	assert.EqualValues(t, "", newJobDto.ErrorMsg)
	assert.EqualValues(t, "", newJobDto.TechInfo)
}

func Test_ParseStatusRequest_WrongStatus_Returns_Badrequest(t *testing.T) {
	request := dto.JobStatusUpdateRequest{
		Status: "wrong_value",
		ErrMsg: "",
	}
	jobUpd, err := ParseStatusRequest(request)

	assert.Nil(t, jobUpd)
	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("Could not parse status value %v", request.Status), err.Message())
	assert.EqualValues(t, http.StatusBadRequest, err.StatusCode())
}

func Test_ParseStatusRequest_Created_Returns_NoError(t *testing.T) {
	request := dto.JobStatusUpdateRequest{
		Status: "created",
		ErrMsg: "",
	}
	jobUpd, err := ParseStatusRequest(request)

	assert.NotNil(t, jobUpd)
	assert.Nil(t, err)
	assert.EqualValues(t, JobStatusCreated, jobUpd.newStatus)
	assert.EqualValues(t, "", jobUpd.errMsg)
}

func Test_ParseStatusRequest_Queued_Returns_NoError(t *testing.T) {
	request := dto.JobStatusUpdateRequest{
		Status: "queued",
		ErrMsg: "",
	}
	jobUpd, err := ParseStatusRequest(request)

	assert.NotNil(t, jobUpd)
	assert.Nil(t, err)
	assert.EqualValues(t, JobStatusQueued, jobUpd.newStatus)
	assert.EqualValues(t, "", jobUpd.errMsg)
}

func Test_ParseStatusRequest_Running_Returns_NoError(t *testing.T) {
	request := dto.JobStatusUpdateRequest{
		Status: "running",
		ErrMsg: "",
	}
	jobUpd, err := ParseStatusRequest(request)

	assert.NotNil(t, jobUpd)
	assert.Nil(t, err)
	assert.EqualValues(t, JobStatusRunning, jobUpd.newStatus)
	assert.EqualValues(t, "", jobUpd.errMsg)
}

func Test_ParseStatusRequest_Paused_Returns_NoError(t *testing.T) {
	request := dto.JobStatusUpdateRequest{
		Status: "paused",
		ErrMsg: "",
	}
	jobUpd, err := ParseStatusRequest(request)

	assert.NotNil(t, jobUpd)
	assert.Nil(t, err)
	assert.EqualValues(t, JobStatusPaused, jobUpd.newStatus)
	assert.EqualValues(t, "", jobUpd.errMsg)
}

func Test_ParseStatusRequest_Finished_Returns_NoError(t *testing.T) {
	request := dto.JobStatusUpdateRequest{
		Status: "finished",
		ErrMsg: "",
	}
	jobUpd, err := ParseStatusRequest(request)

	assert.NotNil(t, jobUpd)
	assert.Nil(t, err)
	assert.EqualValues(t, JobStatusFinished, jobUpd.newStatus)
	assert.EqualValues(t, "", jobUpd.errMsg)
}

func Test_ParseStatusRequest_Failed_Returns_NoError(t *testing.T) {
	request := dto.JobStatusUpdateRequest{
		Status: "failed",
		ErrMsg: "why-did-it-fail",
	}
	jobUpd, err := ParseStatusRequest(request)

	assert.NotNil(t, jobUpd)
	assert.Nil(t, err)
	assert.EqualValues(t, JobStatusFailed, jobUpd.newStatus)
	assert.EqualValues(t, request.ErrMsg, jobUpd.errMsg)
}
