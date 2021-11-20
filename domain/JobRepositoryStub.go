package domain

import "github.com/johannes-kuhfuss/services_utils/api_error"

type JobRepositoryStub struct {
	jobsList Jobs
}

func (jrs JobRepositoryStub) FindAll() (*Jobs, api_error.ApiErr) {
	return &jrs.jobsList, nil
}

func NewJobRepositoryStub() JobRepositoryStub {
	job1, _ := NewJob("Job 1", "https://server1/path1/file1.ext")
	job2, _ := NewJob("Job 2", "https://server2/path2/file2.ext")
	jList := Jobs{}
	jList = append(jList, *job1)
	jList = append(jList, *job2)
	return JobRepositoryStub{jList}
}
