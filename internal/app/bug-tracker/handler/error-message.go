package handler

type errorMessage struct {
	Message string `json:"message"`
}

func newErrorMessage(message string) *errorMessage {
	return &errorMessage{message}
}
