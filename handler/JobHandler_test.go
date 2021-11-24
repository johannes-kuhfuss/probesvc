package handler

import (
	"net/http"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
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
