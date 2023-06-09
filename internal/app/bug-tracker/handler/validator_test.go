package handler

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mock_handler "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/handler/mocks"
)

func Test_newValidator(t *testing.T) {
	validator := validator.New()
	customValidator := newValidator(validator)

	require.Equal(t, &CustomValidator{validator}, customValidator)
}

func Test_Validate(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	mockValidator := mock_handler.NewMockValidator(c)
	customValidator := newValidator(mockValidator)
	s := struct{}{}
	mockValidator.EXPECT().Struct(s).Return(nil)

	require.NoError(t, customValidator.Validate(s))
}

func Test_ValidateError(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	mockValidator := mock_handler.NewMockValidator(c)
	customValidator := newValidator(mockValidator)
	s := struct{}{}
	err := errors.New("error something went wrong")
	mockValidator.EXPECT().Struct(s).Return(err)

	require.EqualError(t, customValidator.Validate(s), err.Error())
}
