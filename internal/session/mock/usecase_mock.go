// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/engineerXIII/maiSystemBackend/internal/models"
)

// MockUCSession is a mock of UCSession interface.
type MockUCSession struct {
	ctrl     *gomock.Controller
	recorder *MockUCSessionMockRecorder
}

// MockUCSessionMockRecorder is the mock recorder for MockUCSession.
type MockUCSessionMockRecorder struct {
	mock *MockUCSession
}

// NewMockUCSession creates a new mock instance.
func NewMockUCSession(ctrl *gomock.Controller) *MockUCSession {
	mock := &MockUCSession{ctrl: ctrl}
	mock.recorder = &MockUCSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUCSession) EXPECT() *MockUCSessionMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockUCSession) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, session, expire)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockUCSessionMockRecorder) CreateSession(ctx, session, expire interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockUCSession)(nil).CreateSession), ctx, session, expire)
}

// DeleteByID mocks base method.
func (m *MockUCSession) DeleteByID(ctx context.Context, sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", ctx, sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockUCSessionMockRecorder) DeleteByID(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockUCSession)(nil).DeleteByID), ctx, sessionID)
}

// GetSessionByID mocks base method.
func (m *MockUCSession) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionByID", ctx, sessionID)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionByID indicates an expected call of GetSessionByID.
func (mr *MockUCSessionMockRecorder) GetSessionByID(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionByID", reflect.TypeOf((*MockUCSession)(nil).GetSessionByID), ctx, sessionID)
}
