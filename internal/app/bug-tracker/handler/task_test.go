package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func Test_createTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		taskData           *dto.CreateTaskDto
		taskDataJSON       string
		userData           *services.TokenData
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid user data",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *Handler {
				return &Handler{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid json",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			taskData:           nil,
			taskDataJSON:       `{"invalid"}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid create task data",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			taskData:           nil,
			taskDataJSON:       `{"name": "N"}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidTaskData.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot create task",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().CreateTask(taskData, userID).Return(uint64(0), err)

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
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().CreateTask(taskData, userID).Return(uint64(1), nil)

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
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "1" + "\n",
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

			handler := test.mockBehaviour(c, test.taskData, userID)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.POST(create, handler.createTask)

			req := httptest.NewRequest(http.MethodPost, create, strings.NewReader(test.taskDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
			echoCtx.Set(userDataCtx, test.userData)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.createTask(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_workOnTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		workOnTaskData     *dto.WorkOnTaskDto
		workOnTaskDataJSON string
		userData           *services.TokenData
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid user data",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler {
				return &Handler{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid json",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			workOnTaskData:     nil,
			workOnTaskDataJSON: `{"invalid"}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid work on task data",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			workOnTaskData:     nil,
			workOnTaskDataJSON: `{"taskId": 1}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidTaskData.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot work on task",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().WorkOnTask(workOnTaskData, userID).Return(err)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, nil}
			},
			workOnTaskData: &dto.WorkOnTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
			workOnTaskDataJSON: `{"taskId": 1, "projectId": 1}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().WorkOnTask(workOnTaskData, userID).Return(nil)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, nil}
			},
			workOnTaskData: &dto.WorkOnTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
			workOnTaskDataJSON: `{"taskId": 1, "projectId": 1}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "true" + "\n",
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

			handler := test.mockBehaviour(c, test.workOnTaskData, userID)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.POST(workOnTask, handler.workOnTask)

			req := httptest.NewRequest(http.MethodPost, workOnTask, strings.NewReader(test.workOnTaskDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
			echoCtx.Set(userDataCtx, test.userData)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.workOnTask(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_stopWorkOnTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		workOnTaskData     *dto.WorkOnTaskDto
		workOnTaskDataJSON string
		userData           *services.TokenData
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid user data",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler {
				return &Handler{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid json",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			workOnTaskData:     nil,
			workOnTaskDataJSON: `{"invalid"}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid work on task data",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			workOnTaskData:     nil,
			workOnTaskDataJSON: `{"taskId": 1}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidTaskData.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot stop work on task",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().StopWorkOnTask(workOnTaskData, userID).Return(err)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, nil}
			},
			workOnTaskData: &dto.WorkOnTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
			workOnTaskDataJSON: `{"taskId": 1, "projectId": 1}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().StopWorkOnTask(workOnTaskData, userID).Return(nil)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, nil}
			},
			workOnTaskData: &dto.WorkOnTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
			workOnTaskDataJSON: `{"taskId": 1, "projectId": 1}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "true" + "\n",
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

			handler := test.mockBehaviour(c, test.workOnTaskData, userID)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.POST(stopWorkOnTask, handler.stopWorkOnTask)

			req := httptest.NewRequest(http.MethodPost, stopWorkOnTask, strings.NewReader(test.workOnTaskDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
			echoCtx.Set(userDataCtx, test.userData)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.stopWorkOnTask(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_updateTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, taskData *dto.UpdateTaskDto, userID uint64) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		taskData           *dto.UpdateTaskDto
		taskDataJSON       string
		userData           *services.TokenData
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid user data",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.UpdateTaskDto, userID uint64) *Handler {
				return &Handler{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid json",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.UpdateTaskDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			taskData:           nil,
			taskDataJSON:       `{"invalid"}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid update task data",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.UpdateTaskDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			taskData:           nil,
			taskDataJSON:       `{"name": "N"}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidTaskData.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot update task",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.UpdateTaskDto, userID uint64) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().UpdateTask(taskData, userID).Return(uint64(0), err)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, nil}
			},
			taskData: &dto.UpdateTaskDto{
				Name:         "name",
				Description:  "description",
				TaskPriority: "high",
				TaskID:       1,
				ProjectID:    1,
				TaskType:     "TO DO",
			},
			taskDataJSON:       `{"name": "name", "description": "description", "taskPriority": "high", "taskId": 1, "projectId": 1, "taskType": "TO DO"}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.UpdateTaskDto, userID uint64) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().UpdateTask(taskData, userID).Return(uint64(1), nil)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, nil}
			},
			taskData: &dto.UpdateTaskDto{
				Name:         "name",
				Description:  "description",
				TaskPriority: "high",
				TaskID:       1,
				ProjectID:    1,
				TaskType:     "TO DO",
			},
			taskDataJSON:       `{"name": "name", "description": "description", "taskPriority": "high", "taskId": 1, "projectId": 1, "taskType": "TO DO"}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "1" + "\n",
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

			handler := test.mockBehaviour(c, test.taskData, userID)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.PUT(update, handler.updateTask)

			req := httptest.NewRequest(http.MethodPost, update, strings.NewReader(test.taskDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
			echoCtx.Set(userDataCtx, test.userData)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.updateTask(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_getTaskById(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, id uint64, ctx echo.Context) *Handler
	err := errors.New("error")
	successReturnBody := `{"id":1,"name":"name","description":"description","priority":"high","projectId":1,"taskType":"TO DO","assignee":{"Int64":1,"Valid":true},"createdAt":{"Time":"1111-11-11T11:11:11Z","Valid":true},"performTo":{"Time":"1111-11-11T11:11:11Z","Valid":true}}` + "\n"

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
			name: "Error cannot get task",
			mockBehaviour: func(c *gomock.Controller, id uint64, ctx echo.Context) *Handler {
				task := mock_services.NewMockTask(c)
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(id, nil)

				task.EXPECT().GetTaskById(id).Return(nil, err)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, params}
			},
			id:                 1,
			paramId:            "1",
			expectedStatusCode: http.StatusNotFound,
			expectedReturnBody: `{"message":"` + errTaskNotFound.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, id uint64, ctx echo.Context) *Handler {
				task := mock_services.NewMockTask(c)
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(id, nil)

				task.EXPECT().GetTaskById(id).Return(&models.Task{
					ID:          1,
					Name:        "name",
					Description: "description",
					Priority:    "high",
					ProjectID:   1,
					TaskType:    "TO DO",
					Assignee:    sql.NullInt64{Int64: 1, Valid: true},
					CreatedAt: sql.NullTime{
						Time:  time.Date(1111, 11, 11, 11, 11, 11, 0, time.UTC),
						Valid: true,
					},
					PerformTo: sql.NullTime{
						Time:  time.Date(1111, 11, 11, 11, 11, 11, 0, time.UTC),
						Valid: true,
					},
				},
					nil,
				)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, params}
			},
			id:                 1,
			paramId:            "1",
			expectedStatusCode: http.StatusFound,
			expectedReturnBody: successReturnBody,
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
			e.GET(id, handler.getTaskById)

			echoCtx.SetPath("/:id")
			echoCtx.SetParamNames("id")
			echoCtx.SetParamValues(test.paramId)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.getTaskById(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_deleteTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, taskData *dto.DeleteTaskDto, userID uint64) *Handler
	err := errors.New("error")

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		taskData           *dto.DeleteTaskDto
		taskDataJSON       string
		userData           *services.TokenData
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid user data",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.DeleteTaskDto, userID uint64) *Handler {
				return &Handler{}
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid json",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.DeleteTaskDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			taskData:           nil,
			taskDataJSON:       `{"invalid"}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid delete task data",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.DeleteTaskDto, userID uint64) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil, nil}
			},
			taskData:           nil,
			taskDataJSON:       `{"taskId": 1}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidTaskData.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot delete task",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.DeleteTaskDto, userID uint64) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().DeleteTask(taskData, userID).Return(err)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, nil}
			},
			taskData: &dto.DeleteTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
			taskDataJSON:       `{"taskId": 1, "projectId": 1}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + err.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.DeleteTaskDto, userID uint64) *Handler {
				task := mock_services.NewMockTask(c)

				task.EXPECT().DeleteTask(taskData, userID).Return(nil)

				serv := &services.Service{Task: task}

				return &Handler{serv, nil, nil, nil}
			},
			taskData: &dto.DeleteTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
			taskDataJSON:       `{"taskId": 1, "projectId": 1}`,
			userData:           &services.TokenData{UserID: 1},
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "true" + "\n",
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

			handler := test.mockBehaviour(c, test.taskData, userID)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.DELETE(empty, handler.deleteTask)

			req := httptest.NewRequest(http.MethodDelete, empty, strings.NewReader(test.taskDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
			echoCtx.Set(userDataCtx, test.userData)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.deleteTask(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}
