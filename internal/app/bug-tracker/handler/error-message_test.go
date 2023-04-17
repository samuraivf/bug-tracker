package handler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newErrorMessage(t *testing.T) {
	str := "error something went wrong"
	msg := newErrorMessage(errors.New(str))

	require.Equal(t, &errorMessage{Message: str}, msg)
}
