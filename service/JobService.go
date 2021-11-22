package service

import (
	"fmt"

	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
)

type JobService interface {
	GetAllJobs(string) (*[]dto.JobResponse, api_error.ApiErr)
	GetJobById(string) (*dto.JobResponse, api_error.ApiErr)
	CreateJob(dto.NewJobRequest) (*dto.JobResponse, api_error.ApiErr)
	DeleteJobById(string) api_error.ApiErr
	GetNextJob() (*dto.JobResponse, api_error.ApiErr)
}

type DefaultJobService struct {
	repo domain.JobRepository
}

func NewJobService(repository domain.JobRepository) DefaultJobService {
	return DefaultJobService{repository}
}

func (s DefaultJobService) GetAllJobs(status string) (*[]dto.JobResponse, api_error.ApiErr) {
	jobs, err := s.repo.FindAll(status)
	if err != nil {
		return nil, err
	}
	response := make([]dto.JobResponse, 0)
	for _, job := range *jobs {
		response = append(response, job.ToDto())
	}
	return &response, nil
}

func (s DefaultJobService) GetJobById(id string) (*dto.JobResponse, api_error.ApiErr) {
	job, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	response := job.ToDto()
	return &response, nil
}

func (s DefaultJobService) CreateJob(jobreq dto.NewJobRequest) (*dto.JobResponse, api_error.ApiErr) {
	newJob, err := domain.NewJob(jobreq.Name, jobreq.SrcUrl)
	if err != nil {
		return nil, err
	}
	err = s.repo.Create(*newJob)
	if err != nil {
		return nil, err
	}
	response := newJob.ToDto()
	return &response, nil
}

func (s DefaultJobService) DeleteJobById(id string) api_error.ApiErr {
	_, err := s.GetJobById(id)
	if err != nil {
		return api_error.NewNotFoundError(fmt.Sprintf("Job with id %v does not exist", id))
	}
	err = s.repo.DeleteById(id)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultJobService) GetNextJob() (*dto.JobResponse, api_error.ApiErr) {
	job, err := s.repo.GetNextJob()
	if err != nil {
		return nil, err
	}
	response := job.ToDto()
	return &response, nil
}
