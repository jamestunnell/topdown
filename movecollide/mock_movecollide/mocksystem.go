// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jamestunnell/topdown/movecollide (interfaces: System)

// Package mock_movecollide is a generated GoMock package.
package mock_movecollide

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	movecollide "github.com/jamestunnell/topdown/movecollide"
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

// MoveCollide mocks base method.
func (m *MockSystem) MoveCollide(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MoveCollide", arg0)
}

// MoveCollide indicates an expected call of MoveCollide.
func (mr *MockSystemMockRecorder) MoveCollide(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MoveCollide", reflect.TypeOf((*MockSystem)(nil).MoveCollide), arg0)
}

// Raycast mocks base method.
func (m *MockSystem) Raycast(arg0 *movecollide.Ray) (*movecollide.RayHit, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Raycast", arg0)
	ret0, _ := ret[0].(*movecollide.RayHit)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Raycast indicates an expected call of Raycast.
func (mr *MockSystemMockRecorder) Raycast(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Raycast", reflect.TypeOf((*MockSystem)(nil).Raycast), arg0)
}
