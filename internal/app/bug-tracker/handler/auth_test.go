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
	mock_kafka "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/kafka/mocks"
	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
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

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		userData           *dto.SignUpDto
		userDataJSON       string
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
			userData:           nil,
			userDataJSON:       `{"invalid"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid sign up data",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignUpDto) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil}
			},
			userData:           nil,
			userDataJSON:       `{"name": "Name", "email": "email", "password": "password"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidSignUpData.Error() + `"}` + "\n",
		},
		{
			name: "Error user already exists",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignUpDto) *Handler {
				user := mock_services.NewMockUser(c)

				user.EXPECT().GetUserByEmail(userData.Email).Return(&models.User{}, nil)

				serv := &services.Service{User: user}

				return &Handler{serv, nil, nil}
			},
			userData: &dto.SignUpDto{
				Name:     "Name",
				Email:    "email@gmail.com",
				Password: "password",
				Username: "username",
			},
			userDataJSON:       `{"name": "Name", "email": "email@gmail.com", "password": "password", "username": "username"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserEmailAlreadyExists.Error() + `"}` + "\n",
		},
		{
			name: "Error user username exists",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignUpDto) *Handler {
				user := mock_services.NewMockUser(c)

				user.EXPECT().GetUserByEmail(userData.Email).Return(nil, errors.New("no user"))
				user.EXPECT().GetUserByUsername(userData.Username).Return(&models.User{}, nil)

				serv := &services.Service{User: user}

				return &Handler{serv, nil, nil}
			},
			userData: &dto.SignUpDto{
				Name:     "Name",
				Email:    "email@gmail.com",
				Password: "password",
				Username: "username",
			},
			userDataJSON:       `{"name": "Name", "email": "email@gmail.com", "password": "password", "username": "username"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserUsernameAlreadyExists.Error() + `"}` + "\n",
		},
		{
			name: "Error no verified email",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignUpDto) *Handler {
				user := mock_services.NewMockUser(c)
				redis := mock_services.NewMockRedis(c)
				ctx := context.Background()

				user.EXPECT().GetUserByEmail(userData.Email).Return(nil, errors.New("no user"))
				user.EXPECT().GetUserByUsername(userData.Username).Return(nil, errors.New("no user"))
				redis.EXPECT().Get(ctx, userData.Email).Return("", errors.New("no value"))

				serv := &services.Service{User: user, Redis: redis}

				return &Handler{serv, nil, nil}
			},
			userData: &dto.SignUpDto{
				Name:     "Name",
				Email:    "email@gmail.com",
				Password: "password",
				Username: "username",
			},
			userDataJSON:       `{"name": "Name", "email": "email@gmail.com", "password": "password", "username": "username"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errEmailIsNotVerified.Error() + `"}` + "\n",
		},
		{
			name: "Error cannot create user",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignUpDto) *Handler {
				user := mock_services.NewMockUser(c)
				redis := mock_services.NewMockRedis(c)
				ctx := context.Background()

				user.EXPECT().GetUserByEmail(userData.Email).Return(nil, errors.New("no user"))
				user.EXPECT().GetUserByUsername(userData.Username).Return(nil, errors.New("no user"))
				redis.EXPECT().Get(ctx, userData.Email).Return("verified", nil)
				user.EXPECT().CreateUser(userData).Return(uint64(0), errors.New("cannot create user"))

				serv := &services.Service{User: user, Redis: redis}

				return &Handler{serv, nil, nil}
			},
			userData: &dto.SignUpDto{
				Name:     "Name",
				Email:    "email@gmail.com",
				Password: "password",
				Username: "username",
			},
			userDataJSON:       `{"name": "Name", "email": "email@gmail.com", "password": "password", "username": "username"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + errInternalServerError.Error() + `"}` + "\n",
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

func TestVerifyEmail(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	log := mock_log.NewMockLog(c)
	kafka := mock_kafka.NewMockKafka(c)

	testVerifyEmailJSON := `{
		"email": "email@gmail.com"
	}`

	testVerifyEmail := &dto.VerifyEmail{
		Email: "email@gmail.com",
	}

	kafka.EXPECT().Write(testVerifyEmail.Email).Return(nil)
	log.EXPECT().Infof("[Kafka] Sent message: %s", testVerifyEmail.Email)

	handler := &Handler{nil, log, kafka}

	e := echo.New()
	defer e.Close()

	validator := validator.New()
	e.Validator = newValidator(validator)
	e.POST(verify, handler.verifyEmail)

	req := httptest.NewRequest(http.MethodPost, verify, strings.NewReader(testVerifyEmailJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)

	defer rec.Result().Body.Close()
	req.Close = true

	require.NoError(t, handler.verifyEmail(echoCtx))
	require.Equal(t, http.StatusOK, rec.Result().StatusCode)
	require.Equal(t, "null\n", rec.Body.String())
}

func TestVerifyEmailError(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, verifyEmail *dto.VerifyEmail) *Handler

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		verifyEmail        *dto.VerifyEmail
		verifyEmailJSON    string
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid json",
			mockBehaviour: func(c *gomock.Controller, verifyEmail *dto.VerifyEmail) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil}
			},
			verifyEmail: nil,
			verifyEmailJSON: `{"invalid"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error in kafka",
			mockBehaviour: func(c *gomock.Controller, verifyEmail *dto.VerifyEmail) *Handler {
				log := mock_log.NewMockLog(c)
				kafka := mock_kafka.NewMockKafka(c)

				kafka.EXPECT().Write(verifyEmail.Email).Return(errors.New("error"))
				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, kafka}
			},
			verifyEmail: &dto.VerifyEmail{Email: "email@gmail.com"},
			verifyEmailJSON: `{"email": "email@gmail.com"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + errInternalServerError.Error() + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			handler := test.mockBehaviour(c, test.verifyEmail)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.POST(verify, handler.verifyEmail)

			req := httptest.NewRequest(http.MethodPost, verify, strings.NewReader(test.verifyEmailJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.verifyEmail(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}
