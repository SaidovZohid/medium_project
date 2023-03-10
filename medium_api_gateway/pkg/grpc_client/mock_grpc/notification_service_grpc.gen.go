// Code generated by MockGen. DO NOT EDIT.
// Source: ../../genproto/notification_service/notification_service_grpc.pb.go

// Package mock_grpc is a generated GoMock package.
package mock_grpc

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	empty "github.com/golang/protobuf/ptypes/empty"
	notification_service "gitlab.com/medium-project/medium_api_gateway/genproto/notification_service"
	grpc "google.golang.org/grpc"
)

// MockNotificationServiceClient is a mock of NotificationServiceClient interface.
type MockNotificationServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceClientMockRecorder
}

// MockNotificationServiceClientMockRecorder is the mock recorder for MockNotificationServiceClient.
type MockNotificationServiceClientMockRecorder struct {
	mock *MockNotificationServiceClient
}

// NewMockNotificationServiceClient creates a new mock instance.
func NewMockNotificationServiceClient(ctrl *gomock.Controller) *MockNotificationServiceClient {
	mock := &MockNotificationServiceClient{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationServiceClient) EXPECT() *MockNotificationServiceClientMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockNotificationServiceClient) SendEmail(ctx context.Context, in *notification_service.SendEmailRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendEmail", varargs...)
	ret0, _ := ret[0].(*empty.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockNotificationServiceClientMockRecorder) SendEmail(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockNotificationServiceClient)(nil).SendEmail), varargs...)
}

// MockNotificationServiceServer is a mock of NotificationServiceServer interface.
type MockNotificationServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceServerMockRecorder
}

// MockNotificationServiceServerMockRecorder is the mock recorder for MockNotificationServiceServer.
type MockNotificationServiceServerMockRecorder struct {
	mock *MockNotificationServiceServer
}

// NewMockNotificationServiceServer creates a new mock instance.
func NewMockNotificationServiceServer(ctrl *gomock.Controller) *MockNotificationServiceServer {
	mock := &MockNotificationServiceServer{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationServiceServer) EXPECT() *MockNotificationServiceServerMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockNotificationServiceServer) SendEmail(arg0 context.Context, arg1 *notification_service.SendEmailRequest) (*empty.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", arg0, arg1)
	ret0, _ := ret[0].(*empty.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockNotificationServiceServerMockRecorder) SendEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockNotificationServiceServer)(nil).SendEmail), arg0, arg1)
}

// mustEmbedUnimplementedNotificationServiceServer mocks base method.
func (m *MockNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNotificationServiceServer")
}

// mustEmbedUnimplementedNotificationServiceServer indicates an expected call of mustEmbedUnimplementedNotificationServiceServer.
func (mr *MockNotificationServiceServerMockRecorder) mustEmbedUnimplementedNotificationServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNotificationServiceServer", reflect.TypeOf((*MockNotificationServiceServer)(nil).mustEmbedUnimplementedNotificationServiceServer))
}

// MockUnsafeNotificationServiceServer is a mock of UnsafeNotificationServiceServer interface.
type MockUnsafeNotificationServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeNotificationServiceServerMockRecorder
}

// MockUnsafeNotificationServiceServerMockRecorder is the mock recorder for MockUnsafeNotificationServiceServer.
type MockUnsafeNotificationServiceServerMockRecorder struct {
	mock *MockUnsafeNotificationServiceServer
}

// NewMockUnsafeNotificationServiceServer creates a new mock instance.
func NewMockUnsafeNotificationServiceServer(ctrl *gomock.Controller) *MockUnsafeNotificationServiceServer {
	mock := &MockUnsafeNotificationServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeNotificationServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeNotificationServiceServer) EXPECT() *MockUnsafeNotificationServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedNotificationServiceServer mocks base method.
func (m *MockUnsafeNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNotificationServiceServer")
}

// mustEmbedUnimplementedNotificationServiceServer indicates an expected call of mustEmbedUnimplementedNotificationServiceServer.
func (mr *MockUnsafeNotificationServiceServerMockRecorder) mustEmbedUnimplementedNotificationServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNotificationServiceServer", reflect.TypeOf((*MockUnsafeNotificationServiceServer)(nil).mustEmbedUnimplementedNotificationServiceServer))
}
