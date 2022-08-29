package config

import (
	"sync"
)

var once sync.Once
var config *App1Config

type App1Config struct {
}

func GetConfig() *App1Config {
	once.Do(func() {
		config = &App1Config{}
		// TODO Load...
	})
	return config
}
