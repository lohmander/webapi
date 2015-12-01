# webapi

[![Build Status](https://travis-ci.org/lohmander/webapi.svg)](https://travis-ci.org/lohmander/webapi)
[![GoDoc](https://godoc.org/github.com/lohmander/webapi?status.svg)](https://godoc.org/github.com/lohmander/webapi)

Minimalist web API framework written in Go. Inspired by [`sleepy`](https://github.com/dougblack/sleepy), but with some batteries included. No external dependencies.

- [Usage](#usage)
    + [Endpoints](#endpoints)
    + [Middleware](#middleware)
    + [URL parameters](#url-parameters)
    + [Request body](#request-body)
- [Example](#example)
- [License](#license)



## Usage

To create a new api, simply use `webapi.NewAPI` and use the returned API as a handler.

```go
api := webapi.NewAPI()
http.ListenAndServe(":3002", api)
```

Or prefix your API by passing it to `http.Handle` (notice the trailing slash).

```go
api := webapi.NewAPI()
http.Handle("/api/", api)
http.ListenAndServe(":3002", nil)
```

### Endpoints

```go
api.Add(`/items/(?P<id>\d+)$`, &Item{})

// ...

type Item struct {}

func (item Item) Get(r *webapi.Request) (int, webapi.Response) {
    someData := map[string]string{
        "key": "value"
    }
    return 200, webapi.Response{
        Data: someData
    }
}
```

#### Multiple handlers for one endpoint

Having multiple handlers for the same endpoint can be useful for instance when you need to maintain multiple versions of the same endpoint in your API.

It can be done using the `webapi.Handlers` function. 

```go
func create_v100(r *webapi.Request) (int, webapi.Response) {
    // ...
}

func create_v110(r *webapi.Request) (int, webapi.Response) {
    // ...
}

func (item Item) Post(r *webapi.Request) (int, webapi.Response) {
    return webapi.Handlers(r, []Handler{
        webapi.Apply(create_v100, Version("1.0.0")),
        webapi.Apply(create_v110, Version("1.1.0")),
    })
}

// versioning middleware
func Version(version string) Middleware {
    return func(handler webapi.Handler) webapi.Handler {
        return func(r *webapi.Request) (int, webapi.Response) {
            requestedVersion := "1.1.0" // extract this somehow from the request

            if version != requestedVersion {
                return webapi.Next()
            }
            return handler(r)
        }
    }
}
```

So what's happening here...

- We have 2 different implementations of the same endpoint (v100 & v110).
- We pass the request object and all of our handlers to the `webapi.Handlers` function.
- We apply our `Version`-middleware to each handler using `webapi.Apply`, which in this case would be the same thing as `Version("1.0.0")(create_v100)`.
- In our middleware we check if the requested version matches the one for the current handler, if not return `webapi.Next()` and otherwise return the handler response. 

### Middleware

Any function that takes and returns a `webapi.Handler`.

```go
// always return I'm a teapot status code
func Teapot(handler webapi.Handler) webapi.Handler {
    return func(r *webapi.Request) (int, webapi.Response) {
        _, data := handler(r)
        return 418, data
    }
}
```

And apply with either

```go
api.Apply(Teapot)
```

which will apply it to all resources added after it, or 

```go
api.Add(`/items/(?P<id>\d+)$`, &Item{}, Teapot)
```

to just apply it to the given resources. You can add any number of middleware.

```go
api.Add(`/items/(?P<id>\d+)$`, &Item{}, Teapot, AnotherMiddleware, AndSoOn)
```

You can also apply middleware directly to a handler with the `webapi.Apply`-function, such in the multiple handlers per endpoint example above.

### URL parameters

You can use any valid regex when specifying paths for your resources. If you specify some named groups you can also access those as URL parameters.

```go
api.Add(`/items/(?P<id>\d+)$`, &Item{})

// ... later inside of a request handler

id := request.Param("id")
```

### Request body

You can easily unmarshal an JSON request body using `request.UnmarshalBody`.

```go
var data interface{}

err := request.UnmarshalBody(&data)
if err != nil {
    // handle appropriately 
}
```

## Example

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/lohmander/webapi"
)

func main() {
    api := webapi.NewAPI()
    api.Apply(Logger)
    api.Add(`/items/(?P<id>\d+)$`, &Item{}, Teapot)

    http.Handle("/api/", api)
    http.ListenAndServe(":3002", nil)
}

type Item struct{}

func (item Item) Post(request *webapi.Request) (int, webapi.Response) {
    var body interface{}

    err := request.UnmarshalBody(&body)
    if err != nil {
        return 500, webapi.Response{
            Error: err,
        }
    }

    return 200, webapi.Response{
        Data: map[string]interface{}{
            "body":    body,
            "idParam": request.Param("id"),
        },
    }
}

// some middleware

func Logger(handler webapi.Handler) webapi.Handler {
    return func(r *webapi.Request) (int, webapi.Response) {
        code, data := handler(r)
        fmt.Println(code, r.Method, r.URL.Path)
        return code, data
    }
}

func Teapot(handler webapi.Handler) webapi.Handler {
    return func(r *webapi.Request) (int, webapi.Response) {
        _, data := handler(r)
        return 418, data
    }
}

```
```sh
> curl -i -X POST -d '{"greeting": "Hello world!"}' http://localhost:3002/api/items/123

HTTP/1.1 418 I'm a teapot
Content-Type: application/json
Date: Thu, 12 Nov 2015 22:37:02 GMT
Content-Length: 61

{"data":{"body":{"greeting":"Hello world!"},"idParam":"123"}}
```

## License

[MIT](https://github.com/lohmander/webapi/blob/master/LICENSE), see **LICENSE** file.
