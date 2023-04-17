package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services"
	mock_services "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/services/mocks"
)

func Test_Logger(t *testing.T) {
	c := gomock.NewController(t)
	log := mock_log.NewMockLog(c)

	middlewareFunc := Logger(log)
	
	e := echo.New()
	e.Use(middlewareFunc)

	require.NotNil(t, e.Logger)
}

func Test_isUnauthorized(t *testing.T) {
	tests := []struct{
		name string
		authorizationHeader bool
		refreshTokenCookie *http.Cookie
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "Authorization Header exists",
			authorizationHeader: true,
			refreshTokenCookie: &http.Cookie{},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserIsAuthorized.Error() + `"}` + "\n",
		},
		{
			name: "Refresh token exists",
			authorizationHeader: false,
			refreshTokenCookie: &http.Cookie{
				Name:     "refreshToken",
				Value:    "token",
				Expires:  time.Now().Add(time.Minute * 10),
				HttpOnly: true,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedReturnBody: `{"message":"` + errUserIsAuthorized.Error() + `"}` + "\n",
		},
		{
			name: "OK",
			authorizationHeader: false,
			refreshTokenCookie: &http.Cookie{},
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: `null` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := NewHandler(nil, nil, nil)
			e := echo.New()
			middleware := h.isUnauthorized(func(c echo.Context) error {
				return c.JSON(http.StatusOK, nil)
			})
		
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.AddCookie(test.refreshTokenCookie)

			if test.authorizationHeader {
				req.Header.Set(authorizationHeader, "Bearer")
			}

			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
		
			defer rec.Result().Body.Close()
			req.Close = true
		
			require.NoError(t, middleware(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_isAuthorized(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, token string) *Handler

	tests := []struct{
		name string
		authorizationHeader bool
		mockBehaviour mockBehaviour
		token string
		expectedStatusCode int
		expectedReturnBody string
	}{
		{
			name: "No Authorization Header",
			authorizationHeader: false,
			mockBehaviour: func(c *gomock.Controller, token string) *Handler {
				return &Handler{}
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedReturnBody: `{"message":"` + errInvalidAuthHeader.Error() + `"}` + "\n",
		},
		{
			name: "Invalid Authorization header [Not Bearer]",
			authorizationHeader: true,
			mockBehaviour: func(c *gomock.Controller, token string) *Handler {
				return &Handler{}
			},
			token: "Berer token",
			expectedStatusCode: http.StatusUnauthorized,
			expectedReturnBody: `{"message":"` + errInvalidAuthHeader.Error() + `"}` + "\n",
		},
		{
			name: "Invalid Authorization header [Parts Length<2]",
			authorizationHeader: true,
			mockBehaviour: func(c *gomock.Controller, token string) *Handler {
				return &Handler{}
			},
			token: "Bearer",
			expectedStatusCode: http.StatusUnauthorized,
			expectedReturnBody: `{"message":"` + errInvalidAuthHeader.Error() + `"}` + "\n",
		},
		{
			name: "Invalid Authorization header [Parts Length>2]",
			authorizationHeader: true,
			mockBehaviour: func(c *gomock.Controller, token string) *Handler {
				return &Handler{}
			},
			token: "Bearer token token",
			expectedStatusCode: http.StatusUnauthorized,
			expectedReturnBody: `{"message":"` + errInvalidAuthHeader.Error() + `"}` + "\n",
		},
		{
			name: "Invalid Authorization header [Parts[1] Length == 0]",
			authorizationHeader: true,
			mockBehaviour: func(c *gomock.Controller, token string) *Handler {
				return &Handler{}
			},
			token: "Bearer ",
			expectedStatusCode: http.StatusUnauthorized,
			expectedReturnBody: `{"message":"` + errTokenIsEmpty.Error() + `"}` + "\n",
		},
		{
			name: "Error in ParseAccessToken",
			authorizationHeader: true,
			mockBehaviour: func(c *gomock.Controller, token string) *Handler {
				auth := mock_services.NewMockAuth(c)

				headerParts := strings.Split(token, " ")

				auth.EXPECT().ParseAccessToken(headerParts[1]).Return(nil, errors.New("error"))
				
				return &Handler{service: &services.Service{Auth: auth}}
			},
			token: "Bearer token",
			expectedStatusCode: http.StatusUnauthorized,
			expectedReturnBody: `{"message":"error"}` + "\n",
		},
		{
			name: "OK",
			authorizationHeader: true,
			mockBehaviour: func(c *gomock.Controller, token string) *Handler {
				auth := mock_services.NewMockAuth(c)

				headerParts := strings.Split(token, " ")

				auth.EXPECT().ParseAccessToken(headerParts[1]).Return(&services.TokenData{}, nil)
				
				return &Handler{service: &services.Service{Auth: auth}}
			},
			token: "Bearer token",
			expectedStatusCode: http.StatusOK,
			expectedReturnBody: "null" + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			h := test.mockBehaviour(c, test.token)
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			if test.authorizationHeader {
				req.Header.Set(authorizationHeader, test.token)
			}

			middleware := h.isAuthorized(func(c echo.Context) error {
				return c.JSON(http.StatusOK, nil)
			})

			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
		
			defer rec.Result().Body.Close()
			req.Close = true
		
			require.NoError(t, middleware(echoCtx))
			require.Equal(t, test.expectedStatusCode, echoCtx.Response().Status)
			require.Equal(t, test.expectedReturnBody, rec.Body.String())
		})
	}
}

func Test_getUserData(t *testing.T) {
	tests := []struct{
		name string
		userData interface{}
		expectedUserData *services.TokenData
		expectedError error
	}{
		{
			name: "No userData",
			userData: nil,
			expectedUserData: nil,
			expectedError: errUserNotFound,
		},
		{
			name: "UserData invalid type",
			userData: "userData",
			expectedUserData: nil,
			expectedError: errUserDataInvalidType,
		},
		{
			name: "OK",
			userData: &services.TokenData{
				TokenID: "id",
				UserID: uint64(1),
				Username: "username",
			},
			expectedUserData: &services.TokenData{
				TokenID: "id",
				UserID: uint64(1),
				Username: "username",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)

			echoCtx.Set(userDataCtx, test.userData)

			userData, err := getUserData(echoCtx)

			require.Equal(t, test.expectedUserData, userData)
			require.Equal(t, test.expectedError, err)
		})
	}
}