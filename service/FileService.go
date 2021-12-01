package service

import (
	"time"

	"github.com/johannes-kuhfuss/probesvc/config"
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/dto"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type FileService interface {
	Run()
}

type DefaultFileService struct {
	repo   domain.FileRepository
	jobSrv DefaultJobService
}

func NewFileService(repository domain.FileRepository, jobSrv DefaultJobService) DefaultFileService {
	return DefaultFileService{repository, jobSrv}
}

func (s DefaultFileService) Run() {
	for !config.Shutdown {
		job, err := s.jobSrv.GetNextJob()
		if err != nil {
			logger.Debug(err.Message())
			time.Sleep(time.Second * time.Duration(config.NoJobWaitTime))
		} else {
			jobStatus := dto.JobStatusUpdateRequest{
				Status: "running",
				ErrMsg: "",
			}
			s.jobSrv.SetStatus(job.Id, jobStatus)
			time.Sleep(time.Second * time.Duration(config.NoJobWaitTime))
			jobStatus.Status = "failed"
			jobStatus.ErrMsg = "processing not implemented"
			s.jobSrv.SetStatus(job.Id, jobStatus)
		}
	}
}
