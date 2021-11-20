package service

import (
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/services_utils/api_error"
)

type JobService interface {
	GetAllJobs() (*domain.Jobs, api_error.ApiErr)
}

type DefaultJobService struct {
	repo domain.JobRepository
}

func (s DefaultJobService) GetAllJobs() (*domain.Jobs, api_error.ApiErr) {
	return s.repo.FindAll()
}

func NewJobService(repository domain.JobRepository) DefaultJobService {
	return DefaultJobService{repository}
}
