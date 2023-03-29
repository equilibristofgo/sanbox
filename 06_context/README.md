# CONTEXT
In this sandbox we play different aprox to use context to transfer info between http handler/middleware

Remember ... to generate this documentation from code
- sandbox/06_context/01_example$ gomarkdoc --output doc.md .
- sandbox/06_context/02_example$ gomarkdoc --output doc.md .

## Initial example
In this simple example, show how to insert key/value inside context to store infro inter methods..

See [doc](01_example/doc.md)

## More complex example
This server using echo, use a middleware that implement a [proxy](https://echo.labstack.com/middleware/proxy/) and could be configure to [transform the result](https://github.com/labstack/echo/discussions/1992). In general, when we read the body twice (inside a middleware than could be chained with other that could need the body too) when need [regenerate](https://stackoverflow.com/questions/31535569/how-to-read-response-body-of-reverseproxy) the readen body.

In this case, as we use Echo, the [context](https://echo.labstack.com/guide/context/) is diffent from standard, this add extra difficult to overwrite the functionality. This make us to store values in underlaying context (native) through echo context [as described here](https://stackoverflow.com/questions/69326129/does-set-method-of-echo-context-saves-the-value-to-the-underlying-context-cont)

To understant how handlers can be chained [here](https://gist.github.com/husobee/fd23681261a39699ee37), [here](https://medium.com/@chrisgregory_83433/chaining-middleware-in-go-918cfbc5644d) and [here](https://fideloper.com/golang-context-http-middleware)

Some other feature of proxy middleware are left commented: url rewrite and overwrite of http transport with RourdTrip more info [here](https://lanre.wtf/blog/2017/07/24/roundtripper-go) and [here](https://echorand.me/posts/go-http-client-middleware/)

A [thoughts](https://medium.com/ymedialabs-innovation/reverse-proxy-in-go-d26482acbcad) of ReverseProxy... and finally a Google product, that born from testing but could be use in other ways [Martian Proxy](https://github.com/google/martian/blob/0f7e6797a04da412118541344bbe0d65945e24c9/cmd/proxy/main.go#L288)

The sample came with a upstream server (port 8080) where the proxy points... and where we can send info

```
curl --header "Content-Type: text/xml;charset=UTF-8" --data @sample.xml localhost:1323
```

a .vscode\launch.json are provided.

See [doc](02_example/doc.md)