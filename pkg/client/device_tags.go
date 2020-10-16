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
		modifiers: NewODataModifiers(c),
	}
}

func (c *Client) DeviceTags() *DeviceTagsResource {
	return NewDeviceTagsResource(c)
}

func (r *DeviceTagsResource) FindByID(deviceTagID int) *DeviceTagResource {
	return NewDeviceTagResource(
		r.client,
		deviceTagID,
	)
}

func (r *DeviceTagsResource) Get() (map[int]models.DeviceTag, error) {

	tags := make(map[int]models.DeviceTag)

	resp, err := r.client.get(r.endpoint, r.modifiers)

	if err != nil {
		return tags, err
	}

	// TODO: decide whether we need gjson dependency, or can do it easily, with a bit more wrangling, ourselves
	data := gjson.GetBytes(resp.Body(), "d") // get data; a list of results

	for _, t := range data.Array() {
		tag := models.DeviceTag{}
		if err := njson.Unmarshal([]byte(t.Raw), &tag); err != nil {
			return tags, err // TODO: don't do early return, but just skip this one and aggregate error?
		}
		tags[tag.ID] = tag
	}

	return tags, nil

}

func (r *DeviceTagsResource) Create(deviceID int, key string, value string) (models.DeviceTag, error) {

	tag := models.DeviceTag{}

	body := map[string]interface{}{
		"device":  deviceID,
		"tag_key": key,
		"value":   value,
	}

	resp, err := r.client.post(r.endpoint, r.modifiers, body)
	if err != nil {
		return tag, err
	}

	if err := njson.Unmarshal(resp.Body(), &tag); err != nil {
		return tag, err
	}

	return tag, nil
}

func (r *DeviceTagsResource) Update(value string) error {

	// TODO: should modifiers (also) be handled here?

	body := map[string]interface{}{
		"value": value,
	}

	_, err := r.client.patch(r.endpoint, r.modifiers, body)
	if err != nil {
		return err
	}

	return nil
}

func (r *DeviceTagsResource) Delete() error {

	// TODO: should modifiers (also) be handled here?

	_, err := r.client.delete(r.endpoint, r.modifiers)
	if err != nil {
		return err
	}

	return nil
}

func (r *DeviceTagsResource) Select(s string) *DeviceTagsResource {
	r.modifiers.AddSelect(s)
	return r
}

func (r *DeviceTagsResource) Filter(s string) *DeviceTagsResource {
	r.modifiers.AddFilter(s)
	return r
}
