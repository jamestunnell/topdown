// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jamestunnell/topdown/resource (interfaces: Resource)

// Package mock_resource is a generated GoMock package.
package mock_resource

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resource "github.com/jamestunnell/topdown/resource"
)

// MockResource is a mock of Resource interface.
type MockResource struct {
	ctrl     *gomock.Controller
	recorder *MockResourceMockRecorder
}

// MockResourceMockRecorder is the mock recorder for MockResource.
type MockResourceMockRecorder struct {
	mock *MockResource
}

// NewMockResource creates a new mock instance.
func NewMockResource(ctrl *gomock.Controller) *MockResource {
	mock := &MockResource{ctrl: ctrl}
	mock.recorder = &MockResourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResource) EXPECT() *MockResourceMockRecorder {
	return m.recorder
}

// Initialize mocks base method.
func (m *MockResource) Initialize(arg0 resource.Manager) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockResourceMockRecorder) Initialize(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockResource)(nil).Initialize), arg0)
}