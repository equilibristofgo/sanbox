// This server using echo, use a middleware that implement a proxy and could be configure to transform the result
// Also, show how to implement a context (that extend echo) tha can store the body of request to retrive insede middleware
// Some other feature of proxy middleware are left commented: url rewrite and overwrite of http transport with RourdTrip
package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type SomeClient struct{}

func (t *SomeClient) RoundTrip(r *http.Request) (*http.Response, error) {
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))

	r.Body = rdr1

	response, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		print("\n\ncame in error resp here", err)
		return nil, err //Server is not reachable. Server not working
	}

	return response, err
}

type proxyContext struct {
	echo.Context
}

// Get retrieves data from the context.
func (ctx proxyContext) Get(key string) interface{} {
	// get old context value
	val := ctx.Context.Get(key)
	if val != nil {
		return val
	}
	return ctx.Request().Context().Value(key)
}

// Set saves data in the context.
func (ctx proxyContext) Set(key string, val interface{}) {
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), key, val)))
}

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Setup proxy
	url1, err := url.Parse("http://localhost:8081")
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Store request body to extract info in response phase...
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &proxyContext{c}
			buf, _ := ioutil.ReadAll(c.Request().Body)
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			c.Request().Body = rdr1
			cc.Set("bodyRequest", string(buf))
			return next(cc)
		}
	})

	// Proxy middleware
	e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{

		// Transport: &SomeClient{},

		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{
				Name: "proxy",
				URL:  url1,
			},
		}),

		ModifyResponse: func(resp *http.Response) error {
			ctx := resp.Request.Context()
			originalBody := ctx.Value("bodyRequest").(string)

			// Generate new body...
			bodyChanged := originalBody

			body := ioutil.NopCloser(bytes.NewReader([]byte(bodyChanged)))
			resp.Body = body
			resp.ContentLength = int64(len(bodyChanged))
			resp.Header.Set("Content-Length", strconv.Itoa(len(bodyChanged)))

			return nil
		},

		// Rewrite: map[string]string{
		// 	"^/v1/*": "/v2/$1",
		// },

		// RegexRewrite: map[*regexp.Regexp]string{
		// 	regexp.MustCompile("^/foo/([0-9].*)"):  "/num/$1",
		// 	regexp.MustCompile("^/bar/(.+?)/(.*)"): "/baz/$2/$1",
		// },
	}))

	e.Logger.Fatal(e.Start(":1323"))

}
