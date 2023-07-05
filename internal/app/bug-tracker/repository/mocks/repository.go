// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	models "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUser) CreateUser(userData *dto.SignUpDto) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", userData)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserMockRecorder) CreateUser(userData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUser)(nil).CreateUser), userData)
}

// GetUserByEmail mocks base method.
func (m *MockUser) GetUserByEmail(email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserMockRecorder) GetUserByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUser)(nil).GetUserByEmail), email)
}

// GetUserById mocks base method.
func (m *MockUser) GetUserById(id uint64) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", id)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockUserMockRecorder) GetUserById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockUser)(nil).GetUserById), id)
}

// GetUserByUsername mocks base method.
func (m *MockUser) GetUserByUsername(username string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", username)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockUserMockRecorder) GetUserByUsername(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockUser)(nil).GetUserByUsername), username)
}

// MockProject is a mock of Project interface.
type MockProject struct {
	ctrl     *gomock.Controller
	recorder *MockProjectMockRecorder
}

// MockProjectMockRecorder is the mock recorder for MockProject.
type MockProjectMockRecorder struct {
	mock *MockProject
}

// NewMockProject creates a new mock instance.
func NewMockProject(ctrl *gomock.Controller) *MockProject {
	mock := &MockProject{ctrl: ctrl}
	mock.recorder = &MockProjectMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProject) EXPECT() *MockProjectMockRecorder {
	return m.recorder
}

// AddMember mocks base method.
func (m *MockProject) AddMember(memberData *dto.AddMemberDto, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMember", memberData, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMember indicates an expected call of AddMember.
func (mr *MockProjectMockRecorder) AddMember(memberData, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMember", reflect.TypeOf((*MockProject)(nil).AddMember), memberData, userID)
}

// CreateProject mocks base method.
func (m *MockProject) CreateProject(projectData *dto.CreateProjectDto) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProject", projectData)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProject indicates an expected call of CreateProject.
func (mr *MockProjectMockRecorder) CreateProject(projectData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProject", reflect.TypeOf((*MockProject)(nil).CreateProject), projectData)
}

// DeleteMember mocks base method.
func (m *MockProject) DeleteMember(memberData *dto.AddMemberDto, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMember", memberData, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMember indicates an expected call of DeleteMember.
func (mr *MockProjectMockRecorder) DeleteMember(memberData, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMember", reflect.TypeOf((*MockProject)(nil).DeleteMember), memberData, userID)
}

// DeleteProject mocks base method.
func (m *MockProject) DeleteProject(projectID, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProject", projectID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProject indicates an expected call of DeleteProject.
func (mr *MockProjectMockRecorder) DeleteProject(projectID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProject", reflect.TypeOf((*MockProject)(nil).DeleteProject), projectID, userID)
}

// GetProjectById mocks base method.
func (m *MockProject) GetProjectById(id uint64) (*models.Project, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectById", id)
	ret0, _ := ret[0].(*models.Project)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectById indicates an expected call of GetProjectById.
func (mr *MockProjectMockRecorder) GetProjectById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectById", reflect.TypeOf((*MockProject)(nil).GetProjectById), id)
}

// LeaveProject mocks base method.
func (m *MockProject) LeaveProject(projectID, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeaveProject", projectID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// LeaveProject indicates an expected call of LeaveProject.
func (mr *MockProjectMockRecorder) LeaveProject(projectID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveProject", reflect.TypeOf((*MockProject)(nil).LeaveProject), projectID, userID)
}

// SetNewAdmin mocks base method.
func (m *MockProject) SetNewAdmin(newAdminData *dto.NewAdminDto, adminID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNewAdmin", newAdminData, adminID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetNewAdmin indicates an expected call of SetNewAdmin.
func (mr *MockProjectMockRecorder) SetNewAdmin(newAdminData, adminID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNewAdmin", reflect.TypeOf((*MockProject)(nil).SetNewAdmin), newAdminData, adminID)
}

// UpdateProject mocks base method.
func (m *MockProject) UpdateProject(projectData *dto.UpdateProjectDto, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProject", projectData, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProject indicates an expected call of UpdateProject.
func (mr *MockProjectMockRecorder) UpdateProject(projectData, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProject", reflect.TypeOf((*MockProject)(nil).UpdateProject), projectData, userID)
}

// MockTask is a mock of Task interface.
type MockTask struct {
	ctrl     *gomock.Controller
	recorder *MockTaskMockRecorder
}

// MockTaskMockRecorder is the mock recorder for MockTask.
type MockTaskMockRecorder struct {
	mock *MockTask
}

// NewMockTask creates a new mock instance.
func NewMockTask(ctrl *gomock.Controller) *MockTask {
	mock := &MockTask{ctrl: ctrl}
	mock.recorder = &MockTaskMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTask) EXPECT() *MockTaskMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockTask) CreateTask(taskData *dto.CreateTaskDto, userID uint64) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", taskData, userID)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockTaskMockRecorder) CreateTask(taskData, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTask)(nil).CreateTask), taskData, userID)
}

// DeleteTask mocks base method.
func (m *MockTask) DeleteTask(taskData *dto.DeleteTaskDto, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", taskData, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockTaskMockRecorder) DeleteTask(taskData, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockTask)(nil).DeleteTask), taskData, userID)
}

// GetTaskById mocks base method.
func (m *MockTask) GetTaskById(id uint64) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskById", id)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskById indicates an expected call of GetTaskById.
func (mr *MockTaskMockRecorder) GetTaskById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskById", reflect.TypeOf((*MockTask)(nil).GetTaskById), id)
}

// GetTasksByProjectId mocks base method.
func (m *MockTask) GetTasksByProjectId(id uint64) ([]*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTasksByProjectId", id)
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksByProjectId indicates an expected call of GetTasksByProjectId.
func (mr *MockTaskMockRecorder) GetTasksByProjectId(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTasksByProjectId", reflect.TypeOf((*MockTask)(nil).GetTasksByProjectId), id)
}

// StopWorkOnTask mocks base method.
func (m *MockTask) StopWorkOnTask(workOnTaskData *dto.WorkOnTaskDto, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopWorkOnTask", workOnTaskData, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopWorkOnTask indicates an expected call of StopWorkOnTask.
func (mr *MockTaskMockRecorder) StopWorkOnTask(workOnTaskData, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopWorkOnTask", reflect.TypeOf((*MockTask)(nil).StopWorkOnTask), workOnTaskData, userID)
}

// UpdateTask mocks base method.
func (m *MockTask) UpdateTask(taskData *dto.UpdateTaskDto, userID uint64) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", taskData, userID)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTaskMockRecorder) UpdateTask(taskData, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTask)(nil).UpdateTask), taskData, userID)
}

// WorkOnTask mocks base method.
func (m *MockTask) WorkOnTask(workOnTaskData *dto.WorkOnTaskDto, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WorkOnTask", workOnTaskData, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// WorkOnTask indicates an expected call of WorkOnTask.
func (mr *MockTaskMockRecorder) WorkOnTask(workOnTaskData, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WorkOnTask", reflect.TypeOf((*MockTask)(nil).WorkOnTask), workOnTaskData, userID)
}
