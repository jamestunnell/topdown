// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jamestunnell/topdown/control (interfaces: Controllable)

// Package mock_control is a generated GoMock package.
package mock_control

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	input "github.com/jamestunnell/topdown/input"
)

// MockControllable is a mock of Controllable interface.
type MockControllable struct {
	ctrl     *gomock.Controller
	recorder *MockControllableMockRecorder
}

// MockControllableMockRecorder is the mock recorder for MockControllable.
type MockControllableMockRecorder struct {
	mock *MockControllable
}

// NewMockControllable creates a new mock instance.
func NewMockControllable(ctrl *gomock.Controller) *MockControllable {
	mock := &MockControllable{ctrl: ctrl}
	mock.recorder = &MockControllableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockControllable) EXPECT() *MockControllableMockRecorder {
	return m.recorder
}

// Control mocks base method.
func (m *MockControllable) Control(arg0 float64, arg1 input.Manager) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Control", arg0, arg1)
}

// Control indicates an expected call of Control.
func (mr *MockControllableMockRecorder) Control(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Control", reflect.TypeOf((*MockControllable)(nil).Control), arg0, arg1)
}

// WatchKeys mocks base method.
func (m *MockControllable) WatchKeys() []ebiten.Key {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchKeys")
	ret0, _ := ret[0].([]ebiten.Key)
	return ret0
}

// WatchKeys indicates an expected call of WatchKeys.
func (mr *MockControllableMockRecorder) WatchKeys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchKeys", reflect.TypeOf((*MockControllable)(nil).WatchKeys))
}
