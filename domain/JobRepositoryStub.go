package domain

import "github.com/johannes-kuhfuss/services_utils/api_error"

type JobRepositoryStub struct {
	jobList []Job
}

func (jrs JobRepositoryStub) FindAll() (*[]Job, api_error.ApiErr) {
	return &jrs.jobList, nil
}

func NewJobRepositoryStub() JobRepositoryStub {
	job1, _ := NewJob("Job 1", "https://server1/path1/file1.ext")
	job2, _ := NewJob("Job 2", "https://server2/path2/file2.ext")
	jList := []Job{}
	jList = append(jList, *job1)
	jList = append(jList, *job2)
	return JobRepositoryStub{jList}
}
