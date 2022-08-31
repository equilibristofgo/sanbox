package models

import (
	"reflect"
	"strconv"
)

type HttpHandler struct {
	Type     string `yaml:"type" default:"restful"`
	Url      string `yaml:"url" default:"127.0.0.1"`
	Port     string `yaml:"port" default:"3000"`
	Swagger  string `yaml:"swagger" default:"../docs/swagger.yaml"`
	LogLevel string `yaml:"debug_level" default:"INFO"`
}

type Auth struct {
	Enable bool `yaml:"enable" default:"false"`
}

type Micro struct {
	Name string `yaml:"name" default:"micro-template"`
}

type Repository struct {
	Db string `yaml:"db" default:"postgres"`
}

//DefaultTag - Server model
func (a HttpHandler) DefaultTag() HttpHandler {
	aAux := HttpHandler{}
	typ := reflect.TypeOf(aAux)

	if f, exist := typ.FieldByName("Type"); exist {
		aAux.Type = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Url"); exist {
		aAux.Url = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Port"); exist {
		aAux.Port = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("Swagger"); exist {
		aAux.Swagger = f.Tag.Get("default")
	}

	if f, exist := typ.FieldByName("LogLevel"); exist {
		aAux.LogLevel = f.Tag.Get("default")
	}

	return aAux
}

//DefaultTag - Auth model
func (a Auth) DefaultTag() Auth {

	aAux := Auth{}
	typ := reflect.TypeOf(aAux)

	if f, exist := typ.FieldByName("Enable"); exist {
		if value, err := strconv.ParseBool(f.Tag.Get("default")); err == nil {
			aAux.Enable = value
		}
	}

	return aAux
}

//DefaultTag - Micro model
func (a Micro) DefaultTag() Micro {

	aAux := Micro{}
	typ := reflect.TypeOf(aAux)

	if f, exist := typ.FieldByName("Name"); exist {
		aAux.Name = f.Tag.Get("default")
	}

	return aAux
}

//DefaultTag - Jaeger config
func (a Repository) DefaultTag() Repository {

	aAux := Repository{}
	typ := reflect.TypeOf(aAux)

	if f, exist := typ.FieldByName("Db"); exist {
		aAux.Db = f.Tag.Get("default")
	}

	return aAux
}
