package service

import (
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/services_utils/api_error"
)

type JobService interface {
	GetAllJobs() (*domain.Jobs, api_error.ApiErr)
	GetJobById(string) (*domain.Job, api_error.ApiErr)
}

type DefaultJobService struct {
	repo domain.JobRepository
}

func NewJobService(repository domain.JobRepository) DefaultJobService {
	return DefaultJobService{repository}
}

func (s DefaultJobService) GetAllJobs() (*domain.Jobs, api_error.ApiErr) {
	return s.repo.FindAll()
}

func (s DefaultJobService) GetJobById(id string) (*domain.Job, api_error.ApiErr) {
	return s.repo.FindById(id)
}
