package app

import (
	"github.com/samuraivf/bug-tracker/configs"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/handler"
)

func Run() {
	configs.Init()
	handler.CreateServer()
}
