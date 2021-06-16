package services

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

// NewConfig - Формируем конфиг
func NewConfig(path string) (config *Config) {
	if _, err := toml.DecodeFile(fmt.Sprintf("%sconfig.toml", path), &config); err != nil {
		log.Fatalln(err)
	}
	config.Application.Path = path
	return config
}

// Config - Конфигурационный файл
type Config struct {
	Application struct {
		FirstMessage string
		Status       string
		Path         string
	}

	Storage struct {
		Driver       string
		DSN          string
		MaxOpenConns int
	}
}
