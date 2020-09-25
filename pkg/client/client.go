package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	version = "0.1.0"

	defaultBaseURL           = "https://api.balena-cloud.com/v5"
	defaultUserAgent         = "https://github.com/hslatman/balena-sdk-go"
	defaultContentTypeHeader = "application/json"
	defaultAcceptHeader      = "application/json"
	defaultTimeOut           = 30 * time.Second
)

type endpoint string

const (
	allApplicationsEndpoint endpoint = "/application"
	applicationsEndpoint    endpoint = "/my_application"
	devicesEndpoint         endpoint = "/device"
)

type ClientModifier func(c *Client)

type Client struct {
	rc           *resty.Client
	modifiers    []ClientModifier
	logger       Logger
	debugEnabled bool
	traceEnabled bool
}

func New(token string, modifiers ...ClientModifier) (*Client, error) {

	// TODO: add additional configuration options / modifiers?
	// TODO: default retries? support for proxy? TLS settings? other transports?

	// Creating a new Resty client with defaults for all requests
	rc := resty.New()

	rc.SetHostURL(defaultBaseURL)
	rc.SetHeader("User-Agent", defaultUserAgent) // TODO: add version?
	rc.SetHeader("Content-Type", defaultContentTypeHeader)
	rc.SetHeader("Accept", defaultAcceptHeader)

	rc.SetAuthScheme("Bearer")
	rc.SetAuthToken(token)

	rc.SetTimeout(defaultTimeOut)
	rc.SetDebug(false)

	c := &Client{
		rc:           rc,
		modifiers:    []ClientModifier{},
		debugEnabled: false,
		traceEnabled: false,
	}

	c.modifiers = append(c.modifiers, modifiers...)
	for _, modifier := range c.modifiers {
		modifier(c)
	}

	fmt.Println(rc)
	fmt.Println(c)

	return c, nil
}

func WithLogger(logger Logger) ClientModifier {
	return func(c *Client) {
		c.logger = logger
	}
}

func WithTimeout(timeout time.Duration) ClientModifier {
	return func(c *Client) {
		c.rc.SetTimeout(timeout)
	}
}

func WithDebug() ClientModifier {
	return func(c *Client) {
		c.debugEnabled = true
		c.rc.SetDebug(true)
	}
}

func WithTrace() ClientModifier {
	return func(c *Client) {
		c.traceEnabled = true
		c.rc.EnableTrace()
	}
}

func (c *Client) send(method string, url string) (*resty.Response, error) {

	resp, err := c.rc.R().Execute(method, url)
	if err != nil {
		return nil, err
	}

	if c.debugEnabled {
		c.debug(resp.Status())
		c.debug(resp.String())
	}

	if c.traceEnabled {
		// TODO: do more with the TraceInfo struct?
		c.info("trace: " + fmt.Sprint(resp.Request.TraceInfo()))
	}

	statusCode := resp.StatusCode()
	if statusCode < 200 || statusCode >= 300 { // TODO: check that this is OK
		return nil, fmt.Errorf("response failed with status: %d (%s)", statusCode, http.StatusText(statusCode))
	}

	return resp, nil
}
