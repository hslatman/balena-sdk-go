package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	rc *resty.Client
}

func New(token string) (*Client, error) {

	// TODO: add additional configuration options / modifiers?
	// TODO: add logging?

	// Creating a new Resty client with defaults for all requests
	rc := resty.New()

	rc.SetHostURL("https://api.balena-cloud.com/v5")
	rc.SetHeader("User-Agent", "https://github.com/hslatman/balena-sdk-go") // TODO: add version?
	rc.SetHeader("Content-Type", "application/json")
	rc.SetHeader("Accept", "application/json")

	rc.SetAuthScheme("Bearer")
	rc.SetAuthToken(token)

	rc.SetTimeout(30 * time.Second)

	rc.EnableTrace()  // TODO: make this optional?
	rc.SetDebug(true) // TODO: make this optional?

	// TODO: default retries? support for proxy? TLS settings? other transports?

	fmt.Println(rc)

	return &Client{
		rc: rc,
	}, nil
}

func (c *Client) send(method string, url string) (*resty.Response, error) {

	resp, err := c.rc.R().Execute(method, url)
	if err != nil {
		return nil, err
	}

	// TODO: make these configurable; log to something different, etc.
	fmt.Println(resp.StatusCode())
	fmt.Println(resp)
	fmt.Println(resp.Request.TraceInfo())

	statusCode := resp.StatusCode()
	if statusCode < 200 || statusCode >= 300 { // TODO: check that this is OK
		return nil, fmt.Errorf("response failed with status: %d (%s)", statusCode, http.StatusText(statusCode))
	}

	return resp, nil
}
