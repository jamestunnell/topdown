// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jamestunnell/topdown/control (interfaces: System)

// Package mock_control is a generated GoMock package.
package mock_control

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSystem is a mock of System interface.
type MockSystem struct {
	ctrl     *gomock.Controller
	recorder *MockSystemMockRecorder
}

// MockSystemMockRecorder is the mock recorder for MockSystem.
type MockSystemMockRecorder struct {
	mock *MockSystem
}

// NewMockSystem creates a new mock instance.
func NewMockSystem(ctrl *gomock.Controller) *MockSystem {
	mock := &MockSystem{ctrl: ctrl}
	mock.recorder = &MockSystemMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSystem) EXPECT() *MockSystemMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockSystem) Add(arg0 string, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0, arg1)
}

// Add indicates an expected call of Add.
func (mr *MockSystemMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockSystem)(nil).Add), arg0, arg1)
}

// Clear mocks base method.
func (m *MockSystem) Clear() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Clear")
}

// Clear indicates an expected call of Clear.
func (mr *MockSystemMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockSystem)(nil).Clear))
}

// Control mocks base method.
func (m *MockSystem) Control(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Control", arg0)
}

// Control indicates an expected call of Control.
func (mr *MockSystemMockRecorder) Control(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Control", reflect.TypeOf((*MockSystem)(nil).Control), arg0)
}

// Remove mocks base method.
func (m *MockSystem) Remove(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Remove", arg0)
}

// Remove indicates an expected call of Remove.
func (mr *MockSystemMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockSystem)(nil).Remove), arg0)
}
