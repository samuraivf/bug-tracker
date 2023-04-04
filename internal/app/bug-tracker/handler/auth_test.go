package handler

import (
	"context"
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

func TestSignUp(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	ctx := context.Background()
	user := mock_services.NewMockUser(c)
	redis := mock_services.NewMockRedis(c)
	log := mock_log.NewMockLog(c)

	testUserDataJSON := `{
		"name": "Random",
		"username": "random",
		"email": "randommail@gmail.com",
		"password": "password"
	}`

	testUserData := &dto.SignUpDto{
		Name:     "Random",
		Username: "random",
		Email:    "randommail@gmail.com",
		Password: "password",
	}

	user.EXPECT().GetUserByEmail(testUserData.Email).Return(nil, errors.New("error user not found"))
	user.EXPECT().GetUserByUsername(testUserData.Username).Return(nil, errors.New("error user not found"))
	redis.EXPECT().Get(ctx, testUserData.Email).Return("verified", nil)
	user.EXPECT().CreateUser(testUserData).Return(uint64(1), nil)

	serv := &services.Service{User: user, Redis: redis}
	handler := &Handler{serv, log, nil}

	e := echo.New()
	defer e.Close()

	validator := validator.New()
	e.Validator = newValidator(validator)
	e.POST(signUp, handler.signUp)

	req := httptest.NewRequest(http.MethodPost, signUp, strings.NewReader(testUserDataJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)

	defer rec.Result().Body.Close()
	req.Close = true

	require.NoError(t, handler.signUp(echoCtx))
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Equal(t, "1\n", rec.Body.String())
}

func TestSignUpError(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, userData *dto.SignUpDto) *Handler

	tests := []struct{
		name string
		mockBehaviour mockBehaviour
		userData *dto.SignUpDto
		userDataJSON string
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid json",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignUpDto) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil}
			},
			userData: nil,
			userDataJSON: `{"invalid"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			handler := test.mockBehaviour(c, test.userData)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.POST(signUp, handler.signUp)

			req := httptest.NewRequest(http.MethodPost, signUp, strings.NewReader(test.userDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.signUp(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}
