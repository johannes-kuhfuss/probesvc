package domain

import "github.com/johannes-kuhfuss/services_utils/api_error"

type JobRepositoryMem struct {
	jobList []Job
}

func NewJobRepositoryMem() JobRepositoryMem {
	jList := make([]Job, 0)
	return JobRepositoryMem{jList}
}

func (csm JobRepositoryMem) FindAll(status string) (*[]Job, api_error.ApiErr) {
	return nil, api_error.NewBadRequestError("")
}

func (csm JobRepositoryMem) FindById(id string) (*Job, api_error.ApiErr) {
	return nil, api_error.NewBadRequestError("")
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
