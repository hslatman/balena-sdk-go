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


    // Retrieve and loop through devices
    devices, err := c.Devices().Get()
    if err != nil {
        fmt.Println(err)
    }

    for _, d := range devices {
        fmt.Println(d)
    }

    // Retrieve id, name and type for the device with ID 1337
    dr := c.Device(1337).Select("id,device_name,device_type")
    device, err := dr.Get()
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(device)

    // Retrieve the device tags for the device with ID 1337
    dtr := c.Device(1337).Tags()
    deviceTags, err := dtr.Get()
    if err != nil {
        fmt.Println(err)
    }

    for _, t := range deviceTags {
        fmt.Println(t)
    }

}
```

## OData

This SDK uses the Balena OData API to retrieve and update data.
The OData implementation is not generic and mostly geared towards the Balena API.
It was inspired by the implementation in [gosip](https://github.com/koltyakov/gosip).

## TODO

* Add authentication method; currently requires API token to be retrieved from portal.
* More types of resources; many are currently missing.
* Add operations for changing resources (POST, PATCH, DELETE).
* Nicer [OData](https://www.odata.org/) implementation? Currently no real mature Go library available.
* Add HTTP caching? For example using https://github.com/gregjones/httpcache. 
* Use fastjson? (https://github.com/valyala/fastjson)
* Implement ID type (i.e. DeviceID, TagID, etc.)
* More model fields; they are incomplete now.
* Add code generation for Select(), Filter(), Expand(), etc?
* Add typesafe handling of OData modifiers (i.e. eq, neq, etc.)
* Add chaining of multiple subselects (if possible)
* Add context.Context to requests (in the resources object?)
* "Multi-layer" OData modifier handling