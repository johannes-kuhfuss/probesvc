package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/johannes-kuhfuss/probesvc/config"
	"github.com/johannes-kuhfuss/probesvc/domain"
	"github.com/johannes-kuhfuss/probesvc/dto"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type FileService interface {
	Run()
}

type DefaultFileService struct {
	repo   domain.FileRepository
	jobSrv JobService
}

var (
	binPath   string                     = "./service/ffprobe.exe"
	jobStatus dto.JobStatusUpdateRequest = dto.JobStatusUpdateRequest{}
)

func NewFileService(repository domain.FileRepository, jobSrv JobService) DefaultFileService {
	return DefaultFileService{repository, jobSrv}
}

func (s DefaultFileService) Run() {

	for !config.Shutdown {
		job, err := s.jobSrv.GetNextJob()
		if err != nil {
			logger.Debug(err.Message())
			time.Sleep(time.Second * time.Duration(config.NoJobWaitTime))
		} else {
			jobStatus.Status = "running"
			s.jobSrv.SetStatus(job.Id, jobStatus)
			result, err := s.analyzeFile(job.SrcUrl)
			if err != nil {
				s.failJob(err, job.Id)
			} else {
				err := s.addResultToJob(job.Id, result)
				if err != nil {
					s.failJob(err, job.Id)
				} else {
					s.finishJob(job.Id, job.SrcUrl)
				}
			}
		}
	}
}

func (s DefaultFileService) failJob(err api_error.ApiErr, id string) {
	logger.Error("Error while analyzing file", err)
	jobStatus.Status = "failed"
	jobStatus.ErrMsg = "Error while analyzing file"
	s.jobSrv.SetStatus(id, jobStatus)
}

func (s DefaultFileService) finishJob(id string, srcUrl string) {
	logger.Info(fmt.Sprintf("Finished data extraction for Job ID %v with Source %v", id, srcUrl))
	jobStatus.Status = "failed"
	jobStatus.ErrMsg = "Error while analyzing file"
	s.jobSrv.SetStatus(id, jobStatus)
}

func (s DefaultFileService) addResultToJob(id string, result string) api_error.ApiErr {
	err := s.jobSrv.SetResult(id, result)
	if err != nil {
		return err
	}
	return nil
}

func (s DefaultFileService) analyzeFile(srcUrl string) (string, api_error.ApiErr) {
	ctx := context.Background()
	reader, err := s.getAzureReader(srcUrl)
	if err != nil {
		return "", api_error.NewInternalServerError("could not connect to storage", err)
	}

	ffArgs := []string{"-loglevel", "fatal", "-print_format", "json", "-show_format", "-show_streams", "-"}
	cmd := exec.CommandContext(ctx, binPath, ffArgs...)
	cmd.Stdin = *reader

	result, runErr := runProbe(cmd)
	if runErr != nil {
		return "", api_error.NewInternalServerError("could not extract metadata from file", err)
	}
	return result, nil
}

func (s DefaultFileService) getAzureReader(srcUrl string) (*io.ReadCloser, api_error.ApiErr) {
	url, _ := url.Parse(srcUrl)
	containerName := strings.TrimLeft(filepath.Dir(url.Path), string(os.PathSeparator))
	fileName := filepath.Base(srcUrl)
	ctx := context.Background()
	container := s.repo.GetClient().NewContainerClient(containerName)
	blockBlob := container.NewBlobClient(fileName)
	get, err := blockBlob.Download(ctx, nil)
	if err != nil {
		logger.Error("Cannot access file on storage account", err)
		return nil, api_error.NewBadRequestError("Cannot access file on storage account")
	}
	reader := get.Body(azblob.RetryReaderOptions{})
	return &reader, nil
}

func runProbe(cmd *exec.Cmd) (data string, err error) {
	var outputBuf bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &outputBuf
	cmd.Stderr = &stdErr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running %s [%s] %w", binPath, stdErr.String(), err)
	}

	if stdErr.Len() > 0 {
		return "", fmt.Errorf("ffprobe error: %s", stdErr.String())
	}

	data = outputBuf.String()

	return data, nil
}
