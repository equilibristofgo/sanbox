package main

import (
	"fmt"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Fenix struct {
	Field1 string
}

type testCase struct {
	body           string
	bodyResponseOk string
	bodyResponseKo string
}

var samples = map[string]testCase{
	"one": {
		body: `<?xml version='1.0'?>
		<methodCall>
		<methodName>method1</methodName>
		<params>
		<param>
		<value><string>123456789</string></value>
		</param>
		</params>
		</methodCall>`,
		bodyResponseOk: `<?xml version='1.0'?>
		<methodResponse>
		<params>
		<param>
		<value><int>4517</int></value>
		</param>
		</params>
		</methodResponse>`,
		bodyResponseKo: `<?xml version='1.0'?>
		<methodResponse>
		<fault>
		<value>
		<struct>
		<member>
		<name>faultCode</name>
		<value><int>99</int></value>
		</member>
		<member>
		<name>faultString</name>
		<value><string>Erro</string></value>
		</member>
		</struct>
		</value>
		</fault>
		</methodResponse>`,
	}}

func main() {
	port := os.Args[2]
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.POST("/", func(c echo.Context) error {
		body, _ := io.ReadAll(c.Request().Body)
		sBody := string(body)
		fmt.Println(sBody)
		return c.String(200, samples["one"].bodyResponseOk)
	})

	e.Logger.Fatal(e.Start(port))
}
