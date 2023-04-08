// Code generated by MockGen. DO NOT EDIT.
// Source: kafka.go

// Package mock_kafka is a generated GoMock package.
package mock_kafka

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockKafka is a mock of Kafka interface.
type MockKafka struct {
	ctrl     *gomock.Controller
	recorder *MockKafkaMockRecorder
}

// MockKafkaMockRecorder is the mock recorder for MockKafka.
type MockKafkaMockRecorder struct {
	mock *MockKafka
}

// NewMockKafka creates a new mock instance.
func NewMockKafka(ctrl *gomock.Controller) *MockKafka {
	mock := &MockKafka{ctrl: ctrl}
	mock.recorder = &MockKafkaMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKafka) EXPECT() *MockKafkaMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockKafka) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockKafkaMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockKafka)(nil).Close))
}

// Write mocks base method.
func (m *MockKafka) Write(message string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Write indicates an expected call of Write.
func (mr *MockKafkaMockRecorder) Write(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockKafka)(nil).Write), message)
}