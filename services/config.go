package services

import (
	"os"
)

// NewConfig - Формируем конфиг
func NewConfig() *Config {
	config := &Config{}
	config.Storage.DSN = os.Getenv("DB_DSN")

	return config
}

// Config - Конфигурационный файл
type Config struct {
	// TODO: Исправить с DSN на что-то более понятное
	Storage struct {
		DSN string
	}
}
