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
)

func main() {
	token := "<your-balena-api-token>"
	c, err := client.New(token)
	if err != nil {
		fmt.Println(err)
		return
	}

	apps, err := c.Applications()
	if err != nil {
		fmt.Println(err)
	}

	for _, app := range apps {
		fmt.Println(app)
	}
}
```

## TODO

* Nicer [OData](https://www.odata.org/) implementation. Currently no real mature Go library available.
* Fluent APIs? Seems to fit the OData paradigm.
