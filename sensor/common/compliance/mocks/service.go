// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/stackrox/rox/sensor/common/compliance (interfaces: Service)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	compliance "github.com/stackrox/rox/generated/internalapi/compliance"
	sensor "github.com/stackrox/rox/generated/internalapi/sensor"
	storage "github.com/stackrox/rox/generated/storage"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AuditEvents mocks base method.
func (m *MockService) AuditEvents() chan *sensor.AuditEvents {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuditEvents")
	ret0, _ := ret[0].(chan *sensor.AuditEvents)
	return ret0
}

// AuditEvents indicates an expected call of AuditEvents.
func (mr *MockServiceMockRecorder) AuditEvents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuditEvents", reflect.TypeOf((*MockService)(nil).AuditEvents))
}

// AuthFuncOverride mocks base method.
func (m *MockService) AuthFuncOverride(arg0 context.Context, arg1 string) (context.Context, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthFuncOverride", arg0, arg1)
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthFuncOverride indicates an expected call of AuthFuncOverride.
func (mr *MockServiceMockRecorder) AuthFuncOverride(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthFuncOverride", reflect.TypeOf((*MockService)(nil).AuthFuncOverride), arg0, arg1)
}

// Communicate mocks base method.
func (m *MockService) Communicate(arg0 sensor.ComplianceService_CommunicateServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Communicate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Communicate indicates an expected call of Communicate.
func (mr *MockServiceMockRecorder) Communicate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Communicate", reflect.TypeOf((*MockService)(nil).Communicate), arg0)
}

// NodeInventories mocks base method.
func (m *MockService) NodeInventories() <-chan *storage.NodeInventory {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NodeInventories")
	ret0, _ := ret[0].(<-chan *storage.NodeInventory)
	return ret0
}

// NodeInventories indicates an expected call of NodeInventories.
func (mr *MockServiceMockRecorder) NodeInventories() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NodeInventories", reflect.TypeOf((*MockService)(nil).NodeInventories))
}

// Output mocks base method.
func (m *MockService) Output() chan *compliance.ComplianceReturn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Output")
	ret0, _ := ret[0].(chan *compliance.ComplianceReturn)
	return ret0
}

// Output indicates an expected call of Output.
func (mr *MockServiceMockRecorder) Output() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Output", reflect.TypeOf((*MockService)(nil).Output))
}

// RegisterServiceHandler mocks base method.
func (m *MockService) RegisterServiceHandler(arg0 context.Context, arg1 *runtime.ServeMux, arg2 *grpc.ClientConn) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterServiceHandler", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterServiceHandler indicates an expected call of RegisterServiceHandler.
func (mr *MockServiceMockRecorder) RegisterServiceHandler(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterServiceHandler", reflect.TypeOf((*MockService)(nil).RegisterServiceHandler), arg0, arg1, arg2)
}

// RegisterServiceServer mocks base method.
func (m *MockService) RegisterServiceServer(arg0 *grpc.Server) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterServiceServer", arg0)
}

// RegisterServiceServer indicates an expected call of RegisterServiceServer.
func (mr *MockServiceMockRecorder) RegisterServiceServer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterServiceServer", reflect.TypeOf((*MockService)(nil).RegisterServiceServer), arg0)
}

// RunScrape mocks base method.
func (m *MockService) RunScrape(arg0 *sensor.MsgToCompliance) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunScrape", arg0)
	ret0, _ := ret[0].(int)
	return ret0
}

// RunScrape indicates an expected call of RunScrape.
func (mr *MockServiceMockRecorder) RunScrape(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunScrape", reflect.TypeOf((*MockService)(nil).RunScrape), arg0)
}