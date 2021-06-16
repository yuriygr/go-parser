package container

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/yuriygr/go-loggy"
	"github.com/yuriygr/go-parser/services"
)

// Container - Структура с зависимостями.
// Знаю, знаю. Но не я такой, жизнь такая.
type Container struct {
	Client  *services.Client
	Config  *services.Config
	Logger  *loggy.Logger
	Storage *sqlx.DB
}

// NewContainer - Собираем зависимости
func NewContainer() *Container {
	config := services.NewConfig(os.Getenv("GORKI_PATH"))

	client := services.NewClient()
	logger := loggy.NewLogger(loggy.LoggerConfig{})
	storage := services.NewStorage(config, logger)

	return &Container{client, config, logger, storage}
}
