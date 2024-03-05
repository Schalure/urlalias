// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Schalure/urlalias/internal/app/aliasmaker (interfaces: Storager)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	aliasentity "github.com/Schalure/urlalias/internal/app/models/aliasentity"
)

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorager) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStoragerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorager)(nil).Close))
}

// CreateUser mocks base method.
func (m *MockStorager) CreateUser() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoragerMockRecorder) CreateUser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStorager)(nil).CreateUser))
}

// FindAllByLongURLs mocks base method.
func (m *MockStorager) FindAllByLongURLs(arg0 context.Context, arg1 []string) (map[string]aliasentity.AliasURLModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByLongURLs", arg0, arg1)
	ret0, _ := ret[0].(map[string]aliasentity.AliasURLModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByLongURLs indicates an expected call of FindAllByLongURLs.
func (mr *MockStoragerMockRecorder) FindAllByLongURLs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByLongURLs", reflect.TypeOf((*MockStorager)(nil).FindAllByLongURLs), arg0, arg1)
}

// FindByLongURL mocks base method.
func (m *MockStorager) FindByLongURL(arg0 context.Context, arg1 string) (*aliasentity.AliasURLModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByLongURL", arg0, arg1)
	ret0, _ := ret[0].(*aliasentity.AliasURLModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByLongURL indicates an expected call of FindByLongURL.
func (mr *MockStoragerMockRecorder) FindByLongURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByLongURL", reflect.TypeOf((*MockStorager)(nil).FindByLongURL), arg0, arg1)
}

// FindByShortKey mocks base method.
func (m *MockStorager) FindByShortKey(arg0 context.Context, arg1 string) (*aliasentity.AliasURLModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByShortKey", arg0, arg1)
	ret0, _ := ret[0].(*aliasentity.AliasURLModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByShortKey indicates an expected call of FindByShortKey.
func (mr *MockStoragerMockRecorder) FindByShortKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByShortKey", reflect.TypeOf((*MockStorager)(nil).FindByShortKey), arg0, arg1)
}

// FindByUserID mocks base method.
func (m *MockStorager) FindByUserID(arg0 context.Context, arg1 uint64) ([]aliasentity.AliasURLModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserID", arg0, arg1)
	ret0, _ := ret[0].([]aliasentity.AliasURLModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserID indicates an expected call of FindByUserID.
func (mr *MockStoragerMockRecorder) FindByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserID", reflect.TypeOf((*MockStorager)(nil).FindByUserID), arg0, arg1)
}

// GetLastShortKey mocks base method.
func (m *MockStorager) GetLastShortKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastShortKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetLastShortKey indicates an expected call of GetLastShortKey.
func (mr *MockStoragerMockRecorder) GetLastShortKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastShortKey", reflect.TypeOf((*MockStorager)(nil).GetLastShortKey))
}

// IsConnected mocks base method.
func (m *MockStorager) IsConnected() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsConnected")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsConnected indicates an expected call of IsConnected.
func (mr *MockStoragerMockRecorder) IsConnected() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsConnected", reflect.TypeOf((*MockStorager)(nil).IsConnected))
}

// MarkDeleted mocks base method.
func (m *MockStorager) MarkDeleted(arg0 context.Context, arg1 []uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkDeleted", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkDeleted indicates an expected call of MarkDeleted.
func (mr *MockStoragerMockRecorder) MarkDeleted(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkDeleted", reflect.TypeOf((*MockStorager)(nil).MarkDeleted), arg0, arg1)
}

// Save mocks base method.
func (m *MockStorager) Save(arg0 context.Context, arg1 *aliasentity.AliasURLModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockStoragerMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStorager)(nil).Save), arg0, arg1)
}

// SaveAll mocks base method.
func (m *MockStorager) SaveAll(arg0 context.Context, arg1 []aliasentity.AliasURLModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAll", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAll indicates an expected call of SaveAll.
func (mr *MockStoragerMockRecorder) SaveAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAll", reflect.TypeOf((*MockStorager)(nil).SaveAll), arg0, arg1)
}
