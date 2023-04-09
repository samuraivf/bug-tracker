package handler

import (
	"context"
	"errors"
	"fmt"
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
	mock_kafka "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/kafka/mocks"
	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
	mock_services "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services/mocks"
)

func TestSignUp(t *testing.T) {
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
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignUpDto) *Handler {
				user := mock_services.NewMockUser(c)
				redis := mock_services.NewMockRedis(c)
				ctx := context.Background()

				user.EXPECT().GetUserByEmail(userData.Email).Return(nil, errors.New("no user"))
				user.EXPECT().GetUserByUsername(userData.Username).Return(nil, errors.New("no user"))
				redis.EXPECT().Get(ctx, userData.Email).Return("verified", nil)
				user.EXPECT().CreateUser(userData).Return(uint64(1), nil)

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
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "1" + "\n",
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
			verifyEmail:        nil,
			verifyEmailJSON:    `{"invalid"}`,
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
			verifyEmail:        &dto.VerifyEmail{Email: "email@gmail.com"},
			verifyEmailJSON:    `{"email": "email@gmail.com"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + errInternalServerError.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, verifyEmail *dto.VerifyEmail) *Handler {
				log := mock_log.NewMockLog(c)
				kafka := mock_kafka.NewMockKafka(c)

				kafka.EXPECT().Write(verifyEmail.Email).Return(nil)
				log.EXPECT().Infof("[Kafka] Sent message: %s", verifyEmail.Email)

				return &Handler{nil, log, kafka}
			},
			verifyEmailJSON:    `{"email": "email@gmail.com"}`,
			verifyEmail:        &dto.VerifyEmail{Email: "email@gmail.com"},
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "null" + "\n",
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

func TestSetEmail(t *testing.T) {
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
			verifyEmail:        nil,
			verifyEmailJSON:    `{"invalid"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidJSON.Error() + `"}` + "\n",
		},
		{
			name: "Error in redis",
			mockBehaviour: func(c *gomock.Controller, verifyEmail *dto.VerifyEmail) *Handler {
				log := mock_log.NewMockLog(c)
				redis := mock_services.NewMockRedis(c)
				ctx := context.Background()

				redis.EXPECT().Set(ctx, verifyEmail.Email, "verified", time.Minute*10).Return(errors.New("error"))

				return &Handler{&services.Service{Redis: redis}, log, nil}
			},
			verifyEmail:        &dto.VerifyEmail{Email: "email@gmail.com"},
			verifyEmailJSON:    `{"email": "email@gmail.com"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedReturnBody: `{"message":"` + errInternalServerError.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, verifyEmail *dto.VerifyEmail) *Handler {
				redis := mock_services.NewMockRedis(c)
				ctx := context.Background()

				redis.EXPECT().Set(ctx, verifyEmail.Email, "verified", time.Minute*10).Return(nil)

				return &Handler{&services.Service{Redis: redis}, nil, nil}
			},
			verifyEmail:        &dto.VerifyEmail{Email: "email@gmail.com"},
			verifyEmailJSON:    `{"email": "email@gmail.com"}`,
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "null" + "\n",
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
			e.POST(setEmail, handler.setEmail)

			req := httptest.NewRequest(http.MethodPost, setEmail, strings.NewReader(test.verifyEmailJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.setEmail(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func TestSignIn(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, userData *dto.SignInDto) *Handler

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		userData           *dto.SignInDto
		userDataJSON       string
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Error invalid json",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignInDto) *Handler {
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
			name: "Error invalid sign in data",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignInDto) *Handler {
				log := mock_log.NewMockLog(c)

				log.EXPECT().Error(gomock.Any()).Return()

				return &Handler{nil, log, nil}
			},
			userData:           nil,
			userDataJSON:       `{"email": "email", "password": "password"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidSignInData.Error() + `"}` + "\n",
		},
		{
			name: "Error invalid user",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignInDto) *Handler {
				user := mock_services.NewMockUser(c)
				log := mock_log.NewMockLog(c)

				user.EXPECT().ValidateUser(userData.Email, userData.Password).Return(nil, errors.New("error"))
				log.EXPECT().Error(gomock.Any()).Return()

				serv := &services.Service{User: user}

				return &Handler{serv, log, nil}
			},
			userData: &dto.SignInDto{
				Email:    "email@gmail.com",
				Password: "password",
			},
			userDataJSON:       `{"email": "email@gmail.com", "password": "password"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"error"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignInDto) *Handler {
				user := mock_services.NewMockUser(c)

				user.EXPECT().ValidateUser(userData.Email, userData.Password).Return(&models.User{
					Username: "username",
					ID:       uint64(1),
				}, nil)

				serv := &services.Service{User: user}

				return &Handler{serv, nil, nil}
			},
			userData: &dto.SignInDto{
				Email:    "email@gmail.com",
				Password: "password",
			},
			userDataJSON:       `{"email": "email@gmail.com", "password": "password"}`,
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "null" + "\n",
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
			e.POST(signIn, func(c echo.Context) error {
				return handler.signIn(c, func(c echo.Context, username string, userID uint64) error {
					return c.JSON(http.StatusOK, nil)
				})
			})

			req := httptest.NewRequest(http.MethodPost, signIn, strings.NewReader(test.userDataJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.signIn(echoCtx, func(c echo.Context, username string, userID uint64) error {
				return c.JSON(http.StatusOK, nil)
			}))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func TestRefresh(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, refreshToken string) *Handler

	tests := []struct {
		name               string
		mockBehaviour      mockBehaviour
		refreshToken       string
		refreshTokenCookie *http.Cookie
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "No refresh token",
			mockBehaviour: func(c *gomock.Controller, refreshToken string) *Handler {
				return &Handler{nil, nil, nil}
			},
			refreshToken:       "",
			refreshTokenCookie: &http.Cookie{},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errInvalidRefreshToken.Error() + `"}` + "\n",
		},
		{
			name: "Invalid refresh token",
			mockBehaviour: func(c *gomock.Controller, refreshToken string) *Handler {
				auth := mock_services.NewMockAuth(c)

				auth.EXPECT().ParseRefreshToken(refreshToken).Return(nil, errors.New("error"))

				return &Handler{&services.Service{Auth: auth}, nil, nil}
			},
			refreshToken: "token",
			refreshTokenCookie: &http.Cookie{
				Name:     "refreshToken",
				Value:    "token",
				Expires:  time.Now().Add(time.Minute * 10),
				HttpOnly: true,
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedReturnBody: `{"message":"error"}` + "\n",
		},
		{
			name: "Error in redis",
			mockBehaviour: func(c *gomock.Controller, refreshToken string) *Handler {
				auth := mock_services.NewMockAuth(c)
				redis := mock_services.NewMockRedis(c)
				ctx := context.Background()

				tokenData := &services.TokenData{
					TokenID:  "id",
					Username: "username",
					UserID:   1,
				}
				key := fmt.Sprintf("%s:%s", tokenData.Username, tokenData.TokenID)

				auth.EXPECT().ParseRefreshToken(refreshToken).Return(tokenData, nil)
				redis.EXPECT().GetRefreshToken(ctx, key).Return("", errors.New("error"))

				return &Handler{&services.Service{Auth: auth, Redis: redis}, nil, nil}
			},
			refreshToken: "token",
			refreshTokenCookie: &http.Cookie{
				Name:     "refreshToken",
				Value:    "token",
				Expires:  time.Now().Add(time.Minute * 10),
				HttpOnly: true,
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedReturnBody: `{"message":"` + errTokenDoesNotExist.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, refreshToken string) *Handler {
				auth := mock_services.NewMockAuth(c)
				redis := mock_services.NewMockRedis(c)
				ctx := context.Background()

				tokenData := &services.TokenData{
					TokenID:  "id",
					Username: "username",
					UserID:   1,
				}
				key := fmt.Sprintf("%s:%s", tokenData.Username, tokenData.TokenID)

				auth.EXPECT().ParseRefreshToken(refreshToken).Return(tokenData, nil)
				redis.EXPECT().GetRefreshToken(ctx, key).Return("token", nil)

				return &Handler{&services.Service{Auth: auth, Redis: redis}, nil, nil}
			},
			refreshToken: "token",
			refreshTokenCookie: &http.Cookie{
				Name:     "refreshToken",
				Value:    "token",
				Expires:  time.Now().Add(time.Minute * 10),
				HttpOnly: true,
			},
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "null" + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			handler := test.mockBehaviour(c, test.refreshToken)

			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)
			e.GET(refresh, func(c echo.Context) error {
				return handler.refresh(c, func(c echo.Context, username string, userID uint64) error {
					return c.JSON(http.StatusOK, nil)
				})
			})

			req := httptest.NewRequest(http.MethodGet, refresh, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			req.AddCookie(test.refreshTokenCookie)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)

			defer rec.Result().Body.Close()
			req.Close = true

			require.NoError(t, handler.refresh(echoCtx, func(c echo.Context, username string, userID uint64) error {
				return c.JSON(http.StatusOK, nil)
			}))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}
