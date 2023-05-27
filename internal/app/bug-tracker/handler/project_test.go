package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	mock_handler "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/handler/mocks"
	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
	mock_services "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services/mocks"
)

func Test_createProject(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectData *dto.CreateProjectDto) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		projectData        *dto.CreateProjectDto
		projectDataJSON    string
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid json",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.CreateProjectDto) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			projectData:        nil,
			projectDataJSON:    `{"invalid"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid create project data",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.CreateProjectDto) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			projectData:        nil,
			projectDataJSON:    `{"name": "N"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidProjectData.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot create project",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.CreateProjectDto) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().CreateProject(projectData).Return(uint64(0), err)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, nil}
			},
			projectData: &dto.CreateProjectDto{
				Name:        "name",
				Description: "description",
				AdminID:     1,
			},
			projectDataJSON:    `{"name": "name", "description": "description", "adminId": 1}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + errInternalServerError.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.CreateProjectDto) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().CreateProject(projectData).Return(uint64(1), nil)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, nil}
			},
			projectData: &dto.CreateProjectDto{
				Name:        "name",
				Description: "description",
				AdminID:     1,
			},
			projectDataJSON:    `{"name": "name", "description": "description", "adminId": 1}`,
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "1" + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			handler := test.mockBehaviour(c, test.projectData)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.POST(create, handler.createProject)

			req := httptest.NewRequest(http.MethodPost, create, strings.NewReader(test.projectDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.createProject(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_getProjectById(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, id uint64, ctx echo.Context) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		id                 uint64
		paramId            string
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error in params.GetIdParam",
			mockBehaviour: func(c *gomock.Controller, id uint64, ctx echo.Context) *Handler {
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(uint64(0), err)

				return &Handler{params: params}
			},
			paramId:            "1b",
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot get project",
			mockBehaviour: func(c *gomock.Controller, id uint64, ctx echo.Context) *Handler {
				project := mock_services.NewMockProject(c)
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(id, nil)

				project.EXPECT().GetProjectById(id).Return(nil, err)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, params}
			},
			id:                 1,
			paramId:            "1",
			expectedStatusCode: http.StatusNotFound,
			expectedReturnBody: `{"message":"` + errProjectNotFound.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, id uint64, ctx echo.Context) *Handler {
				project := mock_services.NewMockProject(c)
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(id, nil)

				project.EXPECT().GetProjectById(id).Return(&models.Project{ID: 1, Name: "name", AdminID: 1}, nil)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, params}
			},
			id:                 1,
			paramId:            "1",
			expectedStatusCode: http.StatusFound,
			expectedReturnBody: `{"id":1,"name":"name","description":"","admin":1}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			echoCtx := e.NewContext(req, rec)

			handler := test.mockBehaviour(c, test.id, echoCtx)
			e.GET(id, handler.getProjectById)

			echoCtx.SetPath("/:id")
			echoCtx.SetParamNames("id")
			echoCtx.SetParamValues(test.paramId)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.getProjectById(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_deleteProject(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectID, userID uint64, ctx echo.Context) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		id                 uint64
		userData           *services.TokenData
		paramId            string
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error in params.GetIdParam",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64, ctx echo.Context) *Handler {
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(uint64(0), err)

				return &Handler{params: params}
			},
			paramId:            "1b",
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid user data",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64, ctx echo.Context) *Handler {
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(projectID, nil)

				return &Handler{params: params}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot delete project",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64, ctx echo.Context) *Handler {
				project := mock_services.NewMockProject(c)
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(projectID, nil)

				project.EXPECT().DeleteProject(projectID, userID).Return(err)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, params}
			},
			id: 1,
			userData: &services.TokenData{
				UserID: 1,
			},
			paramId:            "1",
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64, ctx echo.Context) *Handler {
				project := mock_services.NewMockProject(c)
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(projectID, nil)

				project.EXPECT().DeleteProject(projectID, userID).Return(nil)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, params}
			},
			id: 1,
			userData: &services.TokenData{
				UserID: 1,
			},
			paramId:            "1",
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: `true` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userID := uint64(0)
			if test.userData != nil {
				userID = test.userData.UserID
			}

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			echoCtx := e.NewContext(req, rec)

			handler := test.mockBehaviour(c, test.id, userID, echoCtx)
			e.DELETE(id, handler.deleteProject)

			echoCtx.SetPath("/:id")
			echoCtx.SetParamNames("id")
			echoCtx.SetParamValues(test.paramId)
			echoCtx.Set(userDataCtx, test.userData)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.deleteProject(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_updateProject(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		projectData        *dto.UpdateProjectDto
		projectDataJSON    string
		userData           *services.TokenData
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid user data",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *Handler {
				return &Handler{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid JSON",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{log: log}
			},
			userData:           &services.TokenData{UserID: 1},
			projectDataJSON:    "{invalid}",
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot update project",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().UpdateProject(projectData, userID).Return(err)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, nil}
			},
			projectData:     &dto.UpdateProjectDto{Description: "description", ProjectID: 1},
			projectDataJSON: `{"projectId": 1, "description": "description"}`,
			userData: &services.TokenData{
				UserID: 1,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + errInternalServerError.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().UpdateProject(projectData, userID).Return(nil)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, nil}
			},
			projectData:     &dto.UpdateProjectDto{Description: "description", ProjectID: 1},
			projectDataJSON: `{"projectId": 1, "description": "description"}`,
			userData: &services.TokenData{
				UserID: 1,
			},
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: `true` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userID := uint64(0)
			if test.userData != nil {
				userID = test.userData.UserID
			}

			handler := test.mockBehaviour(c, test.projectData, userID)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.PUT(id, handler.updateProject)

			req := httptest.NewRequest(http.MethodPut, update, strings.NewReader(test.projectDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
			echoCtx.Set(userDataCtx, test.userData)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.updateProject(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_addMember(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		memberData         *dto.AddMemberDto
		memberDataJSON     string
		userData           *services.TokenData
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid user data",
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *Handler {
				return &Handler{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid JSON",
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{log: log}
			},
			userData:           &services.TokenData{UserID: 1},
			memberDataJSON:     "{invalid}",
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error memberID == adminID",
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *Handler {
				return &Handler{}
			},
			userData:           &services.TokenData{UserID: 1},
			memberDataJSON:     `{"projectId": 1, "memberId": 1}`,
			memberData:         &dto.AddMemberDto{ProjectID: 1, MemberID: 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidOperation.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot add member to project",
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().AddMember(memberData, userID).Return(err)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, nil}
			},
			memberData:     &dto.AddMemberDto{MemberID: 2, ProjectID: 1},
			memberDataJSON: `{"projectId": 1, "memberId": 2}`,
			userData: &services.TokenData{
				UserID: 1,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + errInternalServerError.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().AddMember(memberData, userID).Return(nil)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, nil}
			},
			memberData:     &dto.AddMemberDto{MemberID: 2, ProjectID: 1},
			memberDataJSON: `{"projectId": 1, "memberId": 2}`,
			userData: &services.TokenData{
				UserID: 1,
			},
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: `true` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userID := uint64(0)
			if test.userData != nil {
				userID = test.userData.UserID
			}

			handler := test.mockBehaviour(c, test.memberData, userID)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.POST(addMember, handler.addMember)

			req := httptest.NewRequest(http.MethodPost, addMember, strings.NewReader(test.memberDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
			echoCtx.Set(userDataCtx, test.userData)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.addMember(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_leaveProject(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectID, userID uint64, ctx echo.Context) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		id                 uint64
		userData           *services.TokenData
		paramId            string
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error in params.GetIdParam",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64, ctx echo.Context) *Handler {
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(uint64(0), err)

				return &Handler{params: params}
			},
			paramId:            "1b",
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid user data",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64, ctx echo.Context) *Handler {
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(projectID, nil)

				return &Handler{params: params}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot leave project",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64, ctx echo.Context) *Handler {
				project := mock_services.NewMockProject(c)
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(projectID, nil)

				project.EXPECT().LeaveProject(projectID, userID).Return(err)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, params}
			},
			id: 1,
			userData: &services.TokenData{
				UserID: 1,
			},
			paramId:            "1",
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64, ctx echo.Context) *Handler {
				project := mock_services.NewMockProject(c)
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(projectID, nil)

				project.EXPECT().LeaveProject(projectID, userID).Return(nil)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, params}
			},
			id: 1,
			userData: &services.TokenData{
				UserID: 1,
			},
			paramId:            "1",
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: `true` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userID := uint64(0)
			if test.userData != nil {
				userID = test.userData.UserID
			}

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			echoCtx := e.NewContext(req, rec)

			handler := test.mockBehaviour(c, test.id, userID, echoCtx)
			e.GET(id, handler.leaveProject)

			echoCtx.SetPath("/:id")
			echoCtx.SetParamNames("id")
			echoCtx.SetParamValues(test.paramId)
			echoCtx.Set(userDataCtx, test.userData)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.leaveProject(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_setNewAdmin(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, newAdminData *dto.NewAdminDto, userID uint64) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		newAdminData       *dto.NewAdminDto
		newAdminDataJSON   string
		userData           *services.TokenData
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid user data",
			mockBehaviour: func(c *gomock.Controller, newAdminData *dto.NewAdminDto, userID uint64) *Handler {
				return &Handler{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid JSON",
			mockBehaviour: func(c *gomock.Controller, newAdminData *dto.NewAdminDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{log: log}
			},
			userData:           &services.TokenData{UserID: 1},
			newAdminDataJSON:   "{invalid}",
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot set new admin",
			mockBehaviour: func(c *gomock.Controller, newAdminData *dto.NewAdminDto, userID uint64) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().SetNewAdmin(newAdminData, userID).Return(err)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, nil}
			},
			newAdminData:     &dto.NewAdminDto{ProjectID: 1, NewAdminID: 2},
			newAdminDataJSON: `{"projectId": 1, "newAdminId": 2}`,
			userData: &services.TokenData{
				UserID: 1,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, newAdminData *dto.NewAdminDto, userID uint64) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().SetNewAdmin(newAdminData, userID).Return(nil)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil, nil}
			},
			newAdminData:     &dto.NewAdminDto{ProjectID: 1, NewAdminID: 2},
			newAdminDataJSON: `{"projectId": 1, "newAdminId": 2}`,
			userData: &services.TokenData{
				UserID: 1,
			},
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: `true` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userID := uint64(0)
			if test.userData != nil {
				userID = test.userData.UserID
			}

			handler := test.mockBehaviour(c, test.newAdminData, userID)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.POST(setAdmin, handler.setNewAdmin)

			req := httptest.NewRequest(http.MethodPost, addMember, strings.NewReader(test.newAdminDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
			echoCtx.Set(userDataCtx, test.userData)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.setNewAdmin(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}
