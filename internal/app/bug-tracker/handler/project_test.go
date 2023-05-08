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
			expectedReturnBody: `{"message":"` + errInvalidCreateProjectData.Error() + `"}` + "\n",
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