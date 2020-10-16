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

	"github.com/hslatman/balena-sdk-go/pkg/models"
	"github.com/m7shapan/njson"
	"github.com/tidwall/gjson"
)

type DeviceTagResource struct {
	client      *Client
	endpoint    string
	deviceTagID int
	modifiers   *ODataModifiers
}

func NewDeviceTagResource(c *Client, deviceTagID int) *DeviceTagResource {
	return &DeviceTagResource{
		client:      c,
		endpoint:    fmt.Sprintf("%s(%d)", deviceTagsEndpoint, deviceTagID),
		deviceTagID: deviceTagID,
		modifiers:   NewODataModifiers(c),
	}
}

func (c *Client) DeviceTag(deviceTagID int) *DeviceTagResource {
	return NewDeviceTagResource(c, deviceTagID)
}

func (r *DeviceTagResource) Get() (models.DeviceTag, error) {

	tag := models.DeviceTag{}

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return tag, err
	}

	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results
	first := data.Get("0")                   // first (and only) tag, if found

	if !first.Exists() {
		return tag, fmt.Errorf("%s not found", r.endpoint)
	}

	if err := njson.Unmarshal([]byte(first.Raw), &tag); err != nil {
		return tag, err
	}

	return tag, nil
}

func (r *DeviceTagResource) Update(value string) error {

	body := map[string]interface{}{
		"value": value,
	}

	_, err := r.client.patch(r.endpoint, r.modifiers, body)

	if err != nil {
		return err
	}

	return nil
}

func (r *DeviceTagResource) Delete() error {

	_, err := r.client.delete(r.endpoint, r.modifiers)
	if err != nil {
		return err
	}

	return nil
}

func (r *DeviceTagResource) Select(s string) *DeviceTagResource {
	r.modifiers.AddSelect(s)
	return r
}

func (r *DeviceTagResource) Expand(s string) *DeviceTagResource {
	r.modifiers.AddExpand(s)
	return r
}
