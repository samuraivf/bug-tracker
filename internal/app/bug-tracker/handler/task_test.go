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
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
	mock_services "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services/mocks"
)

func Test_createTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, taskData *dto.CreateTaskDto) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		taskData           *dto.CreateTaskDto
		taskDataJSON       string
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid json",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			taskData:           nil,
			taskDataJSON:       `{"invalid"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid create task data",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			taskData:           nil,
			taskDataJSON:       `{"name": "N"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidTaskData.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot create task",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().CreateTask(taskData).Return(uint64(0), err)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, nil}
			},
			taskData: &dto.CreateTaskDto{
				Name:         "name",
				Description:  "description",
				TaskPriority: "high",
				ProjectID:    1,
				TaskType:     "TO DO",
			},
			taskDataJSON:       `{"name": "name", "description": "description", "taskPriority": "high", "projectId": 1, "taskType": "TO DO"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().CreateTask(taskData).Return(uint64(1), nil)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, nil}
			},
			taskData: &dto.CreateTaskDto{
				Name:         "name",
				Description:  "description",
				TaskPriority: "high",
				ProjectID:    1,
				TaskType:     "TO DO",
			},
			taskDataJSON:       `{"name": "name", "description": "description", "taskPriority": "high", "projectId": 1, "taskType": "TO DO"}`,
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "1" + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			handler := test.mockBehaviour(c, test.taskData)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.POST(create, handler.createTask)

			req := httptest.NewRequest(http.MethodPost, create, strings.NewReader(test.taskDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.createTask(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}
