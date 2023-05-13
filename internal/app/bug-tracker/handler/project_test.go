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

				return &Handler{nil, log, nil}
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

				return &Handler{nil, log, nil}
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

				return &Handler{serv, nil, nil}
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

				return &Handler{serv, nil, nil}
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
	type mockBehaviour func(c *gomock.Controller, id uint64) *Handler
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
			name: "Error empty param",
			mockBehaviour: func(c *gomock.Controller, id uint64) *Handler {
				return &Handler{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errProjectNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid param",
			mockBehaviour: func(c *gomock.Controller, id uint64) *Handler {
				return &Handler{}
			},
			paramId:            "1b",
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errProjectNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot get project",
			mockBehaviour: func(c *gomock.Controller, id uint64) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().GetProjectById(id).Return(nil, err)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil}
			},
			id:                 1,
			paramId:            "1",
			expectedStatusCode: http.StatusNotFound,
			expectedReturnBody: `{"message":"` + errProjectNotFound.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, id uint64) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().GetProjectById(id).Return(&models.Project{ID: 1, Name: "name", AdminID: 1}, nil)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil}
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

			handler := test.mockBehaviour(c, test.id)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.GET(id, handler.getProjectById)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
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
	type mockBehaviour func(c *gomock.Controller, projectID, userID uint64) *Handler
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
			name: "Error invalid user data",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *Handler {
				return &Handler{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error empty param",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *Handler {
				return &Handler{}
			},
			userData: &services.TokenData{
				UserID: 1,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errProjectNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid param",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *Handler {
				return &Handler{}
			},
			userData: &services.TokenData{
				UserID: 1,
			},
			paramId:            "1b",
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errProjectNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot delete project",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().DeleteProject(projectID, userID).Return(err)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil}
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
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().DeleteProject(projectID, userID).Return(nil)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil}
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

			handler := test.mockBehaviour(c, test.id, userID)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.DELETE(id, handler.deleteProject)

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
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
			userData: &services.TokenData{UserID: 1},
			projectDataJSON: "{invalid}",
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot update project",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *Handler {
				project := mock_services.NewMockProject(c)

				project.EXPECT().UpdateProject(projectData, userID).Return(err)

				serv := &services.Service{Project: project}

				return &Handler{serv, nil, nil}
			},
			projectData: &dto.UpdateProjectDto{Description: "description", ProjectID: 1},
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

				return &Handler{serv, nil, nil}
			},
			projectData: &dto.UpdateProjectDto{Description: "description", ProjectID: 1},
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

			req := httptest.NewRequest(http.MethodDelete, update, strings.NewReader(test.projectDataJSON))
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
