// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jamestunnell/topdown/engine (interfaces: Engine)

// Package mock_engine is a generated GoMock package.
package mock_engine

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resource "github.com/jamestunnell/topdown/resource"
)

// MockEngine is a mock of Engine interface.
type MockEngine struct {
	ctrl     *gomock.Controller
	recorder *MockEngineMockRecorder
}

// MockEngineMockRecorder is the mock recorder for MockEngine.
type MockEngineMockRecorder struct {
	mock *MockEngine
}

// NewMockEngine creates a new mock instance.
func NewMockEngine(ctrl *gomock.Controller) *MockEngine {
	mock := &MockEngine{ctrl: ctrl}
	mock.recorder = &MockEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEngine) EXPECT() *MockEngineMockRecorder {
	return m.recorder
}

// Initialize mocks base method.
func (m *MockEngine) Initialize() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize")
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockEngineMockRecorder) Initialize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockEngine)(nil).Initialize))
}

// ResourceManager mocks base method.
func (m *MockEngine) ResourceManager() resource.Manager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResourceManager")
	ret0, _ := ret[0].(resource.Manager)
	return ret0
}

// ResourceManager indicates an expected call of ResourceManager.
func (mr *MockEngineMockRecorder) ResourceManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResourceManager", reflect.TypeOf((*MockEngine)(nil).ResourceManager))
}