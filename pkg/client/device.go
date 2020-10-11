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
	"fmt"
	"strconv"

	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/tidwall/gjson"
)

type DeviceResource struct {
	client    *Client
	endpoint  string
	deviceID  int
	modifiers *ODataModifiers
}

func NewDeviceResource(c *Client, deviceID int) *DeviceResource {
	return &DeviceResource{
		client:    c,
		endpoint:  fmt.Sprintf("%s(%d)", devicesEndpoint, deviceID),
		deviceID:  deviceID,
		modifiers: c.NewODataModifiers(),
	}
}

func (c *Client) Device(deviceID int) *DeviceResource {
	return NewDeviceResource(c, deviceID)
}

func (d *DeviceResource) Get() (models.Device, error) {

	device := models.Device{}

	resp, err := d.client.get(d.endpoint, d.modifiers)

	if err != nil {
		return device, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) device

	if !first.Exists() {
		return device, fmt.Errorf("%s not found", d.endpoint)
	}

	if err := json.Unmarshal([]byte(first.Raw), &device); err != nil {
		return device, err
	}

	return device, nil
}

func (d *DeviceResource) Select(s string) *DeviceResource {
	d.modifiers.AddSelect(s) // TODO: add validation that fields to be selected are valid fields for Device?
	return d
}

func (d *DeviceResource) Tags() *DeviceTagsResource {
	r := NewDeviceTagsResource(
		d.client,
	)
	// TODO: improve the formatting of this filter; Balena API seems to require this like this?
	r.modifiers.AddFilter("device/id%20eq%20'" + strconv.Itoa(d.deviceID) + "'")
	return r
}
