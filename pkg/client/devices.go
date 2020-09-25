package client

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/tidwall/gjson"
)

func (c *Client) Devices() (map[int]models.Device, error) {

	devices := make(map[int]models.Device)

	resp, err := c.send(resty.MethodGet, string(devicesEndpoint))

	if err != nil {
		return devices, err
	}

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

	resp, err := c.send(resty.MethodGet, fmt.Sprintf("%s(%d)", devicesEndpoint, id))

	if err != nil {
		return device, err
	}

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
