package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/tidwall/gjson"
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

func (c *Client) Applications() (map[int]models.Application, error) {

	apps := make(map[int]models.Application)

	resp, err := c.rc.R().
		Get("/my_application")

	if err != nil {
		return apps, err
	}

	// TODO: add $select and $filter

	fmt.Println(resp.StatusCode())
	fmt.Println(resp)
	fmt.Println(resp.Request.TraceInfo())

	// TODO: "fluent" APIs? Like, from Device call, instantly get to Application, or something like that?

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, app := range data.Array() {
		a := models.Application{}
		if err := json.Unmarshal([]byte(app.Raw), &a); err != nil {
			return apps, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		apps[a.ID] = a
	}

	return apps, nil
}

func (c *Client) Application(id int) (models.Application, error) {

	// TODO: also lookup by other identifiers, like UUID, device name, device type, etc

	app := models.Application{}

	resp, err := c.rc.R().
		Get(fmt.Sprintf("/my_application(%d)", id))

	if err != nil {
		return app, err
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp)
	fmt.Println(resp.Request.TraceInfo())

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) application, if found

	if !first.Exists() {
		return app, fmt.Errorf("application (@id=%d) not found", id)
	}

	if err := json.Unmarshal([]byte(first.Raw), &app); err != nil {
		return app, err
	}

	return app, nil
}

func (c *Client) Devices() (map[int]models.Device, error) {

	devices := make(map[int]models.Device)

	resp, err := c.rc.R().
		Get("/device")

	if err != nil {
		return devices, err
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp)
	fmt.Println(resp.Request.TraceInfo())

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, device := range data.Array() {
		d := models.Device{}
		if err := json.Unmarshal([]byte(device.Raw), &d); err != nil {
			return devices, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		devices[d.ID] = d
	}

	return devices, nil
}

func (c *Client) Device(id int) (models.Device, error) {

	// TODO: also lookup by other identifiers, like UUID, device name, device type, etc

	device := models.Device{}

	resp, err := c.rc.R().
		Get(fmt.Sprintf("/device(%d)", id))

	if err != nil {
		return device, err
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp)
	fmt.Println(resp.Request.TraceInfo())

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) device

	if !first.Exists() {
		return device, fmt.Errorf("device (@id=%d) not found", id)
	}

	if err := json.Unmarshal([]byte(first.Raw), &device); err != nil {
		return device, err
	}

	return device, nil
}

func (c *Client) AllApplications() (map[int]models.Application, error) {

	apps := make(map[int]models.Application)

	resp, err := c.rc.R().
		Get("/application")

	if err != nil {
		return apps, err
	}

	fmt.Println(resp.StatusCode())
	fmt.Println(resp)
	fmt.Println(resp.Request.TraceInfo())

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, app := range data.Array() {
		a := models.Application{}
		if err := json.Unmarshal([]byte(app.Raw), &a); err != nil {
			return apps, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		apps[a.ID] = a
	}

	return apps, nil
}
