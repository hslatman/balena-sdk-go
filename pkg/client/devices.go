// Copyright 2020 Herman Slatman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"encoding/json"

	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/tidwall/gjson"
)

type DevicesResource struct {
	client    *Client
	endpoint  string
	modifiers *ODataModifiers

	// TODO: context, configuration
}

func NewDevicesResource(c *Client) *DevicesResource {
	return &DevicesResource{
		client:    c,
		endpoint:  string(devicesEndpoint),
		modifiers: c.NewODataModifiers(),
	}
}

func (c *Client) Devices() *DevicesResource {
	return NewDevicesResource(c)
}

func (d *DevicesResource) Select(s string) *DevicesResource {
	d.modifiers.AddSelect(s)
	return d
}

func (d *DevicesResource) Filter(f string) *DevicesResource {
	d.modifiers.AddFilter(f)
	return d
}

func (d *DevicesResource) Get() (map[int]models.Device, error) {

	devices := make(map[int]models.Device)

	resp, err := d.client.get(d.endpoint, d.modifiers)

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

func (d *DevicesResource) GetByID(deviceID int) *DeviceResource {
	return NewDeviceResource(
		d.client,
		deviceID,
	)
}
