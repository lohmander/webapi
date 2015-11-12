# webapi

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
    var someData interface{} = map[string]string{
        "key": "value"
    }
    return 200, webapi.Response{
        Data: someData
    }
}
```

### Middleware

Any function that takes and returns a `webapi.Handler`.

```go
// always return I'm a teaopot status code
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

which will apply it to all endpoints added after it, or 

```go
api.Add(`/items/(?P<id>\d+)$`, &Item{}, Teapot)
```

to just apply it to the given endpoint. You can add any number of endpoints.

```go
api.Add(`/items/(?P<id>\d+)$`, &Item{}, Teapot, AnotherMiddleware, AndSoOn)
```

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
