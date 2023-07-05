// Code generated by MockGen. DO NOT EDIT.
// Source: services.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	models "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	services "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
)

// MockAuth is a mock of Auth interface.
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth.
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance.
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// GenerateAccessToken mocks base method.
func (m *MockAuth) GenerateAccessToken(username string, userID uint64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAccessToken", username, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateAccessToken indicates an expected call of GenerateAccessToken.
func (mr *MockAuthMockRecorder) GenerateAccessToken(username, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAccessToken", reflect.TypeOf((*MockAuth)(nil).GenerateAccessToken), username, userID)
}

// GenerateRefreshToken mocks base method.
func (m *MockAuth) GenerateRefreshToken(username string, userID uint64) (*services.RefreshTokenData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRefreshToken", username, userID)
	ret0, _ := ret[0].(*services.RefreshTokenData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateRefreshToken indicates an expected call of GenerateRefreshToken.
func (mr *MockAuthMockRecorder) GenerateRefreshToken(username, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRefreshToken", reflect.TypeOf((*MockAuth)(nil).GenerateRefreshToken), username, userID)
}

// GetRefreshTokenTTL mocks base method.
func (m *MockAuth) GetRefreshTokenTTL() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefreshTokenTTL")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// GetRefreshTokenTTL indicates an expected call of GetRefreshTokenTTL.
func (mr *MockAuthMockRecorder) GetRefreshTokenTTL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefreshTokenTTL", reflect.TypeOf((*MockAuth)(nil).GetRefreshTokenTTL))
}

// ParseAccessToken mocks base method.
func (m *MockAuth) ParseAccessToken(accessToken string) (*services.TokenData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseAccessToken", accessToken)
	ret0, _ := ret[0].(*services.TokenData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseAccessToken indicates an expected call of ParseAccessToken.
func (mr *MockAuthMockRecorder) ParseAccessToken(accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseAccessToken", reflect.TypeOf((*MockAuth)(nil).ParseAccessToken), accessToken)
}

// ParseRefreshToken mocks base method.
func (m *MockAuth) ParseRefreshToken(refreshToken string) (*services.TokenData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseRefreshToken", refreshToken)
	ret0, _ := ret[0].(*services.TokenData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseRefreshToken indicates an expected call of ParseRefreshToken.
func (mr *MockAuthMockRecorder) ParseRefreshToken(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseRefreshToken", reflect.TypeOf((*MockAuth)(nil).ParseRefreshToken), refreshToken)
}

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

// ValidateUser mocks base method.
func (m *MockUser) ValidateUser(email, password string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateUser", email, password)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateUser indicates an expected call of ValidateUser.
func (mr *MockUserMockRecorder) ValidateUser(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateUser", reflect.TypeOf((*MockUser)(nil).ValidateUser), email, password)
}

// MockRedis is a mock of Redis interface.
type MockRedis struct {
	ctrl     *gomock.Controller
	recorder *MockRedisMockRecorder
}

// MockRedisMockRecorder is the mock recorder for MockRedis.
type MockRedisMockRecorder struct {
	mock *MockRedis
}

// NewMockRedis creates a new mock instance.
func NewMockRedis(ctrl *gomock.Controller) *MockRedis {
	mock := &MockRedis{ctrl: ctrl}
	mock.recorder = &MockRedisMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedis) EXPECT() *MockRedisMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockRedis) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockRedisMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRedis)(nil).Close))
}

// DeleteRefreshToken mocks base method.
func (m *MockRedis) DeleteRefreshToken(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRefreshToken", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRefreshToken indicates an expected call of DeleteRefreshToken.
func (mr *MockRedisMockRecorder) DeleteRefreshToken(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRefreshToken", reflect.TypeOf((*MockRedis)(nil).DeleteRefreshToken), ctx, key)
}

// Get mocks base method.
func (m *MockRedis) Get(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRedisMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRedis)(nil).Get), ctx, key)
}

// GetRefreshToken mocks base method.
func (m *MockRedis) GetRefreshToken(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefreshToken", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRefreshToken indicates an expected call of GetRefreshToken.
func (mr *MockRedisMockRecorder) GetRefreshToken(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefreshToken", reflect.TypeOf((*MockRedis)(nil).GetRefreshToken), ctx, key)
}

// Set mocks base method.
func (m *MockRedis) Set(ctx context.Context, key, val string, exp time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, val, exp)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockRedisMockRecorder) Set(ctx, key, val, exp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRedis)(nil).Set), ctx, key, val, exp)
}

// SetRefreshToken mocks base method.
func (m *MockRedis) SetRefreshToken(ctx context.Context, key, refreshToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetRefreshToken", ctx, key, refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRefreshToken indicates an expected call of SetRefreshToken.
func (mr *MockRedisMockRecorder) SetRefreshToken(ctx, key, refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRefreshToken", reflect.TypeOf((*MockRedis)(nil).SetRefreshToken), ctx, key, refreshToken)
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
func (m *MockProject) SetNewAdmin(newAdmintData *dto.NewAdminDto, adminID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNewAdmin", newAdmintData, adminID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetNewAdmin indicates an expected call of SetNewAdmin.
func (mr *MockProjectMockRecorder) SetNewAdmin(newAdmintData, adminID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNewAdmin", reflect.TypeOf((*MockProject)(nil).SetNewAdmin), newAdmintData, adminID)
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
