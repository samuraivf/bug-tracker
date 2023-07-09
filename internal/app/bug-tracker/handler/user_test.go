package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	mock_handler "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/handler/mocks"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
	mock_services "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services/mocks"
	"github.com/stretchr/testify/require"
)

func Test_getUserById(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, id uint64, ctx echo.Context) *Handler
	err := errors.New("error")
	successReturnBody := `{"id":1,"name":"name","username":"username","email":"email@gmail.com"}` + "\n"

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
			name: "Error cannot get user",
			mockBehaviour: func(c *gomock.Controller, id uint64, ctx echo.Context) *Handler {
				user := mock_services.NewMockUser(c)
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(id, nil)

				user.EXPECT().GetUserById(id).Return(nil, err)

				serv := &services.Service{User: user}

				return &Handler{serv, nil, nil, params}
			},
			id:                 1,
			paramId:            "1",
			expectedStatusCode: http.StatusNotFound,
			expectedReturnBody: `{"message":"` + errUserNotFound.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, id uint64, ctx echo.Context) *Handler {
				user := mock_services.NewMockUser(c)
				params := mock_handler.NewMockParams(c)

				params.EXPECT().GetIdParam(ctx).Return(id, nil)

				user.EXPECT().GetUserById(id).Return(&models.User{
					ID:       1,
					Name:     "name",
					Username: "username",
					Email:    "email@gmail.com",
				},
					nil,
				)

				serv := &services.Service{User: user}

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
			e.GET(id, handler.getUserById)

			echoCtx.SetPath("/:id")
			echoCtx.SetParamNames("id")
			echoCtx.SetParamValues(test.paramId)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.getUserById(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}
