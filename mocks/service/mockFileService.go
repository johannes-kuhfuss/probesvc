// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/johannes-kuhfuss/probesvc/service (interfaces: FileService)

// Package service is a generated GoMock package.
package service

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/johannes-kuhfuss/probesvc/dto"
	api_error "github.com/johannes-kuhfuss/services_utils/api_error"
)

// MockFileService is a mock of FileService interface.
type MockFileService struct {
	ctrl     *gomock.Controller
	recorder *MockFileServiceMockRecorder
}

// MockFileServiceMockRecorder is the mock recorder for MockFileService.
type MockFileServiceMockRecorder struct {
	mock *MockFileService
}

// NewMockFileService creates a new mock instance.
func NewMockFileService(ctrl *gomock.Controller) *MockFileService {
	mock := &MockFileService{ctrl: ctrl}
	mock.recorder = &MockFileServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileService) EXPECT() *MockFileServiceMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockFileService) Run() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run")
}

// Run indicates an expected call of Run.
func (mr *MockFileServiceMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockFileService)(nil).Run))
}

// addResultToJob mocks base method.
func (m *MockFileService) addResultToJob(arg0 *dto.JobResponse, arg1 string) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "addResultToJob", arg0, arg1)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// addResultToJob indicates an expected call of addResultToJob.
func (mr *MockFileServiceMockRecorder) addResultToJob(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "addResultToJob", reflect.TypeOf((*MockFileService)(nil).addResultToJob), arg0, arg1)
}

// failJob mocks base method.
func (m *MockFileService) failJob(arg0 *dto.JobResponse, arg1 api_error.ApiErr) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "failJob", arg0, arg1)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// failJob indicates an expected call of failJob.
func (mr *MockFileServiceMockRecorder) failJob(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "failJob", reflect.TypeOf((*MockFileService)(nil).failJob), arg0, arg1)
}

// finishJob mocks base method.
func (m *MockFileService) finishJob(arg0 *dto.JobResponse) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "finishJob", arg0)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// finishJob indicates an expected call of finishJob.
func (mr *MockFileServiceMockRecorder) finishJob(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "finishJob", reflect.TypeOf((*MockFileService)(nil).finishJob), arg0)
}

// getAzureReader mocks base method.
func (m *MockFileService) getAzureReader(arg0 string) (*io.ReadCloser, api_error.ApiErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getAzureReader", arg0)
	ret0, _ := ret[0].(*io.ReadCloser)
	ret1, _ := ret[1].(api_error.ApiErr)
	return ret0, ret1
}

// getAzureReader indicates an expected call of getAzureReader.
func (mr *MockFileServiceMockRecorder) getAzureReader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getAzureReader", reflect.TypeOf((*MockFileService)(nil).getAzureReader), arg0)
}

// startJob mocks base method.
func (m *MockFileService) startJob(arg0 *dto.JobResponse) api_error.ApiErr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "startJob", arg0)
	ret0, _ := ret[0].(api_error.ApiErr)
	return ret0
}

// startJob indicates an expected call of startJob.
func (mr *MockFileServiceMockRecorder) startJob(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "startJob", reflect.TypeOf((*MockFileService)(nil).startJob), arg0)
}
