package container

import (
	"github.com/yuriygr/go-loggy"
	"github.com/yuriygr/go-parser/services"
)

// Container - Структура с зависимостями.
// Знаю, знаю. Но не я такой, жизнь такая.
type Container struct {
	Client  *services.Client
	Config  *services.Config
	Logger  *loggy.Logger
	Storage *services.Storage
}

// NewContainer - Собираем зависимости
func NewContainer() *Container {
	config := services.NewConfig()
	client := services.NewClient()
	logger := loggy.NewLogger()
	storage := services.NewStorage(config)

	return &Container{client, config, logger, storage}
}
