package domain

import (
	"fmt"
	"strings"
	"time"

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
	Id         ksuid.KSUID `db:"job_id"`
	Name       string      `db:"name"`
	CreatedAt  time.Time   `db:"created_at"`
	CreatedBy  string      `db:"created_by"`
	ModifiedAt time.Time   `db:"modified_at"`
	ModifiedBy string      `db:"modified_by"`
	SrcUrl     string      `db:"src_url"`
	Status     JobStatus   `db:"status"`
	ErrorMsg   string      `db:"error_msg"`
	TechInfo   string      `db:"tech_info"`
}

type Jobs []Job

type JobRepository interface {
	FindAll(string) (*Jobs, api_error.ApiErr)
	FindById(string) (*Job, api_error.ApiErr)
}

func createJobName(name string) string {
	var jobName string
	if strings.TrimSpace(name) == "" {
		newDate, _ := date.GetNowLocalString("")
		jobName = fmt.Sprintf("new job @ %s", *newDate)
	} else {
		jobName = name
	}
	return jobName
}

func NewJob(name string, srcurl string) (*Job, api_error.ApiErr) {
	if strings.TrimSpace(srcurl) == "" {
		return nil, api_error.NewBadRequestError("Job must have a source URL")
	}

	return &Job{
		Id:         ksuid.New(),
		Name:       createJobName(name),
		CreatedAt:  date.GetNowUtc(),
		CreatedBy:  "",
		ModifiedAt: time.Time{},
		ModifiedBy: "",
		SrcUrl:     srcurl,
		Status:     JobStatusCreated,
		ErrorMsg:   "",
		TechInfo:   "",
	}, nil
}
