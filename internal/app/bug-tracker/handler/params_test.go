package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func Test_GetIdParam(t *testing.T) {
	tests := []struct {
		name           string
		paramId        string
		expectedResult uint64
		expectedError  error
	}{
		{
			name: "Error empty param",
			paramId: "",
			expectedResult: 0,
			expectedError: errInvalidParam,
		},
		{
			name: "Error invalid param",
			paramId: "1b",
			expectedResult: 0,
			expectedError: errInvalidParam,
		},
		{
			name: "OK",
			paramId: "1",
			expectedResult: 1,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			defer e.Close()

			validator := validator.New()
			e.Validator = newValidator(validator)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			echoCtx := e.NewContext(req, rec)
			echoCtx.SetPath("/:id")
			echoCtx.SetParamNames("id")
			echoCtx.SetParamValues(test.paramId)

			defer rec.Result().Body.Close()
			req.Close = true

			p := &params{}
			id, err := p.GetIdParam(echoCtx)

			require.Equal(t, test.expectedResult, id)
			require.Equal(t, test.expectedError, err)
		})
	}
}
