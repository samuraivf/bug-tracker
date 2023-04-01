package handler

type errorMessage struct {
	Message string `json:"message"`
}

func newErrorMessage(err error) *errorMessage {
	return &errorMessage{Message: err.Error()}
}
