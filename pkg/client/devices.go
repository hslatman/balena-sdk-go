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
	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/m7shapan/njson"
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
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) Devices() *DevicesResource {
	return NewDevicesResource(c)
}

func (r *DevicesResource) Select(s string) *DevicesResource {
	r.modifiers.AddSelect(s)
	return r
}

func (r *DevicesResource) Filter(f string) *DevicesResource {
	r.modifiers.AddFilter(f)
	return r
}

func (r *DevicesResource) Get() (map[int]models.Device, error) {

	devices := make(map[int]models.Device)

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return devices, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, d := range data.Array() {
		device := models.Device{}
		if err := njson.Unmarshal([]byte(d.Raw), &device); err != nil {
			return devices, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		devices[device.ID] = device
	}

	return devices, nil
}

func (r *DevicesResource) FindByID(deviceID int) *DeviceResource {
	return NewDeviceResource(
		r.client,
		deviceID,
	)
}
