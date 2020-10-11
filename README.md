# balena-sdk-go

A Balena SDK in Go

## Description

An SDK for the Balena API written in Go(lang).
Currently it's a work in progress and the client implementation as well as its (public) API will change.

## Usage

```go
import (
    "fmt"
    "github.com/hslatman/balena-sdk-go/pkg/client"
    "github.com/hslatman/balena-sdk-go/pkg/logger"
)

func main() {

    // Create a new Client with ClientModifiers
    logger := logger.NullLogger{}
    token := "<your-balena-api-token>"
    c, err := client.New(
        token,
        client.WithLogger(logger),
        client.WithTimeout(45*time.Second),
        client.WithDebug(),
        client.WithTrace(),
    )

    if err != nil {
        fmt.Println(err)
        return
    }

    // Retrieve and loop through applications
    apps, err := c.Applications()
    if err != nil {
        fmt.Println(err)
        return
    }

    for _, app := range apps {
        fmt.Println(app)
    }
}
```

## OData

A [PR](https://github.com/hslatman/balena-sdk-go/pull/4) is in progress to handle the OData API in a more fluent manner.

## TODO

* More types of resources
* Add operations for changing resources (POST, PATCH, DELETE)
* Add handling of $select, $filter, $expand, etc. (see above)
