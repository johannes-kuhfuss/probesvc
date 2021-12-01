package service

import (
	"time"

	"github.com/johannes-kuhfuss/probesvc/config"
	"github.com/johannes-kuhfuss/probesvc/domain"
)

type FileService interface {
	Run()
}

type DefaultFileService struct {
	repo domain.FileRepository
}

func NewFileService(repository domain.FileRepository) DefaultFileService {
	return DefaultFileService{repository}
}

func (s DefaultFileService) Run() {
	for !config.Shutdown {
		time.Sleep(time.Second * time.Duration(config.NoJobWaitTime))
	}
}
