package service

import "github.com/johannes-kuhfuss/probesvc/domain"

type FileService interface {
	Run()
}

type DefaultFileService struct {
	repo domain.FileRepository
}

func NewFileService(repository domain.FileRepository) DefaultFileService {
	return DefaultFileService{repository}
}

func (s DefaultJobService) Run() {
}
