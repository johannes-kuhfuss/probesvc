package domain

import (
	"fmt"
	"strings"

	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/date"
	"github.com/segmentio/ksuid"
)

type JobStatus string

const (
	JobStatusCreated  JobStatus = "created"
	JobStatusQueued   JobStatus = "queued"
	JobStatusRunning  JobStatus = "running"
	JobStatusPaused   JobStatus = "paused"
	JobStatusFinished JobStatus = "finished"
	JobStatusFailed   JobStatus = "failed"
)

type Job struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  string    `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	ModifiedAt string    `json:"modified_at"`
	ModifiedBy string    `json:"modified_by"`
	SrcUrl     string    `json:"src_url"`
	Status     JobStatus `json:"status"`
	ErrorMsg   string    `json:"error_msg"`
	TechInfo   string    `json:"tech_info"`
}

func NewJob(name string, srcurl string) (*Job, api_error.ApiErr) {
	if strings.TrimSpace(srcurl) == "" {
		return nil, api_error.NewBadRequestError("Job must have a source URL")
	}

	var jobName string
	if strings.TrimSpace(name) == "" {
		jobName = fmt.Sprintf("new job @ %s", date.GetNowUtcString())
	} else {
		jobName = name
	}

	return &Job{
		Id:         ksuid.New().String(),
		Name:       jobName,
		CreatedAt:  date.GetNowUtcString(),
		CreatedBy:  "",
		ModifiedAt: "",
		ModifiedBy: "",
		SrcUrl:     srcurl,
		Status:     JobStatusCreated,
		ErrorMsg:   "",
		TechInfo:   "",
	}, nil
}
