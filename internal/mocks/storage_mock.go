package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/madcarpet/metrics/internal/entity"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockRepository) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockRepositoryMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepository)(nil).Close))
}

// ExportToFile mocks base method.
func (m *MockRepository) ExportToFile() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExportToFile")
	ret0, _ := ret[0].(error)
	return ret0
}

// ExportToFile indicates an expected call of ExportToFile.
func (mr *MockRepositoryMockRecorder) ExportToFile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportToFile", reflect.TypeOf((*MockRepository)(nil).ExportToFile))
}

// GetAllMetrics mocks base method.
func (m *MockRepository) GetAllMetrics(arg0 context.Context) []entity.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMetrics", arg0)
	ret0, _ := ret[0].([]entity.Metric)
	return ret0
}

// GetAllMetrics indicates an expected call of GetAllMetrics.
func (mr *MockRepositoryMockRecorder) GetAllMetrics(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMetrics", reflect.TypeOf((*MockRepository)(nil).GetAllMetrics), arg0)
}

// GetByNameAndType mocks base method.
func (m *MockRepository) GetByNameAndType(arg0 context.Context, arg1 string, arg2 int64) (entity.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNameAndType", arg0, arg1, arg2)
	ret0, _ := ret[0].(entity.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNameAndType indicates an expected call of GetByNameAndType.
func (mr *MockRepositoryMockRecorder) GetByNameAndType(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNameAndType", reflect.TypeOf((*MockRepository)(nil).GetByNameAndType), arg0, arg1, arg2)
}

// ImportFromFile mocks base method.
func (m *MockRepository) ImportFromFile() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImportFromFile")
	ret0, _ := ret[0].(error)
	return ret0
}

// ImportFromFile indicates an expected call of ImportFromFile.
func (mr *MockRepositoryMockRecorder) ImportFromFile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImportFromFile", reflect.TypeOf((*MockRepository)(nil).ImportFromFile))
}

// IsConnected mocks base method.
func (m *MockRepository) IsConnected(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsConnected", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// IsConnected indicates an expected call of IsConnected.
func (mr *MockRepositoryMockRecorder) IsConnected(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsConnected", reflect.TypeOf((*MockRepository)(nil).IsConnected), arg0)
}

// UpdateMetric mocks base method.
func (m *MockRepository) UpdateMetric(arg0 context.Context, arg1 entity.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetric", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetric indicates an expected call of UpdateMetric.
func (mr *MockRepositoryMockRecorder) UpdateMetric(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetric", reflect.TypeOf((*MockRepository)(nil).UpdateMetric), arg0, arg1)
}

// UpdateMetrics mocks base method.
func (m *MockRepository) UpdateMetrics(arg0 context.Context, arg1 []entity.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetrics", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetrics indicates an expected call of UpdateMetrics.
func (mr *MockRepositoryMockRecorder) UpdateMetrics(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetrics", reflect.TypeOf((*MockRepository)(nil).UpdateMetrics), arg0, arg1)
}
