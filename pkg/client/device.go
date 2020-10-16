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
	"fmt"
	"strconv"

	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/m7shapan/njson"
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
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) Device(deviceID int) *DeviceResource {
	return NewDeviceResource(c, deviceID)
}

func (r *DeviceResource) Get() (models.Device, error) {

	device := models.Device{}

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return device, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) device

	if !first.Exists() {
		return device, fmt.Errorf("%s not found", r.endpoint)
	}

	if err := njson.Unmarshal([]byte(first.Raw), &device); err != nil {
		return device, err
	}

	return device, nil
}

func (r *DeviceResource) Select(s string) *DeviceResource {
	r.modifiers.AddSelect(s) // TODO: add validation that fields to be selected are valid fields for Device?
	return r
}

func (r *DeviceResource) Tags() *DeviceTagsResource {
	tr := NewDeviceTagsResource(
		r.client,
	)
	// TODO: improve the formatting of this filter; Balena API seems to require this like this?
	tr.modifiers.AddFilter("device/id%20eq%20'" + strconv.Itoa(r.deviceID) + "'")
	return tr
}

func (r *DeviceResource) AddTag(key string, value string) (models.DeviceTag, error) {

	tr := NewDeviceTagsResource(
		r.client,
	)

	return tr.Create(r.deviceID, key, value)
}

func (r *DeviceResource) UpdateTag(key string, value string) error {

	tr := NewDeviceTagsResource(
		r.client,
	)

	tr.modifiers.AddFilter("device/id%20eq%20'" + strconv.Itoa(r.deviceID) + "'%20and%20tag_key%20eq%20'" + key + "'")

	return tr.Update(value)
}

func (r *DeviceResource) DeleteTag(key string) error {

	tr := NewDeviceTagsResource(
		r.client,
	)

	tr.modifiers.AddFilter("device/id%20eq%20'" + strconv.Itoa(r.deviceID) + "'%20and%20tag_key%20eq%20'" + key + "'")

	return tr.Delete()
}
