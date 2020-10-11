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

type DeviceTagsResource struct {
	client    *Client
	endpoint  string
	modifiers *ODataModifiers
}

func NewDeviceTagsResource(c *Client) *DeviceTagsResource {
	return &DeviceTagsResource{
		client:    c,
		endpoint:  string(deviceTagsEndpoint),
		modifiers: c.NewODataModifiers(),
	}
}

func (c *Client) DeviceTags() *DeviceTagsResource {
	return NewDeviceTagsResource(c)
}

func (d *DeviceTagsResource) GetByID(deviceTagID int) *DeviceTagResource {
	return NewDeviceTagResource(
		d.client,
		deviceTagID,
	)
}

func (d *DeviceTagsResource) Get() (map[int]models.DeviceTag, error) {

	tags := make(map[int]models.DeviceTag)

	resp, err := d.client.get(d.endpoint, d.modifiers)

	if err != nil {
		return tags, err
	}

	// TODO: decide whether we need gjson dependency, or can do it easily, with a bit more wrangling, ourselves
	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, tag := range data.Array() {
		t := models.DeviceTag{}
		if err := njson.Unmarshal([]byte(tag.Raw), &t); err != nil {
			return tags, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		tags[t.ID] = t
	}

	return tags, nil

}

func (d *DeviceTagsResource) Select(s string) *DeviceTagsResource {
	d.modifiers.AddSelect(s)
	return d
}
