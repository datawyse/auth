// Code generated by MockGen. DO NOT EDIT.
// Source: core/ports/uuid.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockUUIDService is a mock of UUIDService interface.
type MockUUIDService struct {
	ctrl     *gomock.Controller
	recorder *MockUUIDServiceMockRecorder
}

// MockUUIDServiceMockRecorder is the mock recorder for MockUUIDService.
type MockUUIDServiceMockRecorder struct {
	mock *MockUUIDService
}

// NewMockUUIDService creates a new mock instance.
func NewMockUUIDService(ctrl *gomock.Controller) *MockUUIDService {
	mock := &MockUUIDService{ctrl: ctrl}
	mock.recorder = &MockUUIDServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUUIDService) EXPECT() *MockUUIDServiceMockRecorder {
	return m.recorder
}

// FromString mocks base method.
func (m *MockUUIDService) FromString(id string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FromString", id)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FromString indicates an expected call of FromString.
func (mr *MockUUIDServiceMockRecorder) FromString(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FromString", reflect.TypeOf((*MockUUIDService)(nil).FromString), id)
}

// IsValidUUID mocks base method.
func (m *MockUUIDService) IsValidUUID(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValidUUID", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsValidUUID indicates an expected call of IsValidUUID.
func (mr *MockUUIDServiceMockRecorder) IsValidUUID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValidUUID", reflect.TypeOf((*MockUUIDService)(nil).IsValidUUID), id)
}

// NewUUID mocks base method.
func (m *MockUUIDService) NewUUID() (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewUUID")
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewUUID indicates an expected call of NewUUID.
func (mr *MockUUIDServiceMockRecorder) NewUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewUUID", reflect.TypeOf((*MockUUIDService)(nil).NewUUID))
}
