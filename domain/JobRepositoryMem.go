package domain

import (
	"fmt"
	"strings"

	"github.com/johannes-kuhfuss/services_utils/api_error"
)

type JobRepositoryMem struct {
	jobList *[]Job
}

func NewJobRepositoryMem() JobRepositoryMem {
	jList := make([]Job, 0)
	return JobRepositoryMem{&jList}
}

func (csm JobRepositoryMem) FindAll(status string) (*[]Job, api_error.ApiErr) {
	if len(*csm.jobList) == 0 {
		return nil, api_error.NewNotFoundError("no jobs in joblist")
	}
	if strings.TrimSpace(status) == "" {
		return csm.jobList, nil
	} else {
		return filterByStatus(csm.jobList, status)
	}
}

func filterByStatus(jobList *[]Job, status string) (*[]Job, api_error.ApiErr) {
	filteredByStatusList := make([]Job, 0)
	for _, curJob := range *jobList {
		if curJob.Status == JobStatus(status) {
			filteredByStatusList = append(filteredByStatusList, curJob)
		}
	}
	if len(filteredByStatusList) == 0 {
		return nil, api_error.NewNotFoundError(fmt.Sprintf("no jobs with status %v in joblist", status))
	} else {
		return &filteredByStatusList, nil
	}
}

func (csm JobRepositoryMem) FindById(id string) (*Job, api_error.ApiErr) {
	if len(*csm.jobList) == 0 {
		return nil, api_error.NewNotFoundError("no jobs in joblist")
	}
	return filterById(csm.jobList, id)
}

func filterById(jobList *[]Job, id string) (*Job, api_error.ApiErr) {
	for _, curJob := range *jobList {
		if curJob.Id.String() == id {
			return &curJob, nil
		}
	}
	return nil, api_error.NewNotFoundError(fmt.Sprintf("no job with id %v in joblist", id))
}

func (csm JobRepositoryMem) Create(job Job) api_error.ApiErr {
	return api_error.NewBadRequestError("")
}

func (csm JobRepositoryMem) DeleteById(id string) api_error.ApiErr {
	return api_error.NewBadRequestError("")
}

func (csm JobRepositoryMem) GetNextJob() (*Job, api_error.ApiErr) {
	return nil, api_error.NewBadRequestError("")
}
