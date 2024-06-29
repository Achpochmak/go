// Code generated by MockGen. DO NOT EDIT.
// Source: ./init.go

// Package mock_cli is a generated GoMock package.
package mock_cli

import (
	models "HOMEWORK-1/internal/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockModule is a mock of Module interface.
type MockModule struct {
	ctrl     *gomock.Controller
	recorder *MockModuleMockRecorder
}

// MockModuleMockRecorder is the mock recorder for MockModule.
type MockModuleMockRecorder struct {
	mock *MockModule
}

// NewMockModule creates a new mock instance.
func NewMockModule(ctrl *gomock.Controller) *MockModule {
	mock := &MockModule{ctrl: ctrl}
	mock.recorder = &MockModuleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModule) EXPECT() *MockModuleMockRecorder {
	return m.recorder
}

// AddOrder mocks base method.
func (m *MockModule) AddOrder(arg0 context.Context, arg1 models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockModuleMockRecorder) AddOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockModule)(nil).AddOrder), arg0, arg1)
}

// DeleteOrder mocks base method.
func (m *MockModule) DeleteOrder(arg0 context.Context, arg1 models.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrder indicates an expected call of DeleteOrder.
func (mr *MockModuleMockRecorder) DeleteOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrder", reflect.TypeOf((*MockModule)(nil).DeleteOrder), arg0, arg1)
}

// DeliverOrder mocks base method.
func (m *MockModule) DeliverOrder(arg0 context.Context, arg1 []int, arg2 int) ([]*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeliverOrder", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeliverOrder indicates an expected call of DeliverOrder.
func (mr *MockModuleMockRecorder) DeliverOrder(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeliverOrder", reflect.TypeOf((*MockModule)(nil).DeliverOrder), arg0, arg1, arg2)
}

// GetOrderByID mocks base method.
func (m *MockModule) GetOrderByID(arg0 context.Context, arg1 models.ID) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByID", arg0, arg1)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByID indicates an expected call of GetOrderByID.
func (mr *MockModuleMockRecorder) GetOrderByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByID", reflect.TypeOf((*MockModule)(nil).GetOrderByID), arg0, arg1)
}

// GetOrdersByCustomer mocks base method.
func (m *MockModule) GetOrdersByCustomer(arg0 context.Context, arg1, arg2 int) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersByCustomer", arg0, arg1, arg2)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersByCustomer indicates an expected call of GetOrdersByCustomer.
func (mr *MockModuleMockRecorder) GetOrdersByCustomer(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersByCustomer", reflect.TypeOf((*MockModule)(nil).GetOrdersByCustomer), arg0, arg1, arg2)
}

// ListOrder mocks base method.
func (m *MockModule) ListOrder(arg0 context.Context) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOrder", arg0)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListOrder indicates an expected call of ListOrder.
func (mr *MockModuleMockRecorder) ListOrder(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOrder", reflect.TypeOf((*MockModule)(nil).ListOrder), arg0)
}

// ListRefund mocks base method.
func (m *MockModule) ListRefund(arg0 context.Context, arg1, arg2 int) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRefund", arg0, arg1, arg2)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRefund indicates an expected call of ListRefund.
func (mr *MockModuleMockRecorder) ListRefund(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRefund", reflect.TypeOf((*MockModule)(nil).ListRefund), arg0, arg1, arg2)
}

// Refund mocks base method.
func (m *MockModule) Refund(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refund", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Refund indicates an expected call of Refund.
func (mr *MockModuleMockRecorder) Refund(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refund", reflect.TypeOf((*MockModule)(nil).Refund), arg0, arg1, arg2)
}
