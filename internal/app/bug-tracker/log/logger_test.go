package log

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type Writer struct {
	M map[string]string
}

func (w *Writer) Write(p []byte) (n int, err error) {
	splitted := strings.Split(string(p), " ")
	msg := splitted[len(splitted)-1]
	msg = stripansi.Strip(msg)
	fmt.Println(msg)
	w.M[msg] = msg
	return len(p), nil
}

func Test_New(t *testing.T) {
	log := zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()

	expected := &Logger{
		log: &log,
	}

	require.Equal(t, expected, New(os.Stderr))
}

func Test_Info(t *testing.T) {
	out := &Writer{M: map[string]string{}}
	log := New(out)

	msg := "Info"
	log.Info(msg)

	require.Equal(t, msg+"\n", out.M[msg+"\n"])
}

func Test_Infof(t *testing.T) {
	out := &Writer{M: map[string]string{}}
	log := New(out)

	msg := "Info%s"
	log.Infof(msg, "info")

	expected := fmt.Sprintf(msg, "info") + "\n"

	require.Equal(t, expected, out.M[expected])
}

func Test_Error(t *testing.T) {
	out := &Writer{M: map[string]string{}}
	log := New(out)

	msg := errors.New("Error")
	log.Error(msg)

	expected := "error=" + msg.Error() + "\n"

	require.Equal(t, expected, out.M[expected])
}

func Test_Internal(t *testing.T) {
	log := New(os.Stderr)

	require.NotNil(t, log.Internal())
	require.IsType(t, reflect.TypeOf(zerolog.New(os.Stderr)), reflect.TypeOf(log.Internal()))
}
