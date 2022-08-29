package config

import (
	"sync"
)

var once sync.Once
var config *App2Config

type App2Config struct {
}

func GetConfig() *App2Config {
	once.Do(func() {
		config = &App2Config{}
		// TODO Load...
	})
	return config
}
