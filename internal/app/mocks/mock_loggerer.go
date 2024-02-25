// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Schalure/urlalias/internal/app/aliasmaker (interfaces: Loggerer)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLoggerer is a mock of Loggerer interface.
type MockLoggerer struct {
	ctrl     *gomock.Controller
	recorder *MockLoggererMockRecorder
}

// MockLoggererMockRecorder is the mock recorder for MockLoggerer.
type MockLoggererMockRecorder struct {
	mock *MockLoggerer
}

// NewMockLoggerer creates a new mock instance.
func NewMockLoggerer(ctrl *gomock.Controller) *MockLoggerer {
	mock := &MockLoggerer{ctrl: ctrl}
	mock.recorder = &MockLoggererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoggerer) EXPECT() *MockLoggererMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockLoggerer) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockLoggererMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockLoggerer)(nil).Close))
}

// Errorw mocks base method.
func (m *MockLoggerer) Errorw(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorw", varargs...)
}

// Errorw indicates an expected call of Errorw.
func (mr *MockLoggererMockRecorder) Errorw(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorw", reflect.TypeOf((*MockLoggerer)(nil).Errorw), varargs...)
}

// Fatalw mocks base method.
func (m *MockLoggerer) Fatalw(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatalw", varargs...)
}

// Fatalw indicates an expected call of Fatalw.
func (mr *MockLoggererMockRecorder) Fatalw(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatalw", reflect.TypeOf((*MockLoggerer)(nil).Fatalw), varargs...)
}

// Info mocks base method.
func (m *MockLoggerer) Info(arg0 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockLoggererMockRecorder) Info(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLoggerer)(nil).Info), arg0...)
}

// Infow mocks base method.
func (m *MockLoggerer) Infow(arg0 string, arg1 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infow", varargs...)
}

// Infow indicates an expected call of Infow.
func (mr *MockLoggererMockRecorder) Infow(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infow", reflect.TypeOf((*MockLoggerer)(nil).Infow), varargs...)
}