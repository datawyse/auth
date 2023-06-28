// Code generated by MockGen. DO NOT EDIT.
// Source: core/ports/auth.go

// Package mocks is a generated GoMock package.
package mocks

import (
	domain "auth/core/domain"
	http "auth/core/domain/http"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockAuthService) Login(username, password string) (*domain.AuthToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", username, password)
	ret0, _ := ret[0].(*domain.AuthToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthServiceMockRecorder) Login(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthService)(nil).Login), username, password)
}

// RefreshToken mocks base method.
func (m *MockAuthService) RefreshToken(refreshToken string) (*domain.AuthToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", refreshToken)
	ret0, _ := ret[0].(*domain.AuthToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockAuthServiceMockRecorder) RefreshToken(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockAuthService)(nil).RefreshToken), refreshToken)
}

// Signup mocks base method.
func (m *MockAuthService) Signup(ctx context.Context, input *http.SignupInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signup", ctx, input)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Signup indicates an expected call of Signup.
func (mr *MockAuthServiceMockRecorder) Signup(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signup", reflect.TypeOf((*MockAuthService)(nil).Signup), ctx, input)
}

// UpdateUser mocks base method.
func (m *MockAuthService) UpdateUser(ctx context.Context, input *http.UpdateUserInput) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, input)
	ret0, _ := ret[0].(*domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockAuthServiceMockRecorder) UpdateUser(ctx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockAuthService)(nil).UpdateUser), ctx, input)
}

// User mocks base method.
func (m *MockAuthService) User(userId string) (*domain.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "User", userId)
	ret0, _ := ret[0].(*domain.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// User indicates an expected call of User.
func (mr *MockAuthServiceMockRecorder) User(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "User", reflect.TypeOf((*MockAuthService)(nil).User), userId)
}
